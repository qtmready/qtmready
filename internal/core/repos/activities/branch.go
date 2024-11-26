package activities

import (
	"context"
	"fmt"

	git "github.com/jeffwelling/git2go/v37"
	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/fns"
)

type (
	Branch struct{}
)

// Clone clones the repository at the specified branch using a temporary path.  It retrieves a tokenized clone URL,
// clones the repository using git2go, fetches the specified branch, and returns the working directory path.
func (a *Branch) Clone(ctx context.Context, payload *defs.ClonePayload) (string, error) {
	logger := activity.GetLogger(ctx)

	url, err := kernel.Get().RepoHook(payload.Hook).TokenizedCloneUrl(ctx, payload.Repo)
	if err != nil {
		logger.Error("Failed to get tokenized clone URL", "error", err) // Log the error
		return "", err
	}

	logger.Info("cloning ...", "url", url)

	opts := &git.CloneOptions{
		CheckoutOptions: git.CheckoutOptions{
			Strategy:    git.CheckoutSafe,
			NotifyFlags: git.CheckoutNotifyAll,
		},
		CheckoutBranch: payload.Branch,
	}

	cloned, err := git.Clone(url, fmt.Sprintf("/tmp/%s", payload.Path), opts)
	if err != nil {
		logger.Error("Failed to clone repository", "error", err, "url", url, "path", fmt.Sprintf("/tmp/%s", payload.Path))
		return "", err
	}

	defer cloned.Free()

	logger.Info("cloned successfully", "repo", payload.Repo.Url, "cloned", cloned.Workdir())

	return cloned.Workdir(), nil
}

// Diff retrieves the diff between two commits.  Given a repository path, base branch, and SHA, it opens the repo,
// fetches the base branch, resolves commits by SHA, and computes the diff between their trees using `git2go`. The
// resulting diff is currently returned unprocessed.
func (a *Branch) Diff(ctx context.Context, payload *defs.DiffPayload) (*defs.DiffResult, error) {
	logger := activity.GetLogger(ctx)

	repo, err := git.OpenRepository(payload.Path)
	if err != nil {
		logger.Error("Failed to open repository", "error", err, "path", payload.Path)
		return nil, err
	}

	defer repo.Free()

	if err := a.refresh_remote(ctx, repo, payload.Base); err != nil {
		return nil, err
	}

	base, err := a.tree_from_branch(ctx, repo, payload.Base)
	if err != nil {
		logger.Error("unable to process base", "base", base)
		return nil, err
	}

	defer base.Free()

	head, err := a.tree_from_sha(ctx, repo, payload.SHA)
	if err != nil {
		logger.Error("unable to process head", "head", head)
		return nil, err
	}

	defer head.Free()

	diffopts, _ := git.DefaultDiffOptions()

	diff, err := repo.DiffTreeToTree(base, head, &diffopts)
	if err != nil {
		logger.Error("Failed to create diff", "error", err)
		return nil, err
	}

	defer func() { _ = diff.Free() }()

	return a.diff_to_result(ctx, diff)
}

// refresh_remote fetches the latest changes from the remote "origin" for the given branch.
// It looks up the remote, fetches the branch, and updates the local branch reference.
func (a *Branch) refresh_remote(ctx context.Context, repo *git.Repository, branch string) error {
	logger := activity.GetLogger(ctx)

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		logger.Error("failed to set remote", "remote", "origin", "error", err.Error())
		return err
	}

	if err := remote.Fetch([]string{fns.BranchNameToRef(branch)}, &git.FetchOptions{}, ""); err != nil {
		logger.Error("unable to fetch from remote", "error", err.Error())
		return err
	}

	ref, err := repo.References.Lookup(fns.BranchNameToRemoteRef("origin", branch))
	if err != nil {
		logger.Error("unable to lookup remote ref", "error", err.Error())
		return err
	}

	defer ref.Free()

	_, err = repo.References.Create(fns.BranchNameToRef(branch), ref.Target(), true, "")
	if err != nil {
		logger.Error("unable to create ref", "error", err.Error())
		return err
	}

	return nil
}

// tree_from_branch retrieves the tree object associated with the given branch.
// It looks up the branch reference, retrieves the corresponding commit, and returns the commit's tree.
func (a *Branch) tree_from_branch(ctx context.Context, repo *git.Repository, branch string) (*git.Tree, error) {
	logger := activity.GetLogger(ctx)

	ref, err := repo.References.Lookup(fns.BranchNameToRef(branch))
	if err != nil {
		logger.Error("Failed to lookup ref", "error", err, "branch", branch)
		return nil, err
	}

	defer ref.Free()

	commit, err := repo.LookupCommit(ref.Target())
	if err != nil {
		logger.Error("Failed to lookup commit", "error", err, "target", ref.Target())
		return nil, err
	}

	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		logger.Error("Failed to lookup tree", "error", err)
		return nil, err
	}

	return tree, nil
}

// tree_from_sha retrieves the tree object associated with the given SHA.
// It looks up the commit by SHA and returns the commit's tree.
func (a *Branch) tree_from_sha(ctx context.Context, repo *git.Repository, sha string) (*git.Tree, error) {
	logger := activity.GetLogger(ctx)

	oid, err := git.NewOid(sha)
	if err != nil {
		logger.Error("Invalid SHA", "error", err, "sha", sha)
		return nil, err
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		logger.Error("Failed to lookup commit", "error", err, "oid", oid)
		return nil, err
	}

	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		logger.Error("Failed to lookup tree", "error", err)
		return nil, err
	}

	return tree, nil
}

// diff_to_result converts a git.Diff into a DiffResult.
// It iterates through the deltas in the diff, categorizing files based on their status
// (added, deleted, modified, etc.). It also calculates the total number of lines added
// and removed using the diff statistics.
func (a *Branch) diff_to_result(ctx context.Context, diff *git.Diff) (*defs.DiffResult, error) {
	logger := activity.GetLogger(ctx)
	result := &defs.DiffResult{}

	deltas, err := diff.NumDeltas()
	if err != nil {
		logger.Error("Failed to get number of deltas", "error", err)
		return nil, err
	}

	for idx := 0; idx < deltas; idx++ {
		delta, _ := diff.Delta(idx)

		switch delta.Status {
		case git.DeltaAdded:
			result.Files.Added = append(result.Files.Added, delta.NewFile.Path)
		case git.DeltaDeleted:
			result.Files.Deleted = append(result.Files.Deleted, delta.OldFile.Path)
		case git.DeltaModified:
			result.Files.Modified = append(result.Files.Modified, delta.NewFile.Path)
		case git.DeltaRenamed:
			result.Files.Renamed = append(result.Files.Renamed, delta.NewFile.Path)
		case git.DeltaCopied:
			result.Files.Copied = append(result.Files.Copied, delta.NewFile.Path)
		case git.DeltaTypeChange:
			result.Files.TypeChange = append(result.Files.TypeChange, delta.NewFile.Path)
		case git.DeltaUnreadable:
			result.Files.Unreadable = append(result.Files.Unreadable, delta.NewFile.Path)
		case git.DeltaIgnored:
			result.Files.Ignored = append(result.Files.Ignored, delta.NewFile.Path)
		case git.DeltaUntracked:
			result.Files.Untracked = append(result.Files.Untracked, delta.NewFile.Path)
		case git.DeltaConflicted:
			result.Files.Conflicted = append(result.Files.Conflicted, delta.NewFile.Path)
		case git.DeltaUnmodified:
		}
	}

	stats, err := diff.Stats()
	if err != nil {
		return nil, err
	}

	defer func() { _ = stats.Free() }()

	result.Lines.Added = stats.Insertions()
	result.Lines.Removed = stats.Deletions()

	return result, nil
}
