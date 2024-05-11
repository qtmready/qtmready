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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	gh "github.com/google/go-github/v53/github"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	ServerHandler struct{ *auth.SecurityHandler }
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(middleware echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &auth.SecurityHandler{Middleware: middleware},
	}
}

func (s *ServerHandler) GithubCompleteInstallation(ctx echo.Context) error {
	request := &CompleteInstallationRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, err := gocql.ParseUUID(ctx.Get("team_id").(string))

	if err != nil {
		return err
	}

	payload := &CompleteInstallationSignal{request.InstallationID, request.SetupAction, teamID}

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.Itoa(int(payload.InstallationID))),
			shared.WithWorkflowElement(WebhookEventInstallation.String()),
		)

	exe, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			WorkflowSignalCompleteInstallation.String(),
			payload,
			opts,
			workflows.OnInstallationEvent,
		)
	if err != nil {
		return err
	}

	_ = exe.Get(ctx.Request().Context(), nil)

	repos := make([]Repo, 0)

	if err := db.Filter(&Repo{}, &repos, db.QueryParams{"team_id": teamID.String()}); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, repos)
}

func (s *ServerHandler) GithubGetRepos(ctx echo.Context) error {
	result := make([]Repo, 0)
	if err := db.Filter(
		&Repo{},
		&result,
		db.QueryParams{"team_id": ctx.Get("team_id").(string)},
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *ServerHandler) GithubGetInstallations(ctx echo.Context) error {
	result := make([]Installation, 0)
	if err := db.Filter(
		&Installation{},
		&result,
		db.QueryParams{"team_id": ctx.Get("team_id").(string)},
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// GithubArtifactReady API is called by github action after building and pushing the build artifact
// After receiving pull request webhook, Quantum waits for the artifact ready event to start deployment.
func (s *ServerHandler) GithubArtifactReady(ctx echo.Context) error {
	request := &ArtifactReadyRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	workflowID := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowID(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(request.InstallationID),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(request.RepoID),
			shared.WithWorkflowMod(WebhookEventPullRequest.String()),
			shared.WithWorkflowModID(request.PullRequestID),
		)

	payload := &ArtifactReadySignal{Image: request.Image, Digest: request.Digest, Registry: request.Registry.String()}

	err := shared.Temporal().Client().SignalWorkflow(ctx.Request().Context(), workflowID, "", WorkflowSignalArtifactReady.String(), payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: workflowID, Status: WorkflowStatusSignaled})
}

func (s *ServerHandler) CliGitMerge(ctx echo.Context) error {
	shared.Logger().Info("CliGitMerge method triggered.")

	request := &CliGitMerge{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	shared.Logger().Info("CliGitMerge", "request", request)

	name := fmt.Sprintf("'%s/%s'", request.RepoOwner, request.RepoName)
	repo := &Repo{}

	if err := db.Get(repo, db.QueryParams{"full_name": name}); err != nil {
		shared.Logger().Error("Getting Repo data from database failed", "Error", err)
		return err
	}

	client, err := Instance().GetClientFromInstallation(repo.InstallationID)

	if err != nil {
		shared.Logger().Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// Get repository information to find the default branch
	repository, _, err := client.Repositories.Get(ctx.Request().Context(), request.RepoOwner, request.RepoName)
	if err != nil {
		shared.Logger().Error("client.Repositories.Get", "Error", err)
	}

	baseBranch := repository.GetDefaultBranch()
	PROptions := &gh.PullRequestListOptions{
		Base: baseBranch,
	}

	PRs, _, err := client.PullRequests.List(ctx.Request().Context(), request.RepoOwner, request.RepoName, PROptions)
	if err != nil {
		shared.Logger().Error("client.PullRequests.List failed", "Error", err)
		return err
	}

	PullRequestID := -1

	for i := 0; i < len(PRs); i++ {
		if *PRs[i].Head.Ref == request.Branch {
			PullRequestID = (*PRs[i].Number)
			break
		}
	}

	// Create PR for this branch and then label it
	if PullRequestID == -1 {
		// Specify the title and body for the pull request
		prTitle := "Pull Request created by Quantum"
		prBody := "Description of your pull request goes here."

		// Create a new pull request
		newPR := &gh.NewPullRequest{
			Title: &prTitle,
			Body:  &prBody,
			Head:  &request.Branch,
			Base:  &baseBranch,
		}

		pr, _, err := client.PullRequests.Create(ctx.Request().Context(), request.RepoOwner, request.RepoName, newPR)
		if err != nil {
			shared.Logger().Error("CliGitMerge", "Error creating pull request", err)
		}

		PullRequestID = *pr.Number

		shared.Logger().Info("CliGitMerge", "Pull request created", pr.GetNumber())
	}

	// Label the PR
	_, _, err = client.Issues.AddLabelsToIssue(ctx.Request().Context(), request.RepoOwner, request.RepoName, PullRequestID,
		[]string{"quantm ready"})
	if err != nil {
		shared.Logger().Error("CliGitMerge", "Error adding label to PR", err)
	}

	ret := fmt.Sprintf("PR %d is labeled", PullRequestID)

	return ctx.JSON(http.StatusOK, ret)
}

func (s *ServerHandler) GithubActionResult(ctx echo.Context) error {
	return nil
}

func (s *ServerHandler) GithubWebhook(ctx echo.Context) error {
	signature := ctx.Request().Header.Get("X-Hub-Signature-256")

	if signature == "" {
		return ctx.JSON(http.StatusUnauthorized, ErrMissingHeaderGithubSignature)
	}

	// NOTE: We are reading the request body twice. This is not ideal.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := Instance().VerifyWebhookSignature(body, signature); err != nil {
		return shared.NewAPIError(http.StatusUnauthorized, err)
	}

	headerEvent := ctx.Request().Header.Get("X-GitHub-Event")
	if headerEvent == "" {
		return shared.NewAPIError(http.StatusBadRequest, ErrMissingHeaderGithubEvent)
	}

	shared.Logger().Debug("GithubWebhook", "headerEvent", headerEvent)
	// Uncomment for debugging!
	// var jsonMap map[string]interface{}
	// json.Unmarshal([]byte(string(body)), &jsonMap)
	// shared.Logger().Debug("GithubWebhook", "body", jsonMap)

	event := WebhookEvent(headerEvent)
	handlers := WebhookEventHandlers{
		WebhookEventInstallation:             handleInstallationEvent,
		WebhookEventInstallationRepositories: handleInstallationRepositoriesEvent,
		WebhookEventPush:                     handlePushEvent,
		WebhookEventPullRequest:              handlePullRequestEvent,
		WebhookEventWorkflowRun:              handleWorkflowRunEvent,
	}

	if handle, exists := handlers[event]; exists {
		return handle(ctx)
	} else {
		shared.Logger().Warn("Github Webhook: Unsupported event", "event", event)
	}

	return shared.NewAPIError(http.StatusBadRequest, ErrInvalidEvent)
}
