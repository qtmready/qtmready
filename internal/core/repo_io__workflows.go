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

package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/timers"
	"go.breu.io/quantm/internal/shared"
)

type (
	RepoWorkflows struct {
		acts *RepoActivities
	}
)

// RepoCtrl is the controller for all the workflows related to the repository.
//
// NOTE: This workflow is only meant to be started with SignalWithStartWorkflow.
func (w *RepoWorkflows) RepoCtrl(ctx workflow.Context, repo *Repo) error {
	// prelude
	logger := NewRepoIOWorkflowLogger(ctx, repo, "repo_ctrl", "", "")
	selector := workflow.NewSelector(ctx)
	done := false

	// channels
	// push event signal
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, w.onRepoPush(ctx, repo)) // post processing for push event recieved on repo.

	// create_delete
	create_delete := workflow.GetSignalChannel(ctx, RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, w.onRepoCreateOrDelete(ctx, repo))

	// pull request channel
	pr := workflow.GetSignalChannel(ctx, RepoIOSignalPullRequest.String())
	selector.AddReceive(pr, w.onRepoPullRequest(ctx, repo)) // post processing for pull request event recieved on repo.

	logger.Info("init ...", "default_branch", repo.DefaultBranch)

	// TODO: need to come up with logic to shutdown when not required.
	for !done {
		selector.Select(ctx)
	}

	return nil
}

// DefaultBranchCtrl is the controller for the default branch.
func (w *RepoWorkflows) DefaultBranchCtrl(ctx workflow.Context, repo *Repo) error {
	// prelude
	logger := NewRepoIOWorkflowLogger(ctx, repo, "branch_ctrl", "", "default")
	selector := workflow.NewSelector(ctx)
	done := false

	// channels
	// push event signal
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, w.onDefaultBranchPush(ctx, repo)) // post processing for push event recieved on repo.

	logger.Info("init ...")

	for !done {
		selector.Select(ctx)
	}

	return nil
}

// BranchCtrl is the controller for all the branches except the default branch.
func (w *RepoWorkflows) BranchCtrl(ctx workflow.Context, repo *Repo, branch string) error {
	// prelude
	logger := NewRepoIOWorkflowLogger(ctx, repo, "branch_ctrl", "", branch)
	selector := workflow.NewSelector(ctx)
	done := false

	interval := timers.NewInterval(ctx, repo.StaleDuration.Duration)

	// handle stale check.
	workflow.Go(ctx, func(ctx workflow.Context) {
		for !done {
			interval.Next(ctx)

			opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
			ctx = workflow.WithActivityOptions(ctx, opts)

			// get the gtihub repo by ctrl_id
			data := &RepoIORepoData{}
			_ = workflow.ExecuteActivity(ctx, Instance().RepoIO(repo.Provider).GetRepoData, repo.CtrlID.String()).Get(ctx, data)

			// Only send message to provider channel
			msg := NewStaleBranchMessage(data, repo, branch)

			logger.Info("stale branch detected, sending message ...", "stale", msg.RepoUrl)

			_ = workflow.ExecuteActivity(ctx, Instance().MessageIO(repo.MessageProvider).SendStaleBranchMessage, msg)
		}
	})

	// push event signal.
	// detect changes. if changes are greater than threshold, send early warning message.
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, w.onBranchPush(ctx, repo, branch, interval)) // post processing for push event recieved on repo.

	// rebase signal.
	// attempts to rebase the branch with the base branch. if there are merge conflicts, sends message.
	rebase := workflow.GetSignalChannel(ctx, ReopIOSignalRebase.String())
	selector.AddReceive(rebase, w.onBranchRebase(ctx, repo, branch)) // post processing for early warning signal.

	// pr signal.
	pr := workflow.GetSignalChannel(ctx, RepoIOSignalPullRequest.String())
	selector.AddReceive(pr, w.onBranchPullRequest(ctx, repo, branch)) // post processing for pull request event recieved on repo.

	logger.Info("init ...")

	for !done {
		selector.Select(ctx)
	}

	logger.Info("done, exiting ...")

	return nil
}

// onRepoPush is a channel handler that is called when a repository is pushed to.
// It checks if the pushed branch is the default branch, and if so, signals the default branch.
// Otherwise, it signals the feature branch.
func (w *RepoWorkflows) onRepoPush(ctx workflow.Context, repo *Repo) shared.ChannelHandler {
	logger := NewRepoIOWorkflowLogger(ctx, repo, "repo_ctrl", "push", "")
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}

	ctx = workflow.WithActivityOptions(ctx, opts)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPushPayload{}
		channel.Receive(ctx, payload)

		logger.Info("init ...")

		if RefFromBranchName(repo.DefaultBranch) == payload.BranchRef {
			logger.Info("signaling default branch ...")

			err := workflow.ExecuteActivity(ctx, w.acts.SignalDefaultBranch, repo, RepoIOSignalPush, payload).Get(ctx, nil)
			if err != nil {
				logger.Warn("error signaling default branch, retrying ...", "error", err.Error())
			}

			return
		}

		logger.Info("signaling branch ...")

		branch := BranchNameFromRef(payload.BranchRef)

		err := workflow.ExecuteActivity(ctx, w.acts.SignalBranch, repo, RepoIOSignalPush, payload, branch).Get(ctx, nil)
		if err != nil {
			logger.Warn("error signaling branch, retrying ...", "error", err.Error())
		}
	}
}

// onDefaultBranchPush is a workflow channel handler that is triggered when the default branch of a repository is pushed to.
// It retrieves all branches in the repository, and signals for a rebase on any branches that are not the default branch and
// not a Quantm-created branch.
func (w *RepoWorkflows) onDefaultBranchPush(ctx workflow.Context, repo *Repo) shared.ChannelHandler {
	logger := NewRepoIOWorkflowLogger(ctx, repo, "branch_ctrl", "push", "default")
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}

	ctx = workflow.WithActivityOptions(ctx, opts)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPushPayload{}
		channel.Receive(ctx, payload)

		logger.Info("init ...", "sha", payload.After)

		// get all branches
		branches := []string{}
		info := &RepoIOInfoPayload{InstallationID: payload.InstallationID, RepoName: payload.RepoName, RepoOwner: payload.RepoOwner}

		if err := workflow.ExecuteActivity(ctx, Instance().RepoIO(repo.Provider).GetAllBranches, info).Get(ctx, &branches); err != nil {
			logger.Warn("error getting branches, retrying ...", "error", err.Error())
		}

		// signal to rebase branches that are not default and not quantm created
		for _, branch := range branches {
			if branch != repo.DefaultBranch && !IsQuantmBranch(branch) {
				logger.Info("signlaing branch controller to rebase ...", "target_branch", branch, "sha", payload.After)

				if err := workflow.ExecuteActivity(ctx, w.acts.SignalBranch, repo, ReopIOSignalRebase, payload, branch).Get(ctx, nil); err != nil { // nolint:revive
					logger.Warn("error sending rebase signal, retrying ...", "error", err.Error())
				}
			}
		}
	}
}

// onBranchPush is a shared.ChannelHandler that is called when a branch is pushed to a repository.
func (w *RepoWorkflows) onBranchPush(ctx workflow.Context, repo *Repo, branch string, interval timers.Interval) shared.ChannelHandler {
	logger := NewRepoIOWorkflowLogger(ctx, repo, "branch_ctrl", "push", branch)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Minute}

	ctx = workflow.WithActivityOptions(ctx, opts)

	interval.Restart(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPushPayload{}
		channel.Receive(ctx, payload)

		// detect changes payload -> RepoIODetectChangesPayload
		detect := &RepoIODetectChangesPayload{
			InstallationID: payload.InstallationID,
			RepoName:       payload.RepoName,
			RepoOwner:      payload.RepoOwner,
			DefaultBranch:  repo.DefaultBranch,
			TargetBranch:   branch,
		}
		changes := &RepoIOChanges{}

		// check the message provider if not ignore the early warning
		if repo.MessageProvider != MessageProviderNone {
			logger.Info("detecting changes ...", "sha", payload.After)

			_ = workflow.ExecuteActivity(ctx, Instance().RepoIO(repo.Provider).DetectChanges, detect).Get(ctx, changes)

			if changes.Delta > repo.Threshold {
				if payload.User != nil && payload.User.IsMessageProviderLinked {
					msg := NewNumberOfLinesExceedMessage(payload, repo, branch, changes, false)

					logger.Info("threshold exceeded ...", "sha", payload.After, "threshold", repo.Threshold, "delta", changes.Delta)

					_ = workflow.
						ExecuteActivity(ctx, Instance().MessageIO(repo.MessageProvider).SendNumberOfLinesExceedMessage, msg).
						Get(ctx, nil)

					// return the workflow is user exit not send message to channel
					return
				}

				// if user not exit then will send message to channel (repo message provider channel)
				msg := NewNumberOfLinesExceedMessage(payload, repo, branch, changes, true)

				logger.Info("threshold exceeded ...", "sha", payload.After, "threshold", repo.Threshold, "delta", changes.Delta)

				_ = workflow.
					ExecuteActivity(ctx, Instance().MessageIO(repo.MessageProvider).SendNumberOfLinesExceedMessage, msg).
					Get(ctx, nil)

				return
			}

			logger.Info("no changes detected ...", "sha", payload.After)

			return
		}

		// TODO: notify customer that message provider is not set.
		logger.Warn("message provider not set, ignoring early warning ...", "sha", payload.After)
	}
}

// onBranchRebase is a workflow handler that handles the rebase operation for a given repository and branch.
// It clones the repository, fetches the default branch, and then rebases the given branch at the specified commit.
// If a merge conflict is detected during the rebase, it sends a message via the message provider.
// In order to make sure that all the activities are executed on the same node, a session is created.
//
// NOTE:  _ = workflow.ExecuteActivity(ctx, ...) might look blasphempous! right? well, it's not. In the scope of the temporal workflow,
// the returned error is generally a temporal.ApplicationError, and will only happen if the number of retries is exhausted. We generally
// return the error to tell us that workflow has failed. In this case, we are not interested in the error.
func (w *RepoWorkflows) onBranchRebase(ctx workflow.Context, repo *Repo, branch string) shared.ChannelHandler {
	// _logger := workflow.GetLogger(ctx)
	// _log := w.logbranch(_logger, "push", repo.ID.String(), repo.Provider.String(), repo.ProviderID, branch)
	logger := NewRepoIOWorkflowLogger(ctx, repo, "branch_ctrl", "rebase", branch)
	retries := &temporal.RetryPolicy{NonRetryableErrorTypes: []string{"RepoIORebaseError"}}
	sopts := &workflow.SessionOptions{ExecutionTimeout: 30 * time.Minute, CreationTimeout: 60 * time.Minute} // TODO: make it configurable.
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second, RetryPolicy: retries}
	w.acts = &RepoActivities{}

	ctx = workflow.WithActivityOptions(ctx, opts)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPushPayload{}
		data := &RepoIOClonePayload{Repo: repo, Push: payload, Branch: branch, Path: ""}

		channel.Receive(ctx, payload)

		logger.Info("init ...", "sha", payload.After)

		_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any { return "/tmp/" + uuid.New().String() }).Get(&data.Path)

		sessionctx, _ := workflow.CreateSession(ctx, sopts)
		defer workflow.CompleteSession(sessionctx)

		logger.Info("cloning repo at branch ...", "sha", payload.After, "target_branch", branch, "path", data.Path)
		_ = workflow.ExecuteActivity(sessionctx, w.acts.CloneBranch, data).Get(sessionctx, nil)

		logger.Info("fetching default branch ...", "sha", payload.After, "path", data.Path)
		_ = workflow.ExecuteActivity(sessionctx, w.acts.FetchBranch, data).Get(sessionctx, nil)

		logger.Info("rebasing at commit ...", "sha", payload.After, "path", data.Path)

		rebase := &RepoIORebaseAtCommitResponse{}
		if err := workflow.ExecuteActivity(sessionctx, w.acts.RebaseAtCommit, data).Get(sessionctx, rebase); err != nil {
			var apperr *temporal.ApplicationError
			if errors.As(err, &apperr) {
				if apperr.Type() == "RepoIORebaseError" && !rebase.InProgress {
					logger.Info("merge conflict detected ...", "sha", rebase.SHA, "commit_message", rebase.Message, "path", data.Path)

					if payload.User != nil && payload.User.IsMessageProviderLinked {
						msg := NewMergeConflictMessage(payload, repo, branch, false)

						logger.Info("merge conflict detected, sending message ...", "sha", payload.After, payload.RepoName)

						_ = workflow.
							ExecuteActivity(ctx, Instance().MessageIO(repo.MessageProvider).SendMergeConflictsMessage, msg).
							Get(ctx, nil)

						_ = workflow.ExecuteActivity(ctx, w.acts.RemoveClonedAtPath, data.Path).Get(ctx, nil)

						// return the workflow is user exit not send message to channel
						return
					}

					// if user not exit then will send message to channel (repo message provider channel)
					msg := NewMergeConflictMessage(payload, repo, branch, true)

					logger.Info("merge conflict detected, sending message ...", "sha", payload.After, payload.RepoName)

					_ = workflow.ExecuteActivity(ctx, Instance().MessageIO(repo.MessageProvider).SendMergeConflictsMessage, msg)
					_ = workflow.ExecuteActivity(ctx, w.acts.RemoveClonedAtPath, data.Path).Get(ctx, nil)

					return
				}
			}

			_ = workflow.ExecuteActivity(ctx, w.acts.Push, branch, data.Path, true).Get(ctx, nil)
			_ = workflow.ExecuteActivity(ctx, w.acts.RemoveClonedAtPath, data.Path).Get(ctx, nil)

			return
		}

		_ = workflow.ExecuteActivity(ctx, w.acts.RemoveClonedAtPath, data.Path).Get(ctx, nil)
	}
}

func (w *RepoWorkflows) onRepoCreateOrDelete(ctx workflow.Context, repo *Repo) shared.ChannelHandler {
	logger := NewRepoIOWorkflowLogger(ctx, repo, "repo_ctrl", "create_delete", "")
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}

	ctx = workflow.WithActivityOptions(ctx, opts)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalCreatePayload{}
		channel.Receive(ctx, payload)

		logger.Info("init ...")
	}
}

func (w *RepoWorkflows) onRepoPullRequest(ctx workflow.Context, repo *Repo) shared.ChannelHandler {
	logger := NewRepoIOWorkflowLogger(ctx, repo, "repo_ctrl", "pull_request", "")
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}

	ctx = workflow.WithActivityOptions(ctx, opts)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPullRequestPayload{}
		channel.Receive(ctx, payload)

		logger.Info("init ...")

		_ = workflow.ExecuteActivity(ctx, w.acts.SignalBranch, repo, RepoIOSignalPullRequest, payload, payload.HeadBranch).Get(ctx, nil)
	}
}

func (w *RepoWorkflows) onBranchPullRequest(ctx workflow.Context, repo *Repo, branch string) shared.ChannelHandler {
	logger := NewRepoIOWorkflowLogger(ctx, repo, "branch_ctrl", "pull_request", branch)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}

	ctx = workflow.WithActivityOptions(ctx, opts)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPullRequestPayload{}
		channel.Receive(ctx, payload)

		logger.Info("on repo pull request", payload)
		logger.Info("on repo pull request action", payload.Action)
		logger.Info("repo branch on which pull request", branch)

		// TODO - convert to map call repo activites to handle the pr actions
		switch payload.Action {
		case "opened":
			logger.Info("pull request with open action")

		case "labeled":
			logger.Info("pull request with labeled action")

		case "synchronize":
			logger.Info("pull request with synchronize action")

		default:
			logger.Info("handlePullRequest Event default closing...")
		}
	}
}
