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
	"time"

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
	state := NewBranchCtrlState(ctx, repo, branch)

	state.check_stale(ctx) // start the stale check coroutine.

	// push event signal.
	// detect changes. if changes are greater than threshold, send early warning message.
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, w.on_branch_push(ctx, state))

	// rebase signal.
	// attempts to rebase the branch with the base branch. if there are merge conflicts, sends message.
	rebase := workflow.GetSignalChannel(ctx, ReopIOSignalRebase.String())
	selector.AddReceive(rebase, w.on_branch_rebase(ctx, state))

	// create_delete signal.
	// creates or deletes the branch.
	create_delete := workflow.GetSignalChannel(ctx, RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, w.on_branch_create_delete(ctx, state))

	// pr signal.
	pr := workflow.GetSignalChannel(ctx, RepoIOSignalPullRequest.String())
	selector.AddReceive(pr, w.on_branch_pr(ctx, state))

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

		err := workflow.ExecuteActivity(ctx, w.acts.SignalBranch_, repo, RepoIOSignalPush, payload, branch).Get(ctx, nil)
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

				if err := workflow.ExecuteActivity(ctx, w.acts.SignalBranch_, repo, ReopIOSignalRebase, payload, branch).Get(ctx, nil); err != nil { // nolint:revive
					logger.Warn("error sending rebase signal, retrying ...", "error", err.Error())
				}
			}
		}
	}
}

// on_branch_push is a shared.ChannelHandler that is called when commits are pushed to a branch. It handles the logic for
// detecting changes in the pushed branch and warning the user if the changes exceed a configured threshold.
func (w *RepoWorkflows) on_branch_push(ctx workflow.Context, state *BranchCtrlState) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPushPayload{}
		channel.Receive(ctx, payload)

		latest := payload.Commits.Latest()
		if latest != nil {
			state.set_commit(ctx, latest)
		}

		complexity := state.calculate_complexity(ctx, payload)
		if complexity.Delta > state.repo.Threshold {
			state.warn_complexity(ctx, payload, complexity)
		}

		state.interval.Restart(ctx)
	}
}

// on_branch_rebase is a shared.ChannelHandler that is called when a branch needs to be rebased. It handles the logic for
// cloning the repository, fetching the default branch, rebasing the branch at the latest commit, and pushing the rebased
// branch back to the repository.
func (w *RepoWorkflows) on_branch_rebase(ctx workflow.Context, state *BranchCtrlState) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}

		channel.Receive(ctx, push)

		session := state.create_session(ctx)
		defer workflow.CompleteSession(session)

		cloned := state.clone_at_commit(session, push)
		if cloned == nil {
			return
		}

		state.fetch_default_branch(session, cloned)

		if err := state.rebase_at_commit(session, cloned); err != nil {
			state.warn_conflict(ctx, push)
		}

		state.push_branch(session, cloned)
		state.remove_cloned(ctx, cloned)
	}
}

// on_branch_create_delete is a shared.ChannelHandler that is called when a branch is created or deleted. It handles the logic for
// updating the state of the branch control when a create or delete event is received.
//
// If the payload indicates the branch was created, the function sets the created timestamp in the state.
// If the payload indicates the branch was deleted, the function terminates the state.
func (w *RepoWorkflows) on_branch_create_delete(ctx workflow.Context, state *BranchCtrlState) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalCreateOrDeletePayload{}
		channel.Receive(ctx, payload)

		if payload.IsCreated {
			state.set_created_at(ctx, timers.Now(ctx))
		} else {
			state.terminate(ctx)
		}
	}
}

func (w *RepoWorkflows) on_branch_pr(ctx workflow.Context, state *BranchCtrlState) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalPullRequestPayload{}
		channel.Receive(ctx, payload)

		switch payload.Action {
		case "opened":
			state.set_pr(ctx, &RepoIOPullRequest{Number: payload.Number, HeadBranch: payload.HeadBranch, BaseBranch: payload.BaseBranch})
		default:
			state.log(ctx, "info", "pull_request", "unhandled action", "action", payload.Action)
		}
	}
}

func (w *RepoWorkflows) onRepoCreateOrDelete(ctx workflow.Context, repo *Repo) shared.ChannelHandler {
	logger := NewRepoIOWorkflowLogger(ctx, repo, "repo_ctrl", "create_delete", "")
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}

	ctx = workflow.WithActivityOptions(ctx, opts)

	return func(channel workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalCreateOrDeletePayload{}
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

		_ = workflow.ExecuteActivity(ctx, w.acts.SignalBranch_, repo, RepoIOSignalPullRequest, payload, payload.HeadBranch).Get(ctx, nil)
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
