// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	gh "github.com/google/go-github/v53/github"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

const (
	CIEvent     = "CI"
	BuildEvent  = "Build"
	DeployEvent = "Deploy"
)

// handleInstallationEvent handles GitHub App installation event.
func handleInstallationEvent(ctx echo.Context) error {
	payload := &InstallationEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	shared.Logger().Info("installation event received ...")

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement(WebhookEventInstallation.String()),
		)

	exe, err := shared.Temporal().Client().SignalWithStartWorkflow(
		ctx.Request().Context(),
		opts.ID,
		WorkflowSignalInstallationEvent.String(),
		payload,
		opts,
		workflows.OnInstallationEvent,
	)
	if err != nil {
		shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
		return nil
	}

	shared.Logger().Debug("installation event handled ...", "options", opts, "execution", exe.GetRunID())

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
}

func processInProgressWorkflowRun(ghWorkflowEvent *GithubWorkflowEvent, eventType string) error {
	shared.Logger().Debug("processInProgressWorkflowRun started.")

	// only saving to DB here as in-progress along with the type of event it is (e.g CI, build, deploy)
	githubEventState := &GithubEventsState{}
	eventsData := make(map[string]any)

	owner := *ghWorkflowEvent.WR.Repository.Owner.Login
	repo := *ghWorkflowEvent.WR.Repository.Name
	workflowID := *ghWorkflowEvent.WR.WorkflowID
	workflowRunID := *ghWorkflowEvent.WR.ID

	if eventType == CIEvent {
		CIBranchName := *ghWorkflowEvent.WR.HeadBranch
		eventsData["branch"] = CIBranchName
	}

	if eventType == BuildEvent || eventType == DeployEvent {
		ghClient, err := Instance().GetClientFromInstallation(ghWorkflowEvent.Installation.ID)
		if err != nil {
			shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		}

		buildCommit := *ghWorkflowEvent.WR.HeadSHA

		// Get all tags for the repository
		tags, _, err := ghClient.Repositories.ListTags(context.Background(), owner, repo, nil)
		if err != nil {
			return err
		}

		// Iterate over tags to find the latest tag associated with the commit
		var latestTag *gh.RepositoryTag

		// TODO: make sure we get the latest tag associated with the `buildCommit``
		// Filter tags that point to the given commit
		var commitTags []*gh.RepositoryTag

		for _, tag := range tags {
			commit, _, err := ghClient.Repositories.GetCommit(context.Background(), owner, repo, tag.GetCommit().GetSHA(), nil)
			if err != nil {
				return err
			}

			if buildCommit == commit.GetSHA() {
				commitTags = append(commitTags, tag)
			}
		}

		latestTag = commitTags[0]
		shared.Logger().Debug("processInProgressWorkflowRun for "+repo+" event "+eventType, "latestTag", latestTag)

		if latestTag == nil {
			shared.Logger().Error("latest tag not found for " + repo + " event " + eventType)
		}

		eventsData["changesetID"] = *latestTag.Name
	}

	shared.Logger().Debug("processInProgressWorkflowRun for "+repo+" event "+eventType, "eventsData", eventsData)

	// githubEventState.ID, _ = gocql.RandomUUID()
	githubEventState.Status = "Inprogress"
	githubEventState.GithubWorkflowID = workflowID
	githubEventState.GithubWorkflowRunID = workflowRunID
	githubEventState.EventType = eventType
	githubEventState.RepoName = repo
	jsonData, _ := json.Marshal(eventsData)
	githubEventState.EventsData = string(jsonData)

	// save the CI event
	if err := db.Save(githubEventState); err != nil {
		shared.Logger().Error("error saving to github_events_state", "error", err)
		return err
	}

	shared.Logger().Debug("processInProgressWorkflowRun", "githubEventState saved to db", githubEventState)

	return nil
}

func processCompletedWorkflowRun(ctx echo.Context, ghWorkflowEvent *GithubWorkflowEvent, eventType string) error {
	githubEventsState := &GithubEventsState{}
	params := db.QueryParams{
		"github_workflow_run_id": strconv.FormatInt(*ghWorkflowEvent.WR.ID, 10),
	}

	shared.Logger().Debug("processCompletedWorkflowRun"+eventType, "params", params)

	if err := db.Get(githubEventsState, params); err != nil {
		shared.Logger().Error("processCompletedWorkflowRun for "+eventType, "error retrieving from db", err)
		return err
	}

	githubEventsState.Status = "Done"
	shared.Logger().Debug("processCompletedWorkflowRun"+eventType, "githubEventsState", githubEventsState)

	var eventsData map[string]string
	_ = json.Unmarshal([]byte(githubEventsState.EventsData), &eventsData)

	switch eventType {
	case "CI":
		// trigger CI related temporal workflow
		workflows := &Workflows{}
		workflowID := shared.Temporal().
			Queue(shared.ProvidersQueue).
			WorkflowID(
				shared.WithWorkflowBlock("CI-github"),
				shared.WithWorkflowBlockID(strconv.FormatInt(ghWorkflowEvent.Repository.ID, 10)),
				shared.WithWorkflowElement("branch"),
				shared.WithWorkflowElementID(*ghWorkflowEvent.WR.HeadBranch),
			)
		opts := shared.Temporal().
			Queue(shared.ProvidersQueue).
			WorkflowOptions(
				shared.WithWorkflowBlock("CI-github"),
				shared.WithWorkflowBlockID(strconv.FormatInt(ghWorkflowEvent.Repository.ID, 10)),
				shared.WithWorkflowElement("branch"),
				shared.WithWorkflowElementID(*ghWorkflowEvent.WR.HeadBranch),
			)

		shared.Logger().Debug("Triggering CI process action")

		if _, err := shared.Temporal().Client().SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			WorkflowSignalActionResult.String(),
			ghWorkflowEvent,
			opts,
			workflows.OnGithubCIAction,
			ghWorkflowEvent,
		); err != nil {
			shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
			return err
		}

		// save db
		// ghEvents := &GithubEventsState{}
		// params := db.QueryParams{
		// 	// "github_workflow_id":     strconv.FormatInt(*ghWorkflowEvent.WR.WorkflowID, 10),
		// 	"github_workflow_run_id": strconv.FormatInt(*ghWorkflowEvent.WR.ID, 10),
		// 	// "event_type":             "CI",
		// }

		// if err := db.Get(ghEvents, params); err != nil {
		// 	shared.Logger().Error("processCompletedWorkflowRun", "error getting data from db", err)
		// 	return err
		// }

		githubEventsState.Status = "Done"
		if err := db.Update(githubEventsState); err != nil {
			shared.Logger().Error("processCompletedWorkflowRun "+eventType, "error updating db", err)
		}

		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: workflowID, Status: WorkflowStatusSignaled})

	case "Build":
		// trigger Build temporal workflow
		workflows := &Workflows{}
		workflowID := shared.Temporal().
			Queue(shared.ProvidersQueue).
			WorkflowID(
				shared.WithWorkflowBlock("Build-github"),
				shared.WithWorkflowBlockID(strconv.FormatInt(ghWorkflowEvent.Repository.ID, 10)),
				shared.WithWorkflowElement("changeset"),
				shared.WithWorkflowElementID(eventsData["changesetID"]),
			)
		opts := shared.Temporal().
			Queue(shared.ProvidersQueue).
			WorkflowOptions(
				shared.WithWorkflowBlock("Build-github"),
				shared.WithWorkflowBlockID(strconv.FormatInt(ghWorkflowEvent.Repository.ID, 10)),
				shared.WithWorkflowElement("changeset"),
				shared.WithWorkflowElementID(eventsData["changesetID"]),
			)

		if _, err := shared.Temporal().Client().SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			WorkflowSignalActionResult.String(),
			ghWorkflowEvent,
			opts,
			workflows.OnGithubBuildAction,
			ghWorkflowEvent,
			eventsData["changesetID"],
		); err != nil {
			shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
			return err
		}

		// save db
		// ghEvents := &GithubEventsState{}
		// params := db.QueryParams{
		// 	"github_workflow_id":     strconv.FormatInt(*ghWorkflowEvent.WR.WorkflowID, 10),
		// 	"github_workflow_run_id": strconv.FormatInt(*ghWorkflowEvent.WR.ID, 10),
		// 	"event_type":             "CI",
		// }

		// if err := db.Get(ghEvents, params); err != nil {
		// 	return err
		// }

		githubEventsState.Status = "Done"
		if err := db.Update(githubEventsState); err != nil {
			shared.Logger().Error("processCompletedWorkflowRun "+eventType, "error updating db", err)
		}

		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: workflowID, Status: WorkflowStatusSignaled})

	case "Deploy":
		// trigger deploy temporal workflow
		workflows := &Workflows{}
		workflowID := shared.Temporal().
			Queue(shared.ProvidersQueue).
			WorkflowID(
				shared.WithWorkflowBlock("Deploy-github"),
				shared.WithWorkflowBlockID(strconv.FormatInt(ghWorkflowEvent.Repository.ID, 10)),
				shared.WithWorkflowElement("changeset"),
				shared.WithWorkflowElementID(eventsData["changesetID"]),
			)
		opts := shared.Temporal().
			Queue(shared.ProvidersQueue).
			WorkflowOptions(
				shared.WithWorkflowBlock("Deploy-github"),
				shared.WithWorkflowBlockID(strconv.FormatInt(ghWorkflowEvent.Repository.ID, 10)),
				shared.WithWorkflowElement("changeset"),
				shared.WithWorkflowElementID(eventsData["changesetID"]),
			)

		if _, err := shared.Temporal().Client().SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			WorkflowSignalActionResult.String(),
			ghWorkflowEvent,
			opts,
			workflows.OnGithubDeployAction,
			ghWorkflowEvent,
			eventsData["changesetID"],
		); err != nil {
			shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
			return err
		}

		// // save db
		// ghEvents := &GithubEventsState{}
		// params := db.QueryParams{
		// 	"github_workflow_id":     strconv.FormatInt(*ghWorkflowEvent.WR.WorkflowID, 10),
		// 	"github_workflow_run_id": strconv.FormatInt(*ghWorkflowEvent.WR.ID, 10),
		// 	"event_type":             "CI",
		// }

		// if err := db.Get(ghEvents, params); err != nil {
		// 	return err
		// }

		githubEventsState.Status = "Done"
		if err := db.Update(githubEventsState); err != nil {
			shared.Logger().Error("processCompletedWorkflowRun "+eventType, "error updating db", err)
		}

		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: workflowID, Status: WorkflowStatusSignaled})
	}

	return nil
}

func handleWorkflowRunEvent(ctx echo.Context) error {
	shared.Logger().Debug("workflow-run event received.")

	payload := &GithubWorkflowEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	// e := &events.Event{
	// 	Provider:   "github",
	// 	ProviderID: payload.Repository.ID,
	// 	Name:       "github workflow run done",
	// }
	// e.Save()

	// shared.Logger().Debug("handleWorkflowRunEvent", "payload", payload)

	parts := strings.Split(*payload.Workflow.Path, "/")
	wf_file := parts[len(parts)-1]

	var eventType string

	if wf_file == "cicd_quantm.yaml" {
		eventType = "CI"
	} else if wf_file == "build_quantm.yaml" {
		eventType = "Build"
	} else if wf_file == "deploy_quantm.yaml" {
		eventType = "Deploy"
	} else {
		// return "unregistered github workflow event received"
		shared.Logger().Warn("handleWorkflowRunEvent invalid workflow file related event received.")
		return nil
	}

	shared.Logger().Debug("handleWorkflowRunEvent", "eventType", eventType)

	if payload.Action == "in_progress" {
		shared.Logger().Info("workflow_run event in progress")

		if err := processInProgressWorkflowRun(payload, eventType); err != nil {
			shared.Logger().Error("process in-progress workflow", "error", err)
			return err
		}

		return nil
	}

	if payload.Action == "completed" {
		shared.Logger().Info("workflow_run event completed")

		if err := processCompletedWorkflowRun(ctx, payload, eventType); err != nil {
			shared.Logger().Error("process complete workflow", "error", err)
			return err
		}

		return nil
	}

	return nil
}

// handlePushEvent handles GitHub push event.
func handlePushEvent(ctx echo.Context) error {
	payload := &PushEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	// the value will be `NoCommit` if we have a tag push, or squash merge.
	if payload.After == NoCommit {
		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: db.NullUUID, Status: WorkflowStatusSkipped})
	}

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(strconv.FormatInt(payload.Repository.ID, 10)),
			shared.WithWorkflowMod(WebhookEventPush.String()),
			shared.WithWorkflowProp("ref", payload.Ref),
		)

	exe, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.OnPushEvent, payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
}

// handlePullRequestEvent handles GitHub pull request event.
func handlePullRequestEvent(ctx echo.Context) error {
	payload := &PullRequestEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("handlePullRequestEvent executing...")

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(strconv.FormatInt(payload.Repository.ID, 10)),
			shared.WithWorkflowMod(WebhookEventPullRequest.String()),
			shared.WithWorkflowModID(fmt.Sprintf(
				"%s-%s",
				strconv.FormatInt(payload.PullRequest.ID, 10),
				payload.Action,
			)),
		)

	switch payload.Action {
	case "opened":
		// e := &events.Event{
		// 	Provider:   "github",
		// 	ProviderID: payload.Repository.ID,
		// 	Name:       "PR opened event received",
		// }
		// e.Save()
		shared.Logger().Info("PR", "status", "open")
		// exe, err := shared.Temporal().
		// 	Client().
		// 	ExecuteWorkflow(context.Background(), opts, w.OnPullRequestEvent, payload)

		// if err != nil {
		// 	return shared.NewAPIError(http.StatusInternalServerError, err)
		// }

		// return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
		return nil

	case "labeled":
		// e := &events.Event{
		// 	Provider:   "github",
		// 	ProviderID: payload.Repository.ID,
		// 	Name:       "PR labeled event received",
		// }
		// e.Save()
		shared.Logger().Info("PR", "status", "label")
		exe, err := shared.Temporal().
			Client().
			ExecuteWorkflow(context.Background(), opts, w.OnLabelEvent, payload)

		if err != nil {
			return shared.NewAPIError(http.StatusInternalServerError, err)
		}

		return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})

	default:
		shared.Logger().Debug("handlePullRequestEvent default closing...")

		// e := &events.Event{
		// 	Provider:   "github",
		// 	ProviderID: payload.Repository.ID,
		// 	Name:       "Unregistered PR event received",
		// }
		// e.Save()
		// err := shared.Temporal().
		// 	Client().
		// 	SignalWorkflow(context.Background(), opts.ID, "", WebhookEventPullRequest.String(), payload)
		// if err != nil {
		// 	shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
		// 	return shared.NewAPIError(http.StatusInternalServerError, err)
		// }

		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: db.NullUUID, Status: WorkflowStatusSkipped})
	}
}

// handleInstallationRepositoriesEvent handles GitHub installation repositories event.
func handleInstallationRepositoriesEvent(ctx echo.Context) error {
	payload := &InstallationRepositoriesEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement(WebhookEventInstallationRepositories.String()),
		)

	exe, err := shared.Temporal().
		Client().
		ExecuteWorkflow(context.Background(), opts, w.OnInstallationRepositoriesEvent, payload)
	if err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: exe.GetID(), Status: WorkflowStatusQueued})
}
