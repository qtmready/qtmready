// Copyright © 2024, Breu, Inc. <info@breu.io>
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

package code

import (
	"context"

	"go.breu.io/quantm/internal/core/defs" // Import Repo from defs
	"go.breu.io/quantm/internal/db"
)

type (
	// repoio represents the activities for the repo.
	repoio struct{}
)

// RepoIO creates and returns a new RepoIO object.
//
// Example:
//
//	repo_io := code.RepoIO()
func RepoIO() *repoio {
	return &repoio{}
}

// Get retrieves a repo from the database based on the provided parameters.
//
// Example:
//
//	repo, err := code.RepoIO().Get(ctx, db.QueryParams{"id": repo_id})
func (r *repoio) Get(ctx context.Context, params db.QueryParams) (*defs.Repo, error) {
	repo := &defs.Repo{}

	return repo, db.Get(repo, params)
}

// GetByID retrieves a repo from the database by their ID.
//
// Example:
//
//	repo, err := code.RepoIO().GetByID(ctx, repo_id)
func (r *repoio) GetByID(ctx context.Context, id string) (*defs.Repo, error) {
	repo := &defs.Repo{}

	if err := db.Get(repo, db.QueryParams{"id": id}); err != nil {
		return nil, err
	}

	return repo, nil
}

// GetByCtrlID retrieves a repo from the database by their control ID.
//
// Example:
//
//	repo, err := code.RepoIO().GetByCtrlID(ctx, ctrl_id)
func (r *repoio) GetByCtrlID(ctx context.Context, ctrl_id string) (*defs.Repo, error) {
	repo := &defs.Repo{}

	return repo, db.Get(repo, db.QueryParams{"ctrl_id": ctrl_id})
}

// Save saves a repo to the database.
//
// Note: Always provide a pointer to the complete Repo object to avoid creating a copy. The Save method will update the
// provided Repo object with any changes made by the database.
//
// Example:
//
//	repo, err := code.RepoIO().Save(ctx, repo)
func (r *repoio) Save(ctx context.Context, repo *defs.Repo) (*defs.Repo, error) {
	return repo, db.Save(repo)
}

func SaveRepoEvent[T defs.EventPayload, P defs.EventProvider](ctx context.Context, event defs.Event[T, P]) error {
	flat, err := event.Flatten()
	if err != nil {
		return err
	}

	return db.Save(flat)
}
