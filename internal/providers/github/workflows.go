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
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

var (
	activities *Activities
)

type (
	// Workflows is the entry point for all workflows for GitHub.
	Workflows struct{}

	// InstallationWorkflowStatus handles the status of the workflow Workflows.OnInstallationEvent.
	InstallationWorkflowStatus struct {
		WebhookDone bool
		RequestDone bool
	}

	// PullRequestWorkflowStatus handles the status of the workflow Workflows.OnPullRequestEvent.
	PullRequestWorkflowStatus struct {
		Complete bool
	}
)

// OnInstallationEvent workflow is executed when we initiate the installation of GitHub core.
//
// In an ideal world, the complete installation request would hit the API after the installation event has hit the
// webhook, however, there can be number of things that can go wrong, and we can receive the complete installation
// request before the push event. To handle this, we use temporal.io's signal API to provide two possible entry points
// for the system. See the README.md for a detailed explanation on how this workflow works.
//
// NOTE: This workflow is only meant to be started with SignalWithStartWorkflow.
func (w *Workflows) OnInstallationEvent(ctx workflow.Context) (*Installation, error) {
	// prelude
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)
	installation := &Installation{}
	webhook := &InstallationEvent{}
	request := &CompleteInstallationSignal{}
	status := &InstallationWorkflowStatus{WebhookDone: false, RequestDone: false}
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, activityOpts)

	// setting up channels to receive signals
	webhookChannel := workflow.GetSignalChannel(ctx, WorkflowSignalInstallationEvent.String())
	requestChannel := workflow.GetSignalChannel(ctx, WorkflowSignalCompleteInstallation.String())

	// setting up callbacks for the channels
	selector.AddReceive(webhookChannel, onInstallationWebhookSignal(ctx, webhook, status))
	selector.AddReceive(requestChannel, onRequestSignal(ctx, request, status))

	logger.Info("github/installation: waiting for webhook and complete installation request signals ...")

	// keep listening for signals until we have received both the installation id and the team id
	for !(status.WebhookDone && status.RequestDone) {
		selector.Select(ctx)
	}

	logger.Info("github/installation: required signals processed ...")

	switch webhook.Action {
	// NOTE - Since a GitHub organization can only have one active installation at a time, when a new installation is created, it's
	// considered the first app installation for the organization, and we assume no teams have been created yet within the organization.
	//
	// TODO - we need to handle the case when an the app uninstallation and reinstallation case.
	//
	// - when delete event is received, we need to add a db field to mark the installation as deleted.
	// - on the subsequent installation, we need to check if the installation is deleted and update the installation status.
	case "created":
		user := &auth.User{}
		team := &auth.Team{}

		if err := workflow.ExecuteActivity(_ctx, activities.GetUserByID, request.UserID.String()).Get(ctx, user); err != nil {
			return nil, err
		}

		if user.TeamID.String() == db.NullUUID {
			logger.Info("github/installation: no team associated, creating a new team ...")

			team.Name = webhook.Installation.Account.Login

			_ = workflow.ExecuteActivity(_ctx, activities.CreateTeam, team).Get(ctx, team)

			logger.Info("github/installation: team created, assigning to user ...")

			user.TeamID = team.ID
			_ = workflow.ExecuteActivity(_ctx, activities.SaveUser, user).Get(ctx, user)
		} else {
			logger.Warn("github/installation: team already associated, fetching ...")

			_ = workflow.ExecuteActivity(_ctx, activities.GetTeamByID, user.TeamID.String()).Get(ctx, team)
		}

		// Finalizing the installation
		installation.TeamID = team.ID
		installation.InstallationID = webhook.Installation.ID
		installation.InstallationLogin = webhook.Installation.Account.Login
		installation.InstallationType = webhook.Installation.Account.Type
		installation.SenderID = webhook.Sender.ID
		installation.SenderLogin = webhook.Sender.Login
		installation.Status = webhook.Action

		logger.Info("github/installation: creating or updating installation ...")

		if err := workflow.ExecuteActivity(_ctx, activities.CreateOrUpdateInstallation, installation).Get(_ctx, installation); err != nil {
			logger.Error("github/installation: error saving installation ...", "error", err)
		}

		logger.Info("github/installation: updating user associations ...")

		membership := &CreateMembershipsPayload{
			UserID:        user.ID,
			TeamID:        team.ID,
			IsAdmin:       true,
			GithubOrgName: webhook.Installation.Account.Login,
			GithubOrgID:   webhook.Installation.Account.ID,
			GithubUserID:  webhook.Sender.ID,
		}

		if err := workflow.ExecuteActivity(_ctx, activities.CreateMemberships, membership).Get(_ctx, nil); err != nil {
			logger.Error("github/installation: error saving installation ...", "error", err)
		}

		logger.Info("github/installation: saving installation repos ...")

		for _, repo := range webhook.Repositories {
			logger.Info("github/installation: saving repository ...")
			logger.Debug("repository", "repository", repo)

			repo := &Repo{
				GithubID:        repo.ID,
				InstallationID:  installation.InstallationID,
				Name:            repo.Name,
				FullName:        repo.FullName,
				DefaultBranch:   "main",
				HasEarlyWarning: false,
				IsActive:        true,
				TeamID:          installation.TeamID,
			}

			future := workflow.ExecuteActivity(_ctx, activities.CreateOrUpdateGithubRepo, repo)

			// NOTE - ideally, we should use a new selector here, but since there will be no new signals comings in, we know that
			// selector.Select will only be waiting for the futures to complete.
			selector.AddFuture(future, onCreateOrUpdateRepoActivityFuture(ctx, repo))
		}

		logger.Info("github/installation: waiting for repositories to be saved ...")

		for range webhook.Repositories {
			selector.Select(ctx)
		}

		logger.Info("github/installation: installation repositories saved ...")
	case "deleted", "suspend", "unsuspend":
		logger.Warn("github/installation: installation removed, unhandled case ...")
	default:
		logger.Warn("github/installation: unhandled action during installation ...", slog.String("action", webhook.Action))
	}

	logger.Info("github/installation: complete", slog.Any("installation", installation))

	return installation, nil
}

// PostInstall refresh the default branch for all repositories associated with the given teamID and gets orgs users.
// NOTE - this workflow runs complete for the first time but when reinstall the github app and configure the same repos. it will give the,
// It will give the access_token error: could not refresh installation id XXXXXXX's token error.
// TODO - handle when the github app is reinstall and confgure the same repos,
// and also need to test when configure the same repo or new repos.
func (w *Workflows) PostInstall(ctx workflow.Context, payload *Installation) error {
	logger := workflow.GetLogger(ctx)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)

	logger.Info(
		"github/installation/post: starting ...",
		slog.String("installation_id", payload.InstallationID.String()),
		slog.String("installation_login", payload.InstallationLogin),
	)

	// TODO: move this inside a workflow.Go statement
	logger.Info("github/installation/post: syncing repos ...", "installation_id", payload.InstallationID.String())

	sync := &SyncReposFromGithubPayload{
		InstallationID: payload.InstallationID,
		Owner:          payload.InstallationLogin,
		TeamID:         payload.TeamID,
	}
	if err := workflow.ExecuteActivity(_ctx, activities.SyncReposFromGithub, sync).Get(_ctx, nil); err != nil {
		logger.Error("github/installation/post: error syncing repos ...", "error", err)
	}

	// TODO: sync users

	return nil
}

// OnPushEvent checks if the push event is associated with an open pull request.If so, it will get the idempotent key for
// the immutable rollout. Depending upon the target branch, it will either queue the rollout or update the existing
// rollout.
func (w *Workflows) OnPushEvent(ctx workflow.Context, event *PushEvent) error {
	logger := workflow.GetLogger(ctx)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)
	repos := make([]Repo, 0)
	corepo := &core.Repo{}

	logger.Info(
		"github/push: preparing ...",
		slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
		slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
	)

	if err := workflow.
		ExecuteActivity(_ctx, activities.GetReposForInstallation, event.Installation.ID.String(), event.Repository.ID.String()).
		Get(_ctx, &repos); err != nil {
		logger.Warn("github/push: database error, retrying ... ")
	}

	if len(repos) == 0 {
		logger.Warn(
			"github/push: unknown repo",
			slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
		)

		return nil
	}

	// TODO: handle the unique together case during installation.
	if len(repos) > 1 {
		logger.Warn(
			"github/push: multiple repos found",
			slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
		)
	}

	repo := repos[0]

	// TODO: notify the user that there is an unconfigured repo?
	if !repo.HasEarlyWarning || !repo.IsActive {
		logger.Warn(
			"webhook/push: uncofigured repo",
			slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
			slog.String("github_repo__id", repo.ID.String()),
		)

		return nil
	}

	if err := workflow.
		ExecuteActivity(_ctx, activities.GetCoreRepoByCtrlID, repo.ID.String()).
		Get(_ctx, corepo); err != nil {
		logger.Warn(
			"github/push: database error, retrying ... ",
			slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
			slog.String("github_repo__id", repo.ID.String()),
			slog.String("core_repo__id", corepo.ID.String()),
		)
	}

	logger.Info(
		"github/push: signal core repo ...",
		slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
		slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
		slog.String("github_repo__id", repo.ID.String()),
		slog.String("core_repo__id", corepo.ID.String()),
	)

	payload := &core.RepoSignalPushPayload{
		BranchRef:      event.Ref,
		Before:         event.Before,
		After:          event.After,
		RepoName:       event.Repository.Name,
		RepoOwner:      event.Repository.Owner.Login,
		CtrlID:         repo.ID.String(),
		InstallationID: event.Installation.ID,
		ProviderID:     repo.GithubID.String(),
	}

	if err := workflow.
		ExecuteActivity(_ctx, activities.SignalCoreRepoCtrl, corepo, core.RepoIOSignalPush, payload).
		Get(_ctx, nil); err != nil {
		logger.Warn(
			"github/push: signal error, retrying ...",
			slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
			slog.String("github_repo__id", repo.ID.String()),
			slog.String("core_repo__id", corepo.ID.String()),
		)
	}

	return nil
}

func (w *Workflows) OnWorkflowRunEvent(ctx workflow.Context, payload *GithubWorkflowRunEvent) error {
	logger := workflow.GetLogger(ctx)

	if actionWorkflowStatuses[payload.Repository.Name] != nil {
		logger.Debug("Workflow action file:", "action", payload.Action)
		logger.Debug("Workflow action file:", "file", payload.Workflow.Path)
		actionWorkflowStatuses[payload.Repository.Name][payload.Workflow.Path] = payload.Action
	}

	return nil
}

// After the creation of the idempotency key, we pass the idempotency key as a signal to the Aperture Workflow.
func (w *Workflows) OnPullRequestEvent(ctx workflow.Context, payload *PullRequestEvent) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("OnPullRequestEvent workflow started ...")
	// status := &PullRequestWorkflowStatus{Complete: false}

	// wait for artifact to generate and push to registery
	ch := workflow.GetSignalChannel(ctx, WorkflowSignalArtifactReady.String())
	artifact := &ArtifactReadySignal{}
	ch.Receive(ctx, artifact)

	// setting activity options
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	// get core repo
	repo := &Repo{GithubID: payload.Repository.ID}
	coreRepo := &core.Repo{}

	err := workflow.ExecuteActivity(actx, activities.GetCoreRepo, repo).Get(ctx, coreRepo)
	if err != nil {
		logger.Error("error getting core repo", "error", err)
		return err
	}

	// get core workflow ID for this stack
	corePRWfID := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowID(
			shared.WithWorkflowBlock("stack"),
			shared.WithWorkflowBlockID(coreRepo.StackID.String()),
		)

	// payload for core stack workflow
	signalPayload := &shared.PullRequestSignal{
		RepoID:           coreRepo.ID,
		SenderWorkflowID: workflow.GetInfo(ctx).WorkflowExecution.ID,
		TriggerID:        payload.PullRequest.ID,
		Image:            artifact.Image,
		Digest:           artifact.Digest,
		ImageRegistry:    artifact.Registry,
	}

	// signal core stack workflow
	logger.Info("core workflow id", "ID", corePRWfID)

	options := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("stack"),
			shared.WithWorkflowBlockID(coreRepo.StackID.String()),
		)

	cw := &core.StackWorkflows{}
	_, _ = shared.Temporal().Client().SignalWithStartWorkflow(
		context.Background(),
		corePRWfID,
		shared.WorkflowSignalDeploymentStarted.String(),
		signalPayload,
		options,
		cw.StackController,
		coreRepo.StackID.String(),
	)
	// workflow.SignalExternalWorkflow(ctx, corePRWfID, "", shared.WorkflowSignalPullRequest.String(), signalPayload).Get(ctx, nil)
	logger.Debug("Signaled workflow", "ID", signalPayload.SenderWorkflowID, " core repo ID: ", signalPayload.RepoID.String())

	// workflow.GetSignalChannel(ctx, WorkflowSignalPullRequestProcessed.String()).Receive(ctx, &status)

	// signal processor
	// selector.AddReceive(prChannel, onPRSignal(ctx, pr, status))

	// logger.Info("PR created: scheduling new aperture at the application level.")

	// // keep listening to signals until complete = true
	// for !status.Complete {
	// 	selector.Select(ctx)
	// }

	return nil
}

// OnInstallationRepositoriesEvent is responsible when a repository is added or removed from an installation.
func (w *Workflows) OnInstallationRepositoriesEvent(ctx workflow.Context, payload *InstallationRepositoriesEvent) error {
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)

	logger.Info("received installation repositories event ...")

	installation := &Installation{}
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	err := workflow.
		ExecuteActivity(actx, activities.GetInstallation, payload.Installation.ID).
		Get(ctx, installation)
	if err != nil {
		logger.Error("error getting installation", "error", err)
		return err
	}

	for _, repository := range payload.RepositoriesAdded {
		logger.Info("saving repository ...")
		logger.Debug("repository", "repository", repository)

		repo := &Repo{
			GithubID:       repository.ID,
			InstallationID: installation.InstallationID,
			Name:           repository.Name,
			FullName:       repository.FullName,
			TeamID:         installation.TeamID,
		}

		future := workflow.ExecuteActivity(actx, activities.CreateOrUpdateGithubRepo, repo)
		selector.AddFuture(future, onCreateOrUpdateRepoActivityFuture(ctx, repo))
	}

	// wait for all the repositories to be saved.
	for range payload.RepositoriesAdded {
		selector.Select(ctx)
	}

	return nil
}

// onCreateOrUpdateRepoActivityFuture handles post-processing after a repository is saved against an installation.
func onCreateOrUpdateRepoActivityFuture(ctx workflow.Context, payload *Repo) shared.FutureHandler {
	logger := workflow.GetLogger(ctx)
	return func(f workflow.Future) { logger.Info("repository saved ...", "repo", payload.GithubID) }
}

// onInstallationWebhookSignal handles webhook events for installation that is in progress.
func onInstallationWebhookSignal(
	ctx workflow.Context, installation *InstallationEvent, status *InstallationWorkflowStatus,
) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("github/installation: webhook received ...", "action", installation.Action)
		channel.Receive(ctx, installation)

		status.WebhookDone = true

		switch installation.Action {
		case "deleted", "suspend", "unsuspend":
			logger.Info("github/installation: installation removed ....", "action", installation.Action)

			status.RequestDone = true
		case "request":
			logger.Info("github/installation: installation request ...", "action", installation.Action)

			status.RequestDone = true
		default:
			logger.Info("github/installation: create action ...", "action", installation.Action)
		}
	}
}

// onRequestSignal handles new http requests on an installation in progress.
func onRequestSignal(
	ctx workflow.Context, installation *CompleteInstallationSignal, status *InstallationWorkflowStatus,
) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("github/installation: received complete installation request ...")
		channel.Receive(ctx, installation)

		status.RequestDone = true
	}
}

// func (w *Workflows) OnLabelEvent(ctx workflow.Context, payload *PullRequestEvent) error {
// 	logger := workflow.GetLogger(ctx)

// 	logger.Info("received PR label event ...")

// 	installationID := payload.Installation.ID
// 	repoOwner := payload.Repository.Owner.Login
// 	repoName := payload.Repository.Name
// 	pullRequestID := payload.Number
// 	label := payload.Label.Name
// 	branch := payload.PullRequest.Head.Ref

// 	switch label {
// 	case "quantm ready":
// 		logger.Debug("quantm ready label applied")

// 		cw := &core.RepoWorkflows{}
// 		opts := shared.Temporal().
// 			Queue(shared.CoreQueue).
// 			WorkflowOptions(
// 				shared.WithWorkflowBlock("repo"),
// 				shared.WithWorkflowBlockID(payload.Repository.ID.String()),
// 				shared.WithWorkflowElement("PR"),
// 				shared.WithWorkflowElementID(fmt.Sprint(pullRequestID)),
// 				shared.WithWorkflowProp("type", "merge_queue"),
// 			)

// 		payload2 := &shared.MergeQueueSignal{
// 			PullRequestID:  pullRequestID,
// 			InstallationID: installationID,
// 			RepoOwner:      repoOwner,
// 			RepoName:       repoName,
// 			Branch:         branch,
// 			RepoProvider:   "github",
// 		}

// 		if _, err := shared.Temporal().Client().
// 			SignalWithStartWorkflow(
// 				context.Background(),
// 				opts.ID,
// 				shared.MergeQueueStarted.String(),
// 				payload2,
// 				opts,
// 				cw.PollMergeQueue,
// 			); err != nil {
// 			logger.Error("OnLabelEvent: Error signaling workflow", "error", err)
// 			return err
// 		}

// 		logger.Info("PR sent to MergeQueue")

// 	case "quantm now":
// 		logger.Debug("quantm now label applied")

// 		// check if all workflows are completed!
// 		for {
// 			allCompleted := true

// 			for _, value := range actionWorkflowStatuses[repoName] {
// 				if value != "completed" {
// 					// return here since all are not completed
// 					allCompleted = false

// 					logger.Warn("all actions were not successful")

// 					break
// 				}
// 			}

// 			if allCompleted {
// 				break
// 			}

// 			_ = workflow.Sleep(ctx, 30*time.Second)

// 			logger.Debug("checking again all actions statuses")
// 		}

// 		// cw := &core.Workflows{}
// 		opts := shared.Temporal().
// 			Queue(shared.CoreQueue).
// 			WorkflowOptions(
// 				shared.WithWorkflowBlock("repo"),
// 				shared.WithWorkflowBlockID(payload.Repository.ID.String()),
// 				shared.WithWorkflowElement("PR"),
// 				shared.WithWorkflowElementID(fmt.Sprint(pullRequestID)),
// 				shared.WithWorkflowProp("type", "merge_queue"),
// 			)

// 		if err := shared.Temporal().Client().
// 			SignalWorkflow(context.Background(), opts.ID, "", shared.MergeTriggered.String(), nil); err != nil {
// 			logger.Error("OnLabelEvent: Error signaling workflow", "error", err)
// 			return err
// 		}

// 	default:
// 		logger.Debug("undefined label applied!")
// 	}

// 	return nil
// }

// OnPullRequestEvent workflow is responsible to get or create the idempotency key for the changeset controller workflow.
// Regardless of the action on PR, the algorithm needs to arrive at the same idempotency key! One possible way is
// to calculate the checksum  of different components. The trick would be to handle "synchronize" event as this relates
// to a new commit on the PR.
//
//   - One possible way to handle "synchronize" would be to only listen to label events on the PR.
//   - The other possible way to create an idempotency key would be to take the state, create a version set and then tag
//     the git commit with the version set. We can also take a look at aviator.co to see how they are creating version-sets.
//
// // onPRSignal handles incoming signals on open PR.
// func onPRSignal(ctx workflow.Context, pr *PullRequestEvent, status *PullRequestWorkflowStatus) shared.ChannelHandler {
// 	logger := workflow.GetLogger(ctx)

// 	return func(channel workflow.ReceiveChannel, more bool) {
// 		channel.Receive(ctx, pr)

// 		switch pr.Action {
// 		case "closed":
// 			logger.Info("PR closed: scheduling aperture to be abandoned.", "action", pr.Action)

// 			if pr.PullRequest.Merged {
// 				logger.Info("PR merged: scheduling aperture to finish with conclusion.")

// 				// TODO: send the signal to the aperture workflow to finish with conclusion.
// 				status.Complete = true
// 			} else {
// 				logger.Info("PR closed: abort aperture.")

// 				status.Complete = true
// 			}
// 		case "synchronize":
// 			logger.Info("PR updated: checking the status of the environment ...", "action", pr.Action)
// 			// TODO: here we need to check the app associated with repo & get the `release` branch. If the PR branch is not
// 			// the default branch, then we update in place, otherwise, we queue a new rollout.
// 		default:
// 			logger.Info("PR: no action required, skipping ...", "action", pr.Action)
// 		}
// 	}
// }
