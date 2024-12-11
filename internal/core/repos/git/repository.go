package git

import (
	"context"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos/cast"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	Repository struct {
		Entity *entities.Repo
		Branch string
		Path   string

		cloned *gogit.Repository
	}
)

func (r *Repository) Clone(ctx context.Context) error {
	if r.cloned != nil {
		return ErrRepoAlreadyExists
	}

	hook := cast.HookToProto(r.Entity.Hook)
	ref := plumbing.NewBranchReferenceName(fns.BranchNameToRef(r.Branch))

	if err := ref.Validate(); err != nil {
		return ErrInvalidBranch
	}

	url, err := kernel.Get().RepoHook(hook).TokenizedCloneUrl(ctx, r.Entity)
	if err != nil {
		return ErrTokenization
	}

	cloned, err := gogit.PlainCloneContext(ctx, r.Branch, false, &gogit.CloneOptions{
		URL:           url,
		ReferenceName: ref,
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return ErrClone
	}

	r.cloned = cloned

	return nil
}

func (r *Repository) Open() error {
	if r.cloned != nil {
		return nil
	}

	cloned, err := gogit.PlainOpen(r.Path)
	if err != nil {
		return ErrOpen
	}

	r.cloned = cloned

	return nil
}

func NewRepository(entity *entities.Repo, branch, path string) *Repository {
	return &Repository{
		Entity: entity,
		Branch: branch,
		Path:   path,
	}
}
