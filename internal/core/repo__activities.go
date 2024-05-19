package core

import (
	"context"

	"go.breu.io/quantm/internal/shared"
)

type (
	RepoActivities struct{}
)

// SignalDefaultBranch signals the default branch of a repository with a given workflow signal and payload.
// It uses Temporal to queue the workflow and passes the necessary options and parameters.
// If the signal and workflow start are successful, it returns nil. Otherwise, it returns an error.
func (a *RepoActivities) SignalDefaultBranch(ctx context.Context, repo *Repo, signal shared.WorkflowSignal, payload any) error {
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("repo"),
		shared.WithWorkflowBlockID(repo.ID.String()),
		shared.WithWorkflowElement("branch"),
		shared.WithWorkflowElementID(repo.DefaultBranch),
	)

	w := &RepoWorkflows{}

	_, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(context.Background(), opts.ID, signal.String(), payload, opts, w.DefaultBranchCtrl, repo)

	if err != nil {
		return err
	}

	shared.Logger().Info("signaled default branch", "repo", repo.ID, "signal", signal, "payload", payload)

	return nil
}

// SignalBranch signals a branch other than the default branch of a repository.
// This is mostly responsible for handling the early warning system.
//
//   - tries to rebase the commit on main back on to branch. if there are merge conflicts, sends message.
func (a *RepoActivities) SignalBranch(ctx context.Context, repo *Repo, signal shared.WorkflowSignal, payload any, branch string) error {
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("repo"),
		shared.WithWorkflowBlockID(repo.ID.String()),
		shared.WithWorkflowElement("branch"),
		shared.WithWorkflowElementID(branch),
	)

	w := &RepoWorkflows{}

	_, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(context.Background(), opts.ID, signal.String(), payload, opts, w.BranchCtrl, repo)

	if err != nil {
		return err
	}

	return nil
}
