package activities

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	git "github.com/jeffwelling/git2go/v37"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Branch struct{}
)

// Clone clones a repo to a temp path, fetching a specified branch.
func (a *Branch) Clone(ctx context.Context, payload *defs.ClonePayload) (string, error) {
	url, err := kernel.Get().RepoHook(payload.Hook).TokenizedCloneUrl(ctx, payload.Repo)
	if err != nil {
		slog.Warn("clone: unable to get tokenized url ...", "error", err)
		return "", err
	}

	opts := &git.CloneOptions{
		CheckoutOptions: git.CheckoutOptions{
			Strategy:    git.CheckoutSafe,
			NotifyFlags: git.CheckoutNotifyAll,
		},
		CheckoutBranch: payload.Branch,
	}

	cloned, err := git.Clone(url, fmt.Sprintf("/tmp/%s", payload.Path), opts)
	if err != nil {
		slog.Warn("clone: failed ...", "error", err, "url", url, "path", fmt.Sprintf("/tmp/%s", payload.Path))

		return "", err
	}

	defer cloned.Free()

	return cloned.Workdir(), nil
}

// RemoveDir removes a directory and handles potential errors.
func (a *Branch) RemoveDir(ctx context.Context, path string) error {
	slog.Debug("removing directory", "path", path)

	if err := os.RemoveAll(path); err != nil {
		slog.Warn("Failed to remove directory", "error", err, "path", path)
	}

	return nil
}

// Diff computes the diff between two commits using git2go.
func (a *Branch) Diff(ctx context.Context, payload *defs.DiffPayload) (*eventsv1.Diff, error) {
	repo, err := git.OpenRepository(payload.Path)
	if err != nil {
		slog.Warn("Failed to open repository", "error", err, "path", payload.Path)
		return nil, err
	}

	defer repo.Free()

	if err := a.refresh_remote(ctx, repo, payload.Base); err != nil {
		slog.Warn("diff: unable to refresh remote", "path", payload.Path, "error", err.Error())
		return nil, err
	}

	base, err := a.tree_from_branch(ctx, repo, payload.Base)
	if err != nil {
		slog.Warn("diff: unable to process base", "base", base)
		return nil, err
	}

	defer base.Free()

	head, err := a.tree_from_sha(ctx, repo, payload.SHA)
	if err != nil {
		slog.Warn("diff: unable to process head", "head", payload.SHA)
		return nil, err
	}

	defer head.Free()

	opts, _ := git.DefaultDiffOptions()

	diff, err := repo.DiffTreeToTree(base, head, &opts)
	if err != nil {
		slog.Warn("Failed to create diff", "error", err, "base", base, "head", head)
		return nil, err
	}

	defer func() { _ = diff.Free() }()

	return a.diff_to_result(ctx, diff)
}

// Rebase performs a git rebase operation.  Handles conflicts and returns result.
func (a *Branch) Rebase(ctx context.Context, payload *defs.RebasePayload) (*defs.RebaseResult, error) {
	result := defs.NewRebaseResult()

	repo, err := git.OpenRepository(payload.Path)
	if err != nil {
		slog.Warn("rebase: failed to open repository", "error", err, "path", payload.Path)
		return result, err
	}

	defer repo.Free()

	if err := a.refresh_remote(ctx, repo, payload.Rebase.Base); err != nil {
		slog.Warn(
			"rebase: unable to refresh remote",
			"error", err.Error(),
			"branch", payload.Rebase.Base, "sha", payload.Rebase.Head)

		result.Status = defs.RebaseStatusFailure
		result.Error = err.Error()

		return result, nil
	}

	branch, err := a.annotated_commit_from_ref(ctx, repo, payload.Rebase.Base)
	if err != nil {
		slog.Warn(
			"rebase: failed to get annotated commit from ref",
			"error", err.Error(),
			"branch", payload.Rebase.Base, "sha", payload.Rebase.Head,
		)

		result.Status = defs.RebaseStatusFailure
		result.Error = err.Error()

		return result, nil
	}

	defer branch.Free()

	upstream, err := a.annotated_commit_from_oid(ctx, repo, payload.Rebase.Head)
	if err != nil {
		slog.Warn(
			"rebase: failed to get annotated commit from sha",
			"error", err.Error(),
			"branch", payload.Rebase.Base, "sha", payload.Rebase.Head,
		)

		result.Status = defs.RebaseStatusFailure
		result.Error = fmt.Sprintf("Failed to get annotated commit from OID: %v", err)

		return result, nil
	}

	defer upstream.Free()

	opts, err := git.DefaultRebaseOptions()
	if err != nil {
		slog.Warn(
			"rebase: to get default rebase options",
			"error", err.Error(),
			"branch", payload.Rebase.Base, "sha", payload.Rebase.Head,
		)

		result.Status = defs.RebaseStatusFailure
		result.Error = fmt.Sprintf("Failed to get default rebase options: %v", err)

		return result, nil
	}

	rebase, err := repo.InitRebase(branch, upstream, nil, &opts)
	if err != nil {
		slog.Warn(
			"rebase: failed to initialize rebase",
			"error", err.Error(),
			"branch", payload.Rebase.Base, "sha", payload.Rebase.Head,
		)

		result.Status = defs.RebaseStatusFailure
		result.Error = fmt.Sprintf("Failed to initialize rebase: %v", err)

		return result, nil
	}

	defer a.rebase_abort(ctx, rebase)
	defer rebase.Free()

	result.TotalCommits = rebase.OperationCount()

	if err := a.rebase_each(ctx, repo, rebase, result); err != nil {
		slog.Warn(
			"rebase: unable to rebase",
			"error", err.Error(),
			"branch", payload.Rebase.Base, "sha", payload.Rebase.Head,
		)

		result.Status = defs.RebaseStatusFailure
		result.Error = err.Error()

		return result, nil
	}

	if err := rebase.Finish(); err != nil {
		slog.Warn(
			"rebase: unable to finish",
			"error", err.Error(),
			"branch", payload.Rebase.Base, "sha", payload.Rebase.Head,
		)

		return result, nil
	}

	rebase = nil

	return result, nil
}

func (a *Branch) rebase_each(ctx context.Context, repo *git.Repository, rebase *git.Rebase, result *defs.RebaseResult) error {
	for {
		op, err := rebase.Next()
		if err != nil {
			if git.IsErrorCode(err, git.ErrorCodeIterOver) {
				result.SetStatusSuccess()

				break
			}

			return err
		}

		commit, err := repo.LookupCommit(op.Id)
		if err != nil {
			result.AddOperation(op.Type, defs.RebaseStatusFailure, "", commit.Message(), err)
			return err
		}

		defer commit.Free()

		slog.Debug("processing commit", "id", commit.Id().String())

		idx, err := repo.Index()
		if err != nil {
			result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), err)

			return err
		}

		defer idx.Free()

		if conflicts, err := a.get_conflicts(ctx, idx); err != nil {
			result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), err)
			return err
		} else if len(conflicts) > 0 {
			result.Conflicts = conflicts
			result.SetStatusConflicts()
			result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), nil)

			continue
		}

		err = rebase.Commit(commit.Id(), commit.Author(), commit.Committer(), commit.Message())
		if err != nil {
			result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), err)

			return err
		}

		slog.Debug("commit processed", "id", commit.Id().String())
		result.Head = commit.Id().String()
		result.AddOperation(op.Type, defs.RebaseStatusSuccess, commit.Id().String(), commit.Message(), nil)
		result.SetStatusSuccess()
	}

	return nil
}

// rebase_abort aborts a git rebase operation if it's not nil.  Logs a warning if the abort fails.
func (a *Branch) rebase_abort(_ context.Context, rebase *git.Rebase) {
	if rebase != nil {
		if err := rebase.Abort(); err != nil {
			slog.Warn("rebase: unable to abort!", "error", err.Error())
		}
	}
}

// get_conflicts retrieves conflict information from a git index. Returns an empty slice if no conflicts are found.
func (a *Branch) get_conflicts(_ context.Context, idx *git.Index) ([]string, error) {
	conflicts := make([]string, 0)

	if idx == nil {
		return conflicts, nil
	}

	if !idx.HasConflicts() {
		return conflicts, nil
	}

	iter, err := idx.ConflictIterator()
	if err != nil {
		slog.Warn("Failed to create conflict iterator", "error", err)
		return conflicts, fmt.Errorf("failed to create conflict iterator: %w", err)
	}

	defer iter.Free()

	for {
		entry, err := iter.Next()
		if err != nil {
			if git.IsErrorCode(err, git.ErrorCodeIterOver) {
				break
			}

			slog.Warn("Failed to get next conflict entry", "error", err)

			return conflicts, fmt.Errorf("failed to get next conflict entry: %w", err)
		}

		conflicts = append(conflicts, entry.Ancestor.Path)
	}

	return conflicts, nil
}

// annotated_commit_from_ref retrieves an annotated commit from a ref.
func (a *Branch) annotated_commit_from_ref(_ context.Context, repo *git.Repository, branch string) (*git.AnnotatedCommit, error) {
	ref, err := repo.References.Lookup(fns.BranchNameToRef(branch))
	if err != nil {
		slog.Warn("Failed to lookup ref", "error", err, "branch", branch)
		return nil, err
	}

	defer ref.Free()

	commit, err := repo.LookupAnnotatedCommit(ref.Target())
	if err != nil {
		slog.Warn("Failed to lookup base commit", "error", err, "target", ref.Target())
		return nil, err
	}

	return commit, nil
}

// annotated_commit_from_oid retrieves an annotated commit from an OID.
func (a *Branch) annotated_commit_from_oid(_ context.Context, repo *git.Repository, sha string) (*git.AnnotatedCommit, error) {
	id, err := git.NewOid(sha)
	if err != nil {
		slog.Warn("Invalid head SHA", "error", err, "sha", sha)
		return nil, err
	}

	commit, err := repo.LookupAnnotatedCommit(id)
	if err != nil {
		slog.Warn("Failed to lookup head commit", "error", err, "id", id)
		return nil, err
	}

	return commit, nil
}

// refresh_remote fetches a branch from the "origin" remote.
func (a *Branch) refresh_remote(_ context.Context, repo *git.Repository, branch string) error {
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	if err := remote.Fetch([]string{fns.BranchNameToRef(branch)}, &git.FetchOptions{}, ""); err != nil {
		return err
	}

	ref, err := repo.References.Lookup(fns.BranchNameToRemoteRef("origin", branch))
	if err != nil {
		return err
	}

	defer ref.Free()

	_, err = repo.References.Create(fns.BranchNameToRef(branch), ref.Target(), true, "")
	if err != nil {
		return err
	}

	return nil
}

// tree_from_branch gets the tree from a branch ref.
func (a *Branch) tree_from_branch(_ context.Context, repo *git.Repository, branch string) (*git.Tree, error) {
	ref, err := repo.References.Lookup(fns.BranchNameToRef(branch))
	if err != nil {
		slog.Warn("Failed to lookup ref", "error", err, "branch", branch)
		return nil, err
	}

	defer ref.Free()

	commit, err := repo.LookupCommit(ref.Target())
	if err != nil {
		slog.Warn("Failed to lookup commit", "error", err, "target", ref.Target())
		return nil, err
	}

	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		slog.Warn("Failed to lookup tree", "error", err)
		return nil, err
	}

	return tree, nil
}

// tree_from_sha gets the tree from a commit SHA.
func (a *Branch) tree_from_sha(_ context.Context, repo *git.Repository, sha string) (*git.Tree, error) {
	oid, err := git.NewOid(sha)
	if err != nil {
		slog.Warn("Invalid SHA", "error", err, "sha", sha)
		return nil, err
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		slog.Warn("Failed to lookup commit", "error", err, "oid", oid)
		return nil, err
	}

	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		slog.Warn("Failed to lookup tree", "error", err)
		return nil, err
	}

	return tree, nil
}

// diff_to_result converts a git.Diff to a DiffResult.
func (a *Branch) diff_to_result(_ context.Context, diff *git.Diff) (*eventsv1.Diff, error) {
	result := &eventsv1.Diff{Files: &eventsv1.DiffFiles{}, Lines: &eventsv1.DiffLines{}}

	deltas, err := diff.NumDeltas()
	if err != nil {
		slog.Warn("Failed to get number of deltas", "error", err)
		return nil, err
	}

	for idx := 0; idx < deltas; idx++ {
		delta, _ := diff.Delta(idx)

		switch delta.Status { // nolint:exhaustive
		default:
		case git.DeltaAdded:
			result.Files.Added = append(result.Files.Added, delta.NewFile.Path)
		case git.DeltaDeleted:
			result.Files.Deleted = append(result.Files.Deleted, delta.OldFile.Path)
		case git.DeltaModified:
			result.Files.Modified = append(result.Files.Modified, delta.NewFile.Path)
		case git.DeltaRenamed:
			result.Files.Renamed = append(result.Files.Renamed, delta.NewFile.Path)
		}
	}

	stats, err := diff.Stats()
	if err != nil {
		return nil, err
	}

	defer func() { _ = stats.Free() }()

	result.Lines.Added = int32(stats.Insertions())  // nolint:gosec
	result.Lines.Removed = int32(stats.Deletions()) // nolint:gosec

	return result, nil
}

// ExceedLines notifies on chat if lines exceed a limit.
func (a *Branch) ExceedLines(ctx context.Context, event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) error {
	if err := kernel.Get().ChatHook(event.Context.Hook).NotifyLinesExceed(ctx, event); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}

// MergeConflict notifies on chat if merge conflict message.
func (a *Branch) MergeConflict(ctx context.Context, event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) error {
	if err := kernel.Get().ChatHook(event.Context.Hook).NotifyMergeConflict(ctx, event); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}
