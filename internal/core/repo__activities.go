package core

import (
	"context"
	"log/slog"
	"os/exec"

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
		SignalWithStartWorkflow(context.Background(), opts.ID, signal.String(), payload, opts, w.BranchCtrl, repo, branch)

	if err != nil {
		return err
	}

	return nil
}

// CloneBranch clones a repository branch at the temporary location, as specified by the payload.
// It uses the RepoIO interface to get the url with the oauth token in it.
// If an error occurs retrieving the clone URL, it is returned.
func (a *RepoActivities) CloneBranch(ctx context.Context, payload *RepoIOClonePayload) error {
	url, err := instance.
		RepoIO(payload.Repo.Provider).
		TokenizedCloneURL(
			ctx,
			&RepoIOInfoPayload{
				InstallationID: payload.Push.InstallationID,
				RepoName:       payload.Push.RepoName,
				RepoOwner:      payload.Push.RepoOwner,
			},
		)
	if err != nil {
		return err
	}

	slog.Info("clone url", slog.Any("url", url)) // TODO: remove in production.

	// NOTE: probably the method at https://stackoverflow.com/a/7349740 is much faster. Also see the comments.
	cmd := exec.Command("git", "clone", "-b", payload.Branch, "--single-branch", url, payload.Path) //nolint:gosec

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
