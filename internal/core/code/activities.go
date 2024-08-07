package code

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"regexp"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/shared"
)

type (
	RepoActivities struct{}
)

func (a *RepoActivities) SignalBranch(ctx context.Context, payload *defs.RepoIOSignalBranchCtrlPayload) error {
	args := make([]any, 0)
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("repo"),
		shared.WithWorkflowBlockID(payload.Repo.ID.String()),
		shared.WithWorkflowElement("branch"),
		shared.WithWorkflowElementID(payload.Branch),
	)

	args = append(args, payload.Repo)

	var workflow any
	if payload.Repo.DefaultBranch == payload.Branch {
		workflow = TrunkCtrl
	} else {
		workflow = BranchCtrl

		args = append(args, payload.Branch)
	}

	_, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(
			context.Background(),
			opts.ID,
			payload.Signal.String(),
			payload.Payload,
			opts,
			workflow,
			args...,
		)

	return err
}

// CloneBranch clones a repository branch at the temporary location, as specified by the payload.
// It uses the RepoIO interface to get the url with the oauth token in it.
// If an error occurs retrieving the clone URL, it is returned.
func (a *RepoActivities) CloneBranch(ctx context.Context, payload *defs.RepoIOClonePayload) error {
	url, err := kernel.
		Instance().
		RepoIO(payload.Repo.Provider).
		TokenizedCloneURL(
			ctx,
			&defs.RepoIOProviderInfo{
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
func (a RepoActivities) FetchBranch(ctx context.Context, payload *defs.RepoIOClonePayload) error {
	cmd := exec.CommandContext(ctx, "git", "-C", payload.Path, "fetch", "origin", payload.Repo.DefaultBranch) //nolint

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

// RebaseAtCommit attempts to rebase the repository at the given commit.
// If the rebase fails, it returns the SHA and error message of the failed commit.
// If the rebase is in progress, it returns an InProgress flag.
func (a *RepoActivities) RebaseAtCommit(ctx context.Context, payload *defs.RepoIOClonePayload) (*defs.RepoIORebaseAtCommitResponse, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", payload.Path, "rebase", payload.Push.After) // nolint

	var exerr *exec.ExitError

	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.As(err, &exerr) {
			switch exerr.ExitCode() {
			case 1:
				stderr := string(out)
				pattern := `(?m)^Could not apply ([0-9a-fA-F]{7})\.\.\. (.*)$`

				// Compile the regex
				re := regexp.MustCompile(pattern)

				// Find all matches
				matches := re.FindAllStringSubmatch(stderr, -1)

				if len(matches) > 0 {
					sha, msg := matches[0][1], matches[0][2]

					return &defs.RepoIORebaseAtCommitResponse{SHA: sha, Message: msg}, NewRebaseError(sha, msg)
				}
			case 128:
				return &defs.RepoIORebaseAtCommitResponse{InProgress: true}, NewRebaseError("unknown", "unknown")
			default:
				return nil, err // retry
			}
		}

		return nil, err // retry
	}

	return nil, nil // rebase successful
}

// Push pushes the contents of the repository at the given path to the remote.
// If force is true, the push will be forced (--force).
func (a *RepoActivities) Push(ctx context.Context, payload *defs.RepoIOPushBranchPayload) error {
	args := []string{"-C", payload.Path, "push", "origin", payload.Branch}
	if payload.Force {
		args = append(args, "--force")
	}

	cmd := exec.CommandContext(ctx, "git", args...)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

// RemoveClonedAtPath removes the repo that was cloned at the given path.
// If the path does not exist, RemoveClonedAtPath will return a nil error.
func (a *RepoActivities) RemoveClonedAtPath(ctx context.Context, path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}
