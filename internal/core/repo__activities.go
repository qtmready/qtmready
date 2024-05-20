package core

import (
	"context"
	"errors"
	"fmt"
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
//   - on default branch push, rebase the commits from main branch onto the branch. if it fails, send a merge conflict warning.
//   - on push, tries to check if the number of lines changed is greater than the threshold defined on repo.
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
	cmd := exec.CommandContext(ctx, "git", "clone", "-b", payload.Branch, "--single-branch", "--depth", "1", url, payload.Path) //nolint

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	slog.Info(
		"repo cloned",
		slog.String("output", string(out)),
		slog.String("repo_id", payload.Repo.ID.String()),
		slog.String("branch", payload.Branch),
	)

	return nil
}

// FetchBranch fetches the given branch from the origin.
// TODO: right now this fetches the branch but we need to make it generic.
func (a RepoActivities) FetchBranch(ctx context.Context, payload *RepoIOClonePayload) error {
	cmd := exec.CommandContext(ctx, "git", "-C", payload.Path, "fetch", "origin", payload.Repo.DefaultBranch) //nolint

	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	slog.Info(
		"default branch fetched",
		slog.String("output", string(out)),
		slog.String("repo_id", payload.Repo.ID.String()),
		slog.String("branch", payload.Branch),
	)

	return nil
}

func (a *RepoActivities) RebaseAtCommit(ctx context.Context, payload RepoIOClonePayload) error {
	cmd := exec.CommandContext(ctx, "git", "-C", payload.Path, "rebase", payload.Push.After) // nolint

	var exerr *exec.ExitError

	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.As(err, &exerr) {
			slog.Info(
				"rebased",
				slog.String("error", err.Error()),
				slog.String("type", fmt.Sprintf("%T", err)),
				slog.String("output", string(out)),
				slog.String("repo_id", payload.Repo.ID.String()),
				slog.String("branch", payload.Branch),
			)

			// TODO: Read err.Error(). Find all lines containing a line like
			// Could not apply 1ac7130... chore() updating roles.tf file
			// where
			// - "Could not apply" signifies that rebasing failed
			// - 1ac7130 is the commit sha
			// - chore() updaing roles.tf file is the commit message.
			//
			// there are multiple stratagies we can take here.
			return nil
		}

		return err
	}

	slog.Info(
		"rebased",
		slog.String("output", string(out)),
		slog.String("repo_id", payload.Repo.ID.String()),
		slog.String("branch", payload.Branch),
	)

	return nil
}
