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

package auth

import (
	"context"
	"sync"

	"go.breu.io/quantm/internal/db"
)

type (
	// TeamUserIO represents the activities for the team user.
	teamusrio struct{}
)

var (
	teamusronce *teamusrio
	teamusrsync sync.Once
)

// TeamUserIO creates and returns a new TeamUserIO object.
//
// Example:
//
//	team_user_io := auth.TeamUserIO()
func TeamUserIO() *teamusrio {
	teamusrsync.Do(func() {
		teamusronce = &teamusrio{}
	})

	return teamusronce
}

// Get retrieves a team user from the database by their user ID and team ID.
//
// Example:
//
//	team_user, err := auth.TeamUserIO().Get(ctx, user_id, team_id)
func (a *teamusrio) Get(ctx context.Context, user_id, team_id string) (*TeamUser, error) {
	team_user := &TeamUser{}

	return team_user, db.Get(team_user, db.QueryParams{
		"id":      user_id,
		"team_id": team_id,
	})
}

// GetByUserID retrieves a team user from the database by their user ID.
//
// Example:
//
//	team_users, err := auth.TeamUserIO().GetByUserID(ctx, user_id)
func (a *teamusrio) GetByUserID(ctx context.Context, user_id string) ([]TeamUser, error) {
	tus := make([]TeamUser, 0)
	err := db.Filter(&TeamUser{}, tus, db.QueryParams{"id": user_id})

	return tus, err
}

// GetByTeamID retrieves a team user from the database by their team ID.
//
// Example:
//
//	team_users, err := auth.TeamUserIO().GetByTeamID(ctx, team_id)
func (a *teamusrio) GetByTeamID(ctx context.Context, team_id string) ([]TeamUser, error) {
	tus := make([]TeamUser, 0)
	err := db.Filter(&TeamUser{}, tus, db.QueryParams{"team_id": team_id})

	return tus, err
}

// GetByLogin retrieves a team user from the database by their login.
//
// Example:
//
//	team_user, err := auth.TeamUserIO().GetByLogin(ctx, login)
func (a *teamusrio) GetByLogin(ctx context.Context, login string) (*TeamUser, error) {
	team_user := &TeamUser{}

	return team_user, db.Get(team_user, db.QueryParams{"user_login_id": login})
}

// Save creates or updates a team user in the database.
//
// Example:
//
//	team_user, err := auth.TeamUserIO().Save(ctx, team_user)
func (a *teamusrio) Save(ctx context.Context, tu *TeamUser) (*TeamUser, error) {
	fetched := &TeamUser{}

	if err := db.Get(fetched, db.QueryParams{"team_id": tu.TeamID.String(), "user_id": tu.UserID.String()}); err == nil {
		tu.ID = fetched.ID
		tu.CreatedAt = fetched.CreatedAt
	}

	if err := db.Save(tu); err != nil {
		return nil, err
	}

	return tu, nil
}
