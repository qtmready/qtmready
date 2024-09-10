// Copyright © 2023, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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
