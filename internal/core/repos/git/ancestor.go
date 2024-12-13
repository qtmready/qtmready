package git

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (r *Repository) Ancestor(a, b plumbing.Hash) (*object.Commit, error) {
	onto, err := r.cloned.CommitObject(a)
	if err != nil {
		return nil, err
	}

	upstream, err := r.cloned.CommitObject(b)
	if err != nil {
		return nil, err
	}

	ancestors, err := onto.MergeBase(upstream)
	if err != nil {
		return nil, err
	}

	if len(ancestors) == 0 {
		return nil, NewCompareError(r, OpAncestor, a.String(), b.String())
	}

	return ancestors[0], nil
}
