package git

import (
	"context"
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (r *Repository) CherryPick(ctx context.Context, branch, hash string) (*object.Commit, error) {
	if r.cloned == nil {
		if err := r.Open(); err != nil {
			return nil, NewRepositoryError(r, OpOpen).Wrap(err)
		}
	}

	pick, err := r.ResolveCommit(ctx, hash)
	if err != nil {
		return nil, NewResolveError(r, OpResolveCommit, hash).Wrap(err)
	}

	worktree, err := r.cloned.Worktree()
	if err != nil {
		return nil, NewCherryPickError(r, "worktree", hash).Wrap(err)
	}

	err = worktree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Create: false,
	})

	if err != nil {
		return nil, NewCherryPickError(r, "checkout", hash).Wrap(err)
	}

	commit, err := worktree.Commit(pick.Message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  pick.Author.Name,
			Email: pick.Author.Email,
			When:  pick.Author.When,
		},
	})

	if err != nil {
		return nil, NewCherryPickError(r, "commit", hash).Wrap(err)
	}

	cp, err := r.cloned.CommitObject(commit)
	if err != nil {
		return nil, NewCherryPickError(r, "commit_object", hash).Wrap(err)
	}

	if err := worktree.Checkout(&gogit.CheckoutOptions{
		Hash:   cp.Hash,
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	}); err != nil {
		return nil, NewCherryPickError(r, "checkout_post_cherrypick", hash).Wrap(err)
	}

	return cp, nil
}
