package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/diff"

	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// Diff returns a diff between two commits.
func (r *Repository) Diff(ctx context.Context, from, to string) (*eventsv1.Diff, error) {
	if r.cloned == nil {
		if err := r.Open(); err != nil {
			return nil, NewRepositoryError(r, OpOpen).Wrap(err)
		}
	}

	from_commit, err := r.ResolveCommit(ctx, from)
	if err != nil {
		return nil, NewResolveError(r, OpResolveCommit, from).Wrap(err)
	}

	to_commit, err := r.ResolveCommit(ctx, to)
	if err != nil {
		return nil, NewResolveError(r, OpResolveCommit, to).Wrap(err)
	}

	patch, err := from_commit.Patch(to_commit)
	if err != nil {
		return nil, NewCompareError(r, OpDiff, from, to).Wrap(err)
	}

	files, lines := PatchToDiff(patch)

	commits := &eventsv1.DiffCommits{
		Base: from_commit.Hash.String(),
		Head: to_commit.Hash.String(),
	}

	builder := strings.Builder{}
	if patch != nil {
		builder.WriteString(patch.String())

		ancestor, err := r.Ancestor(from_commit.Hash, to_commit.Hash)
		if err != nil {
			if _, ok := err.(*CompareError); !ok {
				err = NewCompareError(r, OpAncestor, from_commit.Hash.String(), to_commit.Hash.String()).Wrap(err)
			}

			return nil, err
		}

		if ancestor != nil {
			commits.ConflictAt = ancestor.Hash.String()
		}
	}

	stats := patch.Stats()

	for _, stat := range stats {
		lines.Added += int32(stat.Addition)
		lines.Removed += int32(stat.Deletion)
	}

	has_conflict := !(commits.ConflictAt != "")

	return &eventsv1.Diff{
		Files:       files,
		Lines:       lines,
		Commits:     commits,
		Patch:       builder.String(),
		HasConflict: has_conflict,
	}, nil
}

// PatchToDiff converts a git patch to a DiffFiles and DiffLines struct.
func PatchToDiff(patch diff.Patch) (*eventsv1.DiffFiles, *eventsv1.DiffLines) {
	files := &eventsv1.DiffFiles{
		Added:    make([]string, 0),
		Deleted:  make([]string, 0),
		Modified: make([]string, 0),
		Renamed:  make([]string, 0),
	}

	lines := &eventsv1.DiffLines{}

	if patch == nil {
		return files, lines
	}

	for _, fp := range patch.FilePatches() {
		from, to := fp.Files()

		if from == nil { // nolint: gocritic
			files.Added = append(files.Added, to.Path())
		} else if to == nil {
			files.Deleted = append(files.Deleted, from.Path())
		} else if from.Path() != to.Path() {
			files.Renamed = append(files.Renamed, fmt.Sprintf("%s => %s", from.Path(), to.Path()))
		} else {
			files.Modified = append(files.Modified, from.Path())
		}
	}

	return files, lines
}
