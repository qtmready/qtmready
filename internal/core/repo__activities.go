package core

import (
	"context"

	"go.breu.io/quantm/internal/shared"
)

type (
	RepoActivities struct{}
)

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

	return nil
}

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
