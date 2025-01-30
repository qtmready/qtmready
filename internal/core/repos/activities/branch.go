package activities

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

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
		slog.Warn("clone: unable to get tokenized url", "error", err)
		return "", err
	}

	opts := &git.CloneOptions{
		CheckoutOptions: git.CheckoutOptions{
			Strategy:    git.CheckoutSafe,
			NotifyFlags: git.CheckoutNotifyAll,
		},
		CheckoutBranch: payload.Branch,
	}

	path := fmt.Sprintf("/tmp/%s", payload.Path)

	cloned, err := git.Clone(url, path, opts)
	if err != nil {
		slog.Warn("clone: failed", "error", err, "url", url, "path", path)
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
		slog.Warn("diff: unable to open repository", "error", err, "path", payload.Path)
		return nil, err
	}

	defer repo.Free()

	if err := a.refresh_remote(ctx, repo, payload.Base); err != nil {
		slog.Warn("diff: unable to refresh remote", "path", payload.Path, "error", err.Error())
		return nil, err
	}

	base, err := a.tree_from_branch(ctx, repo, payload.Base)
	if err != nil {
		slog.Warn("diff: unable to process base", "base", payload.Base, "error", err)
		return nil, err
	}

	defer base.Free()

	head, err := a.tree_from_sha(ctx, repo, payload.SHA)
	if err != nil {
		slog.Warn("diff: unable to process head", "head", payload.SHA, "error", err)
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

// Rebase performs a git rebase operation. Handles conflicts and returns result.
func (a *Branch) Rebase(ctx context.Context, payload *defs.RebasePayload) (*defs.RebaseResult, error) {
	result := defs.NewRebaseResult()

	repo, err := git.OpenRepository(payload.Path)
	if err != nil {
		a.report_rebase_error(ctx, result, "rebase: failed to open repository", err, payload.Rebase.Base, payload.Rebase.Head)
		return result, err
	}

	defer repo.Free()

	if err := a.refresh_remote(ctx, repo, payload.Rebase.Base); err != nil {
		a.report_rebase_error(
			ctx, result,
			"rebase: unable to refresh remote", err, payload.Rebase.Base, payload.Rebase.Head,
		)

		return result, nil
	}

	branch, upstream, err := a.get_annotated_commits(ctx, repo, payload.Rebase.Base, payload.Rebase.Head)
	if err != nil {
		a.report_rebase_error(
			ctx, result,
			"rebase: failed to get annotated commits", err, payload.Rebase.Base, payload.Rebase.Head,
		)

		return result, nil
	}

	defer branch.Free()
	defer upstream.Free()

	opts, err := git.DefaultRebaseOptions()
	if err != nil {
		a.report_rebase_error(
			ctx, result,
			"rebase: failed to get default rebase options", err, payload.Rebase.Base, payload.Rebase.Head,
		)

		result.Error = fmt.Sprintf("Failed to get default rebase options: %v", err)

		return result, nil
	}

	rebase, err := repo.InitRebase(branch, upstream, nil, &opts)
	if err != nil {
		a.report_rebase_error(
			ctx, result,
			"rebase: failed to initialize rebase", err, payload.Rebase.Base, payload.Rebase.Head,
		)

		result.Error = fmt.Sprintf("Failed to initialize rebase: %v", err)

		return result, nil
	}

	defer rebase.Free()

	result.TotalCommits = rebase.OperationCount()

	if err := a.rebase_each(ctx, repo, rebase, result); err != nil {
		a.report_rebase_error(
			ctx, result,
			"rebase: unable to rebase", err, payload.Rebase.Base, payload.Rebase.Head,
		)

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

	return result, nil
}

// NotifyLinesExceeded notifies on chat if lines exceed a limit.
func (a *Branch) NotifyLinesExceeded(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Diff]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyLinesExceed(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}

// NotifyMergeConflict notifies on chat if merge conflict message.
func (a *Branch) NotifyMergeConflict(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Merge]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyMergeConflict(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}

func (a *Branch) AnalyzeChange(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Diff]) error {
	return nil
}

// - Diff Helpers -
// diff_to_result converts a git.Diff to a DiffResult.
func (a *Branch) diff_to_result(_ context.Context, diff *git.Diff) (*eventsv1.Diff, error) {
	result := &eventsv1.Diff{Files: &eventsv1.DiffFiles{}, Lines: &eventsv1.DiffLines{}}
	deltas, err := diff.NumDeltas()

	if err != nil {
		slog.Warn("Failed to get number of deltas", "error", err)
		return nil, err
	}

	// use a sync.Map for concurrent safe updates
	var mutex sync.Mutex

	// use a go routine to handle diff deltas in parallel
	var wg sync.WaitGroup

	wg.Add(deltas)

	for idx := 0; idx < deltas; idx++ {
		go func(i int) {
			defer wg.Done()

			delta, _ := diff.Delta(i)

			mutex.Lock()
			defer mutex.Unlock()

			switch delta.Status { // nolint:exhaustive
			case git.DeltaAdded:
				result.Files.Added = append(result.Files.Added, delta.NewFile.Path)
			case git.DeltaDeleted:
				result.Files.Deleted = append(result.Files.Deleted, delta.OldFile.Path)
			case git.DeltaModified:
				result.Files.Modified = append(result.Files.Modified, delta.NewFile.Path)
			case git.DeltaRenamed:
				result.Files.Renamed = append(result.Files.Renamed, &eventsv1.RenamedFile{Old: delta.OldFile.Path, New: delta.NewFile.Path})
			}
		}(idx)
	}

	wg.Wait()

	stats, err := diff.Stats()
	if err != nil {
		return nil, err
	}

	defer func() { _ = stats.Free() }()

	result.Lines.Added = int32(stats.Insertions())  // nolint:gosec
	result.Lines.Removed = int32(stats.Deletions()) // nolint:gosec

	return result, nil
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

// - Rebase Helpers -

// get_annotated_commits retrieves annotated commits for the base and head of a rebase operation.
func (a *Branch) get_annotated_commits(
	ctx context.Context, repo *git.Repository, base string, head string,
) (*git.AnnotatedCommit, *git.AnnotatedCommit, error) {
	branch, err := a.annotated_commit_from_ref(ctx, repo, base)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get annotated commit from ref: %w", err)
	}

	upstream, err := a.annotated_commit_from_oid(ctx, repo, head)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get annotated commit from sha: %w", err)
	}

	return branch, upstream, nil
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

// rebase_each iterates over a rebase operation, processing each commit.
func (a *Branch) rebase_each(ctx context.Context, repo *git.Repository, rebase *git.Rebase, result *defs.RebaseResult) error {
	for {
		op, err := rebase.Next()
		if err != nil {
			if git.IsErrorCode(err, git.ErrorCodeIterOver) {
				result.SetStatusSuccess()

				break
			}

			a.rebase_abort(ctx, rebase)

			return err
		}

		if err := a.rebase_op(ctx, repo, rebase, op, result); err != nil {
			return err
		}
	}

	return nil
}

// rebase_op processes a rebase operation.
func (a *Branch) rebase_op(
	ctx context.Context, repo *git.Repository, rebase *git.Rebase, op *git.RebaseOperation, result *defs.RebaseResult,
) error {
	commit, err := repo.LookupCommit(op.Id)
	if err != nil {
		result.AddOperation(op.Type, defs.RebaseStatusFailure, "", commit.Message(), err)
		result.SetStatusFailure(err)
		a.rebase_abort(ctx, rebase)

		return err
	}

	defer commit.Free()

	slog.Debug("processing commit", "id", commit.Id().String())

	idx, err := repo.Index()
	if err != nil {
		result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), err)
		result.SetStatusFailure(err)

		a.rebase_abort(ctx, rebase)

		return err
	}
	defer idx.Free()

	conflicts, err := a.get_conflicts(ctx, idx)
	if err != nil {
		result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), err)
		result.SetStatusFailure(err)
		a.rebase_abort(ctx, rebase)

		return err
	} else if len(conflicts) > 0 {
		result.Conflicts = conflicts
		result.SetStatusConflicts()
		result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), nil)

		a.rebase_abort(ctx, rebase)

		return nil
	}

	err = rebase.Commit(commit.Id(), commit.Author(), commit.Committer(), commit.Message())
	if err != nil {
		result.AddOperation(op.Type, defs.RebaseStatusFailure, commit.Id().String(), commit.Message(), err)
		result.SetStatusFailure(err)

		a.rebase_abort(ctx, rebase)

		return err
	}

	slog.Debug("commit processed", "id", commit.Id().String())
	result.Head = commit.Id().String()
	result.AddOperation(op.Type, defs.RebaseStatusSuccess, commit.Id().String(), commit.Message(), nil)
	result.SetStatusSuccess()

	return nil
}

// rebase_abort aborts a git rebase operation if it's not nil.  Logs a warning if the abort fails.
func (a *Branch) rebase_abort(_ context.Context, rebase *git.Rebase) {
	slog.Debug("aborting rebase")

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

// report_rebase_error logs a rebase error and updates the rebase result.
func (a *Branch) report_rebase_error(_ context.Context, result *defs.RebaseResult, message string, err error, base string, head string) {
	slog.Warn(message, "error", err.Error(), "branch", base, "sha", head)

	result.Status = defs.RebaseStatusFailure
	result.Error = err.Error()
}

// - Git Helpers -

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
