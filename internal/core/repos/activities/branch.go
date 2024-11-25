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

// Clone clones the repository at the specified branch using a temporary path.
func (a *Branch) Clone(ctx context.Context, payload *ClonePayload) (string, error) {
	url, err := kernel.Get().RepoHook(payload.Hook).TokenizedCloneUrl(ctx, payload.Repo)
	if err != nil {
		slog.Error("Failed to get tokenized clone URL", "error", err) // Log the error
		return "", err
	}

	slog.Info("cloning ...", "url", url)

	copts := &git.CloneOptions{
		CheckoutOptions: git.CheckoutOptions{
			Strategy:    git.CheckoutSafe,
			NotifyFlags: git.CheckoutNotifyAll,
		},
		CheckoutBranch: payload.Branch,
	}

	cloned, err := git.Clone(url, fmt.Sprintf("/tmp/%s", payload.Path), copts)
	if err != nil {
		slog.Error("Failed to clone repository", "error", err, "url", url, "path", fmt.Sprintf("/tmp/%s", payload.Path))
		return "", err
	}

	defer cloned.Free()

	remote, err := cloned.Remotes.Lookup("origin")
	if err != nil {
		slog.Error("failed to set remote", "remote", "origin", "error", err.Error())
		return "", err
	}

	if err := remote.Fetch([]string{fns.BranchNameToRef(payload.Repo.DefaultBranch)}, &git.FetchOptions{}, ""); err != nil {
		slog.Error("unable to fetch from remote", "error", err.Error())
		return "", err
	}

	slog.Info("cloned successfully", "repo", payload.Repo.Url, "cloned", cloned.Workdir())

	return cloned.Workdir(), nil
}

func (a *Branch) Diff(ctx context.Context, payload *DiffPayload) (string, error) {
	repo, err := git.OpenRepository(payload.Path)
	if err != nil {
		slog.Error("Failed to open repository", "error", err, "path", payload.Path)
		return "", err
	}

	defer repo.Free()

	ls, _ := repo.Remotes.List()
	slog.Info("all remotes", "list", ls)

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		slog.Error("failed to set remote", "remote", "origin", "error", err.Error())
		return "", err
	}

	if err := remote.Fetch([]string{fns.BranchNameToRef(payload.Base)}, &git.FetchOptions{}, ""); err != nil {
		slog.Error("unable to fetch from remote", "error", err.Error())
		return "", err
	}

	baseref, err := repo.References.Lookup(fns.BranchNameToRef(payload.Base))
	if err != nil {
		slog.Error("Failed to lookup base ref", "error", err, "branch", payload.Base)
		return "", err
	}

	defer baseref.Free()

	basecommit, err := repo.LookupCommit(baseref.Target())
	if err != nil {
		slog.Error("Failed to lookup base commit", "error", err, "target", baseref.Target())
		return "", err
	}

	defer basecommit.Free()

	basetree, err := basecommit.Tree()
	if err != nil {
		return "", err
	}

	defer basetree.Free()

	sha, err := git.NewOid(payload.SHA)
	if err != nil {
		slog.Error("Invalid SHA", "error", err, "sha", payload.SHA)
		return "", err
	}

	headcommit, err := repo.LookupCommit(sha)
	if err != nil {
		slog.Error("Failed to lookup head commit", "error", err, "sha", sha)
		return "", err
	}

	defer headcommit.Free()

	head, err := headcommit.Tree()
	if err != nil {
		return "", err
	}
	defer head.Free()

	diffopts, _ := git.DefaultDiffOptions()

	diff, err := repo.DiffTreeToTree(basetree, head, &diffopts)
	if err != nil {
		slog.Error("Failed to create diff", "error", err)
		return "", err
	}

	defer func() { _ = diff.Free() }()

	// Process the diff (diff.NumDeltas(), diff.Delta(i), etc.)  This will depend on your needs.  For now, return an empty string.
	//Example:  You might want to convert it to a unified diff string.

	return "", nil
}
