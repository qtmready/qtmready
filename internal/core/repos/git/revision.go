package git

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
)

// ResolveRevision resolves a revision string to its corresponding commit hash.
//
// Unsupported reference types (trees, annotated tags) are rejected.
//
// Supported revisions:
//   - HEAD, branches, tags,
//   - remote-tracking branches, HEAD~n, HEAD^,
//   - refspec selectors (e.g., HEAD^{/fix bug}),
//   - hash prefixes/full hashes.
func (r *Repository) ResolveRevision(ctx context.Context, revision string) (*plumbing.Hash, error) {
	if r.cloned == nil {
		if err := r.Open(); err != nil {
			return nil, err
		}
	}

	return r.cloned.ResolveRevision(plumbing.Revision(revision))
}
