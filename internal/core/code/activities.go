// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// # Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// # You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.
package code

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	// Activities defines an interface for repository-related actions.
	Activities struct{}

	// RepoEvent defines an interface for repository events. It simplifies type operations by wrapping
	// defs.Event[defs.EventPayload, defs.EventProvider].
	RepoEvent interface {
		Flatten() (db.Entity, error)
		MarshalJSON() ([]byte, error)
		UnmarshalJSON([]byte) error
	}
)

// SignalBranch signals a branch workflow for the given repository.
//
// If the branch is the default branch, it signals the TrunkCtrl workflow.
// Otherwise, it signals the BranchCtrl workflow.
func (a *Activities) SignalBranch(ctx context.Context, payload *defs.RepoIOSignalBranchCtrlPayload) error {
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

// SignalQueue signals a queue workflow for the given repository.
//
// It signals the QueueCtrl workflow with the repository, branch, and a serialized queue state.
func (a *Activities) SignalQueue(ctx context.Context, payload *defs.RepoIOSignalQueueCtrlPayload) error {
	args := make([]any, 0)
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("queue"),
		shared.WithWorkflowBlockID(payload.Repo.ID.String()),
		shared.WithWorkflowElement("branch"),
		shared.WithWorkflowElementID(payload.Branch),
	)

	queue := &QueueCtrlSerializedState{}

	args = append(args, payload.Repo, payload.Branch, queue)

	_, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(
			context.Background(),
			opts.ID,
			payload.Signal.String(),
			payload.Payload,
			opts,
			QueueCtrl,
			args...,
		)

	return err
}

// CloneBranch clones a repository branch to a temporary location.
//
// It retrieves the clone URL, including an OAuth token, using the RepoIO interface.
//
// If an error occurs while retrieving the clone URL, it is returned.
func (a *Activities) CloneBranch(ctx context.Context, payload *defs.RepoIOClonePayload) error {
	url, err := kernel.Instance().RepoIO(payload.Repo.Provider).TokenizedCloneURL(ctx, payload.Info)
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
//
// TODO: currently fetches the default branch, but should be made generic.
func (a Activities) FetchBranch(ctx context.Context, payload *defs.RepoIOClonePayload) error {
	cmd := exec.CommandContext(ctx, "git", "-C", payload.Path, "fetch", "origin", payload.Repo.DefaultBranch) //nolint

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

// RebaseAtCommit attempts to rebase the repository at the given commit.
//
// It returns the SHA and error message of the failed commit if the rebase fails.
//
// If the rebase is in progress, it returns an InProgress flag.
func (a *Activities) RebaseAtCommit(ctx context.Context, payload *defs.RepoIOClonePayload) (*defs.RepoIORebaseAtCommitResponse, error) {
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
				msg := fmt.Sprintf("error rebasing branch %s", payload.Branch)
				return &defs.RepoIORebaseAtCommitResponse{InProgress: true}, NewRebaseError("unknown", msg)
			default:
				return nil, err // retry
			}
		}

		return nil, err // retry
	}

	return nil, nil // rebase successful
}

// Push pushes the contents of the repository at the given path to the remote.
//
// If force is true, the push will be forced (--force).
func (a *Activities) Push(ctx context.Context, payload *defs.RepoIOPushBranchPayload) error {
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

// RemoveClonedAtPath removes the repository cloned at the given path.
//
// If the path does not exist, RemoveClonedAtPath returns a nil.
func (a *Activities) RemoveClonedAtPath(ctx context.Context, path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}

// GetByLogin retrieves a team user by their login ID.
//
// It calls the auth TeamUserIO to retrieve the user.
func (a *Activities) GetByLogin(ctx context.Context, id string) (*auth.TeamUser, error) {
	team_user, err := auth.TeamUserIO().GetByLogin(ctx, id)
	if err != nil {
		return nil, err
	}

	return team_user, nil
}

// SaveRepoEvent persists a repository event to the database.
//
// It converts the event to a db.Entity using the Flatten method and saves it to the database.
func (a *Activities) SaveRepoEvent(ctx context.Context, event RepoEvent) error {
	flat, err := event.Flatten()
	if err != nil {
		return err
	}

	return db.Save(flat)
}
