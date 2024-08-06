package kernel

import (
	"context"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	// RepoIO defines the interface for repository operations.
	RepoIO interface {
		// GetProviderInfo retrieves the name and default branch for the provider repo.
		GetProviderInfo(ctx context.Context, id string) (*defs.RepoIOProviderInfo, error)

		// SetEarlyWarning sets the early warning flag for the provider repo.
		SetEarlyWarning(ctx context.Context, id string, value bool) error

		// GetAllBranches retrieves all the branches for the provider repo.
		GetAllBranches(ctx context.Context, payload *defs.RepoIOProviderInfo) ([]string, error)

		// DetectChanges detects changes in the repository.
		DetectChanges(ctx context.Context, payload *defs.RepoIODetectChangesPayload) (*defs.RepoIOChanges, error)

		// MergePR merges a pull request.
		MergePR(ctx context.Context, payload *defs.RepoIOMergePRPayload) error

		// TokenizedCloneURL returns the URL with an OAuth token included.
		//
		// NOTE: Since the URL contains an OAuth token, it is best not to call this as activity.
		// LINK: https://github.com/orgs/community/discussions/24575#discussioncomment-3244524
		TokenizedCloneURL(ctx context.Context, payload *defs.RepoIOProviderInfo) (string, error)
	}
)
