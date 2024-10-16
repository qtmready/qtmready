// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
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

package code

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"

	"github.com/gocql/gocql"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared/queue"
)

type (
	// Activities defines an interface for repository-related actions.
	Activities struct{}

	// RepoCtrlQueryPayloadForBranchParent defines the payload for querying the RepoCtrl workflow for the parent event ID
	// of a branch.
	RepoCtrlQueryPayloadForBranchParent struct {
		Branch string
		Repo   *defs.Repo
	}

	// RepoCtrlQueryResultForBranchParent defines the result for querying the RepoCtrl workflow for the parent event ID.
	// If the Found flag is false, the parent event ID was not found.
	RepoCtrlQueryResultForBranchParent struct {
		ID    gocql.UUID
		Found bool
	}
)

// --- Git Operations ---

// CloneBranch clones a repository branch to a temporary location.
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
func (a Activities) FetchBranch(ctx context.Context, payload *defs.RepoIOClonePayload) error {
	cmd := exec.CommandContext(ctx, "git", "-C", payload.Path, "fetch", "origin", payload.Repo.DefaultBranch) //nolint

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

// RebaseAtCommit attempts to rebase the repository at the given commit.
func (a *Activities) RebaseAtCommit(ctx context.Context, payload *defs.RepoIOClonePayload) (*defs.RepoIORebaseAtCommitResponse, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", payload.Path, "rebase", payload.Rebase.After) // nolint

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
func (a *Activities) RemoveClonedAtPath(ctx context.Context, path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}

// --- DB Operations ---

// GetByLogin retrieves a team user by their login ID.
func (a *Activities) GetByLogin(ctx context.Context, id string) (*auth.TeamUser, error) {
	team_user, err := auth.TeamUserIO().GetByLogin(ctx, id)
	if err != nil {
		return nil, err
	}

	return team_user, nil
}

// SaveRepoEvent persists a repository event to the database.
func (a *Activities) SaveRepoEvent(ctx context.Context, event *defs.FlatEvent[defs.RepoProvider]) error {
	return db.CreateWithID(event, event.ID)
}

// --- Workflow Operations ---

// SignalBranch signals a branch workflow for the given repository.
func (a *Activities) SignalBranch(ctx context.Context, payload *defs.RepoIOSignalBranchCtrlPayload) error {
	if payload.Repo.DefaultBranch == payload.Branch {
		state := NewTrunkCtrlState(ctx, payload.Repo, payload.Info)
		_, err := queue.Core().SignalWithStartWorkflow(
			ctx,
			TrunkCtrlWorkflowOptions(payload.Repo.TeamID.String(), payload.Repo.Name, payload.Repo.ID),
			payload.Signal,
			payload.Payload,
			TrunkCtrl,
			state,
		)

		return err
	}

	ctx, state := NewBranchCtrlState(ctx, payload.Repo, payload.Info, payload.Branch)
	_, err := queue.Core().SignalWithStartWorkflow(
		ctx,
		BranchCtrlWorkflowOptions(payload.Repo.TeamID.String(), payload.Repo.Name, payload.Repo.ID, payload.Branch),
		payload.Signal,
		payload.Payload,
		BranchCtrl,
		state,
	)

	return err
}

// SignalQueue signals a queue workflow for the given repository.
func (a *Activities) SignalQueue(ctx context.Context, payload *defs.RepoIOSignalQueueCtrlPayload) error {
	args := make([]any, 0)

	_, err := queue.Core().SignalWithStartWorkflow(
		ctx,
		QueueCtrlWorkflowOptions(payload.Repo.TeamID.String(), payload.Repo.Name, payload.Repo.ID),
		payload.Signal,
		payload.Payload,
		QueueCtrl,
		append(args, payload.Repo, payload.Branch, &QueueCtrlSerializedState{})...,
	)

	return err
}

// QueryRepoCtrlForBranchParent queries the RepoCtrl workflow for the parent event ID of a branch.
func (a *Activities) QueryRepoCtrlForBranchParent(
	ctx context.Context, payload *RepoCtrlQueryPayloadForBranchParent,
) (*RepoCtrlQueryResultForBranchParent, error) {
	result := &RepoCtrlQueryResultForBranchParent{}
	opts := RepoCtrlWorkflowOptions(payload.Repo.TeamID.String(), payload.Repo.Name, payload.Repo.ID)

	slog.Info("querying repo ctrl for branch parent", "info", queue.Core().WorkflowID(opts)) // TODO: remove in production.

	return result, nil
}

// QueryRepoCtrlForBranchTriggers queries the RepoCtrl workflow for the branch triggers map.
func (a *Activities) QueryRepoCtrlForBranchTriggers(ctx context.Context, repo *defs.Repo) (BranchTriggers, error) {
	triggers := make(BranchTriggers)
	opts := RepoCtrlWorkflowOptions(repo.TeamID.String(), repo.Name, repo.ID)

	slog.Info("querying repo ctrl for branch parent", "info", queue.Core().WorkflowID(opts)) // TODO: remove in production.

	result, err := queue.Core().QueryWorkflow(ctx, opts, QueryRepoCtrlForBranchTriggers)
	if err != nil {
		return triggers, err
	}

	return triggers, result.Get(&triggers)
}
