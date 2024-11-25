package activities

import (
	"context"
	"fmt"
	"log/slog"

	git "github.com/jeffwelling/git2go/v37"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/db/entities"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	ClonePayload struct {
		Repo   *entities.Repo    `json:"repo"`
		Hook   eventsv1.RepoHook `json:"hook"`
		Branch string            `json:"branch"`
		Path   string            `json:"path"`
		SHA    string            `json:"sha"`
	}

	DiffPayload struct {
		Path string `json:"path"`
		Base string `json:"base"`
		SHA  string `json:"sha"`
	}

	Branch struct{}
)

// Clone clones the repository at the specified branch using a temporary path.  It retrieves a tokenized clone URL,
// clones the repository using `git2go`, fetches the specified branch, and returns the working directory path.
func (a *Branch) Clone(ctx context.Context, payload *ClonePayload) (string, error) {
	url, err := kernel.Get().RepoHook(payload.Hook).TokenizedCloneUrl(ctx, payload.Repo)
	if err != nil {
		slog.Error("Failed to get tokenized clone URL", "error", err) // Log the error
		return "", err
	}

	slog.Info("cloning ...", "url", url)

	opts := &git.CloneOptions{
		CheckoutOptions: git.CheckoutOptions{
			Strategy:    git.CheckoutSafe,
			NotifyFlags: git.CheckoutNotifyAll,
		},
		CheckoutBranch: payload.Branch,
	}

	cloned, err := git.Clone(url, fmt.Sprintf("/tmp/%s", payload.Path), opts)
	if err != nil {
		slog.Error("Failed to clone repository", "error", err, "url", url, "path", fmt.Sprintf("/tmp/%s", payload.Path))
		return "", err
	}

	defer cloned.Free()

	slog.Info("cloned successfully", "repo", payload.Repo.Url, "cloned", cloned.Workdir())

	return cloned.Workdir(), nil
}

// Diff retrieves the diff between two commits.  Given a repository path, base branch, and SHA, it opens the repo,
// fetches the base branch, resolves commits by SHA, and computes the diff between their trees using `git2go`. The
// resulting diff is currently returned unprocessed.
func (a *Branch) Diff(ctx context.Context, payload *DiffPayload) (string, error) {
	repo, err := git.OpenRepository(payload.Path)
	if err != nil {
		slog.Error("Failed to open repository", "error", err, "path", payload.Path)
		return "", err
	}

	defer repo.Free()

	if err := a.refresh_remote(ctx, repo, payload.Base); err != nil {
		return "", err
	}

	base, err := a.tree_from_branch(ctx, repo, payload.Base)
	if err != nil {
		slog.Error("unable to process base", "base", base)
		return "", err
	}

	defer base.Free()

	head, err := a.tree_from_sha(ctx, repo, payload.SHA)
	if err != nil {
		slog.Error("unable to process head", "head", head)
		return "", err
	}

	defer head.Free()

	diffopts, _ := git.DefaultDiffOptions()

	diff, err := repo.DiffTreeToTree(base, head, &diffopts)
	if err != nil {
		slog.Error("Failed to create diff", "error", err)
		return "", err
	}

	defer func() { _ = diff.Free() }()

	// Process the diff (diff.NumDeltas(), diff.Delta(i), etc.)  This will depend on your needs.  For now, return an empty string.
	//Example:  You might want to convert it to a unified diff string.

	return "", nil
}

// refresh_remote fetches the latest changes from the remote "origin" for the given branch.
// It looks up the remote, fetches the branch, and updates the local branch reference.
func (a *Branch) refresh_remote(_ context.Context, repo *git.Repository, branch string) error {
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		slog.Error("failed to set remote", "remote", "origin", "error", err.Error())
		return err
	}

	if err := remote.Fetch([]string{fns.BranchNameToRef(branch)}, &git.FetchOptions{}, ""); err != nil {
		slog.Error("unable to fetch from remote", "error", err.Error())
		return err
	}

	ref, err := repo.References.Lookup(fns.BranchNameToRemoteRef("origin", branch))
	if err != nil {
		slog.Error("unable to lookup remote ref", "error", err.Error())
		return err
	}

	defer ref.Free()

	_, err = repo.References.Create(fns.BranchNameToRef(branch), ref.Target(), true, "")
	if err != nil {
		slog.Error("unable to create ref", "error", err.Error())
		return err
	}

	return nil
}

// tree_from_branch retrieves the tree object associated with the given branch.
// It looks up the branch reference, retrieves the corresponding commit, and returns the commit's tree.
func (a *Branch) tree_from_branch(_ context.Context, repo *git.Repository, branch string) (*git.Tree, error) {
	ref, err := repo.References.Lookup(fns.BranchNameToRef(branch))
	if err != nil {
		slog.Error("Failed to lookup ref", "error", err, "branch", branch)
		return nil, err
	}

	defer ref.Free()

	commit, err := repo.LookupCommit(ref.Target())
	if err != nil {
		slog.Error("Failed to lookup commit", "error", err, "target", ref.Target())
		return nil, err
	}

	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		slog.Error("Failed to lookup tree", "error", err)
		return nil, err
	}

	return tree, nil
}

// tree_from_sha retrieves the tree object associated with the given SHA.
// It looks up the commit by SHA and returns the commit's tree.
func (a *Branch) tree_from_sha(_ context.Context, repo *git.Repository, sha string) (*git.Tree, error) {
	oid, err := git.NewOid(sha)
	if err != nil {
		slog.Error("Invalid SHA", "error", err, "sha", sha)
		return nil, err
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		slog.Error("Failed to lookup commit", "error", err, "oid", oid)
		return nil, err
	}

	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		slog.Error("Failed to lookup tree", "error", err)
		return nil, err
	}

	return tree, nil
}
