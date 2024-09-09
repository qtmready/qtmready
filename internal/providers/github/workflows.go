// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.


package github

import (
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
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

// OnInstallationEvent workflow is executed when we initiate the installation of GitHub defs.
//
// In an ideal world, the complete installation request would hit the API after the installation event has hit the
// webhook, however, there can be number of things that can go wrong, and we can receive the complete installation
// request before the push event. To handle this, we use temporal.io's signal API to provide two possible entry points
// for the system. See the README.md for a detailed explanation on how this workflow works.
//
// NOTE: This workflow is only meant to be started with SignalWithStartWorkflow.
// TODO: Refactor this workflow to reduce complexity.
func (w *Workflows) OnInstallationEvent(ctx workflow.Context) (*Installation, error) { // nolint:funlen
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
	selector.AddReceive(requestChannel, onInstallationRequestSignal(ctx, request, status))

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
		installation.InstallationLogin = webhook.Installation.Account.Login // Github Organization name
		installation.InstallationLoginID = webhook.Installation.Account.ID  // Github organization ID
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

	logger.Info("github/installation/post: syncing github org users ...", "installation_id", payload.InstallationID.String())

	// TODO: sync users
	orgsync := &SyncOrgUsersFromGithubPayload{
		InstallationID: payload.InstallationID,
		GithubOrgName:  payload.InstallationLogin,
		GithubOrgID:    payload.InstallationLoginID,
	}
	if err := workflow.ExecuteActivity(_ctx, activities.SyncOrgUsersFromGithub, orgsync).Get(_ctx, nil); err != nil {
		logger.Error("github/installation/post: error syncing org users ...", "error", err)
	}

	return nil
}

// OnPushEvent is run when ever a repo event is received. Repo Event can be push event or a create event.
func (w *Workflows) OnCreateOrDeleteEvent(ctx workflow.Context, payload *CreateOrDeleteEvent) error {
	logger := workflow.GetLogger(ctx)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)

	state, err := getRepoEventState(ctx, payload)
	if err != nil {
		logger.Warn("github/repo_event: unable to initialize event state ...", "error", err.Error())

		return nil // TODO: We should do some sort of notification because we have a faulty integration.
	}

	event := payload.normalize(state.CoreRepo, state.User)

	if err := workflow.
		ExecuteActivity(_ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalCreateOrDelete, event).
		Get(_ctx, nil); err != nil {
		logger.Warn(
			"github/repo_event: signal error, retrying ...",
			slog.Int64("github_repo__installation_id", payload.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", payload.Repository.ID.Int64()),
			slog.String("github_repo__id", state.Repo.ID.String()),
			slog.String("core_repo__id", state.Repo.ID.String()),
		)

		return err
	}

	return nil
}

// OnPushEvent is run when ever a repo event is received. Repo Event can be push event or a create event.
func (w *Workflows) OnPushEvent(ctx workflow.Context, payload *PushEvent) error {
	logger := workflow.GetLogger(ctx)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)

	state, err := getRepoEventState(ctx, payload)
	if err != nil {
		logger.Error("github/repo_event: unable to initialize event state ...", "error", err.Error())

		return nil // TODO: faulty integration
	}

	event := payload.normalize(state.CoreRepo, state.User)

	if err := workflow.
		ExecuteActivity(_ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPush, event).
		Get(_ctx, nil); err != nil {
		logger.Warn(
			"github/repo_event: signal error, retrying ...",
			slog.Int64("github_repo__installation_id", payload.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", payload.Repository.ID.Int64()),
			slog.String("github_repo__id", state.Repo.ID.String()),
			slog.String("core_repo__id", state.Repo.ID.String()),
		)
	}

	return nil
}

// OnPullRequestEvent normalize the pull request event and then signal the core repo.
func (w *Workflows) OnPullRequestEvent(ctx workflow.Context, payload *PullRequestEvent) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("github/pull_request: preparing ...")

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)

	state, err := getRepoEventState(ctx, payload)
	if err != nil {
		logger.Error("github/pull_request: error preparing ...", "error", err.Error())
		return err
	}

	event_pr := payload.normalize(state.CoreRepo, state.User)
	event_label := payload.as_label(event_pr) // this will be nil if scope is label

	fn := func() workflow.Future {
		if event_label == nil {
			return workflow.
				ExecuteActivity(_ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequest, event_pr)
		}

		return workflow.
			ExecuteActivity(
				_ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequestLabel, event_label,
			)
	}

	if err := fn().Get(_ctx, nil); err != nil {
		logger.Warn(
			"github/pull_request: signal error, retrying ...",
			slog.Int64("github_repo__installation_id", payload.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", payload.Repository.ID.Int64()),
			slog.String("github_repo__id", state.Repo.ID.String()),
			slog.String("core_repo__id", state.Repo.ID.String()),
		)

		return err
	}

	return nil
}

// OnPullRequestReviewEvent normalize the pull request review event and then signal the core repo.
func (w *Workflows) OnPullRequestReviewEvent(ctx workflow.Context, event *PullRequestReviewEvent) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("github/pull_request_review: preparing ...")

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)

	state, err := getRepoEventState(ctx, event)
	if err != nil {
		logger.Error("github/pull_request_review: error preparing ...", "error", err.Error())
		return err
	}

	payload := &defs.RepoIOSignalPullRequestReviewPayload{
		Action:         event.Action,
		Number:         event.Number,
		RepoName:       event.Repository.Name,
		RepoOwner:      event.Repository.Owner.Login,
		BaseBranch:     event.PullRequest.Base.Ref,
		HeadBranch:     event.PullRequest.Head.Ref,
		CtrlID:         state.Repo.ID.String(),
		InstallationID: event.Installation.ID,
		ProviderID:     state.Repo.GithubID.String(),
		User:           state.User,
	}

	if err := workflow.
		ExecuteActivity(_ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequestReview, payload).
		Get(_ctx, nil); err != nil {
		logger.Warn(
			"github/pull_request_review: error signaling repo ctrl ...",
			slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
			slog.String("github_repo__id", state.CoreRepo.ID.String()),
			slog.String("core_repo__id", state.CoreRepo.ID.String()),
		)

		return err
	}

	return nil
}

// OnPullRequestReviewCommentEvent normalize the pull request review comment event and then signal the core repo.
func (w *Workflows) OnPullRequestReviewCommentEvent(ctx workflow.Context, event *PullRequestReviewCommentEvent) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("github/pull_request_review_comment: preparing ...")

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)

	state, err := getRepoEventState(ctx, event)
	if err != nil {
		logger.Error("github/pull_request_review_comment: error preparing ...", "error", err.Error())
		return err
	}

	payload := &defs.RepoIOSignalPullRequestReviewCommentPayload{
		Action:         event.Action,
		Number:         event.Number,
		RepoName:       event.Repository.Name,
		RepoOwner:      event.Repository.Owner.Login,
		BaseBranch:     event.PullRequest.Base.Ref,
		HeadBranch:     event.PullRequest.Head.Ref,
		CtrlID:         state.Repo.ID.String(),
		InstallationID: event.Installation.ID,
		ProviderID:     state.Repo.GithubID.String(),
		User:           state.User,
	}

	if err := workflow.
		ExecuteActivity(_ctx, activities.SignalCoreRepoCtrl, state.CoreRepo, defs.RepoIOSignalPullRequestComment, payload).
		Get(_ctx, nil); err != nil {
		logger.Warn(
			"github/pull_request_review_comment: error signaling repo ctrl ...",
			slog.Int64("github_repo__installation_id", event.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", event.Repository.ID.Int64()),
			slog.String("github_repo__id", state.CoreRepo.ID.String()),
			slog.String("core_repo__id", state.CoreRepo.ID.String()),
		)

		return err
	}

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
		Get(actx, installation)
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

func (w *Workflows) OnWorkflowRunEvent(ctx workflow.Context, pl *GithubWorkflowRunEvent) error {
	logger := workflow.GetLogger(ctx)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, activityOpts)

	state, err := getRepoEventState(ctx, pl)
	if err != nil {
		logger.Error("github/workflow_run: error preparing ...", "error", err.Error())
		return err
	}

	payload := &defs.RepoIOSignalWorkflowRunPayload{
		Action:         pl.Action,
		RepoName:       pl.Repository.Name,
		RepoOwner:      pl.Repository.Owner.Login,
		CtrlID:         state.Repo.ID.String(),
		InstallationID: pl.Installation.ID,
		ProviderID:     state.Repo.GithubID.String(),
		User:           state.User,
	}

	p := &defs.RepoIOWorkflowActionPayload{
		RepoName:       payload.RepoName,
		RepoOwner:      payload.RepoOwner,
		InstallationID: payload.InstallationID,
	}

	winfo := &defs.RepoIOWorkflowInfo{}
	if err := workflow.ExecuteActivity(_ctx, activities.GithubWorkflowInfo, p).Get(ctx, winfo); err != nil {
		logger.Error(
			"github/workflow_run: error github workflow info...",
			slog.Int64("github_repo__installation_id", pl.Installation.ID.Int64()),
			slog.Int64("github_repo__github_id", pl.Repository.ID.Int64()),
			slog.String("github_repo__id", state.CoreRepo.ID.String()),
			slog.String("core_repo__id", state.CoreRepo.ID.String()),
		)

		return err
	}

	payload.WorkflowInfo = winfo
	logger.Info("OnWorkflowRunEvent/payload", "info", payload)

	// TODO - wrokflow logic

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

// onInstallationRequestSignal handles new http requests on an installation in progress.
func onInstallationRequestSignal(
	ctx workflow.Context, installation *CompleteInstallationSignal, status *InstallationWorkflowStatus,
) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("github/installation: received complete installation request ...")
		channel.Receive(ctx, installation)

		status.RequestDone = true
	}
}

func getRepoEventState(ctx workflow.Context, event RepoEvent) (*RepoEventState, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info(
		"github/repo_event: initializing state ...",
		slog.Int64("github_repo__installation_id", event.InstallationID().Int64()),
		slog.Int64("github_repo__github_id", event.RepoID().Int64()),
	)

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	_ctx := workflow.WithActivityOptions(ctx, opts)
	state := &RepoEventState{}
	repos := make([]Repo, 0)

	if err := workflow.
		ExecuteActivity(_ctx, activities.GetReposForInstallation, event.InstallationID().String(), event.RepoID().String()).
		Get(_ctx, &repos); err != nil {
		logger.Error("github/push: temporal error, aborting ... ")

		return state, err
	}

	if len(repos) == 0 {
		logger.Warn(
			"github/repo_event: no repos found ...",
			slog.Int64("github_repo__installation_id", event.InstallationID().Int64()),
			slog.Int64("github_repo__github_id", event.RepoID().Int64()),
		)

		return state, NewRepoNotFoundRepoEventError(event.InstallationID(), event.RepoID(), event.RepoName())
	}

	// TODO: handle the unique together case during installation.
	if len(repos) > 1 {
		logger.Warn(
			"github/repo_event: multiple repos found",
			slog.Int64("github_repo__installation_id", event.InstallationID().Int64()),
			slog.Int64("github_repo__github_id", event.RepoID().Int64()),
		)

		return state, NewMultipleReposFoundRepoEventError(event.InstallationID(), event.RepoID(), event.RepoName())
	}

	state.Repo = &repos[0]

	if !state.Repo.IsActive {
		logger.Warn(
			"github/repo_event: repo is not active",
			slog.Int64("github_repo__installation_id", event.InstallationID().Int64()),
			slog.Int64("github_repo__github_id", event.RepoID().Int64()),
		)

		return state, NewInactiveRepoRepoEventError(event.InstallationID(), event.RepoID(), event.RepoName())
	}

	if !state.Repo.HasEarlyWarning {
		logger.Warn(
			"github/repo_event: repo has no early warning",
			slog.Int64("github_repo__installation_id", event.InstallationID().Int64()),
			slog.Int64("github_repo__github_id", event.RepoID().Int64()),
		)

		return state, NewHasNoEarlyWarningRepoEventError(event.InstallationID(), event.RepoID(), event.RepoName())
	}

	if err := workflow.
		ExecuteActivity(_ctx, activities.GetCoreRepoByCtrlID, state.Repo.ID.String()).
		Get(_ctx, &state.CoreRepo); err != nil {
		logger.Warn(
			"github/repo_event: database error, retrying ... ",
			slog.Int64("github_repo__installation_id", event.InstallationID().Int64()),
			slog.Int64("github_repo__github_id", event.RepoID().Int64()),
			slog.String("github_repo__id", state.Repo.ID.String()),
			slog.String("core_repo__id", state.CoreRepo.ID.String()),
		)

		return state, err
	}

	if err := workflow.
		ExecuteActivity(_ctx, activities.GetTeamUserByLoginID, event.SenderID()).Get(_ctx, &state.User); err != nil {
		logger.Warn(
			"github/repo_event: database error, retrying ... ",
			slog.Int64("github_repo__installation_id", event.InstallationID().Int64()),
			slog.Int64("github_repo__github_id", event.RepoID().Int64()),
			slog.String("github_repo__id", state.Repo.ID.String()),
			slog.String("core_repo__id", state.CoreRepo.ID.String()),
		)

		return state, err
	}

	return state, nil
}
