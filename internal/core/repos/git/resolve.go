package git

import (
	"context"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

// ResolveCommit resolves a revision string to its corresponding commit object.
//
// The same rules as ResolveRevision apply here.
func (r *Repository) ResolveCommit(ctx context.Context, revision string) (*object.Commit, error) {
	hash, err := r.ResolveRevision(ctx, revision)
	if err != nil {
		return nil, err
	}

	commit, err := r.cloned.CommitObject(*hash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
