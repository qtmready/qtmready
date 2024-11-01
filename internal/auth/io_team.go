// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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

	"github.com/scylladb/gocqlx/v2/qb"

	"go.breu.io/quantm/internal/db"
)

type (
	// teamio represents the activities for the team.
	teamio struct{}
)

var (
	teamonce *teamio
	teamsync sync.Once
)

// TeamIO creates and returns a new TeamIO object.
//
// Example:
//
//	team_io := auth.TeamIO()
func TeamIO() *teamio {
	teamsync.Do(func() {
		teamonce = &teamio{}
	})

	return teamonce
}

// Get retrieves a team from the database based on the provided parameters.
//
// Example:
//
//	team, err := auth.TeamIO().Get(ctx, db.QueryParams{"id": team_id})
func (a *teamio) Get(ctx context.Context, params db.QueryParams) (*Team, error) {
	team := &Team{}

	return team, db.Get(team, params)
}

// GetByID retrieves a team from the database by their ID.
//
// Example:
//
//	team, err := auth.TeamIO().GetByID(ctx, team_id)
func (a *teamio) GetByID(ctx context.Context, id string) (*Team, error) {
	team := &Team{}

	return team, db.Get(team, db.QueryParams{"id": id})
}

// GetByName retrieves a team from the database by their name.
//
// Example:
//
//	team, err := auth.TeamIO().GetByName(ctx, "My Team")
func (a *teamio) GetByName(ctx context.Context, name string) (*Team, error) {
	team := &Team{}

	return team, db.Get(team, db.QueryParams{"name": name})
}

// GetUsers retrieves all users associated with a team.
//
// Example:
//
//	users, err := auth.TeamIO().GetUsers(ctx, team)
func (a *teamio) GetUsers(ctx context.Context, team *Team) ([]User, error) {
	entity := &User{}
	users := make([]User, 0)
	ids := make([]string, 0)

	tus, err := TeamUserIO().GetByTeamID(ctx, team.ID.String())
	if err != nil {
		return nil, err
	}

	for _, tu := range tus {
		ids = append(ids, tu.UserID.String())
	}

	query := db.
		SelectBuilder(entity.GetTable().Name()).
		Where(qb.In("id"))

	err = db.Cassandra().
		Session.
		Query(query.ToCql()).BindMap(qb.M{"id": ids}).
		GetRelease(users)

	return users, err
}

// Save saves a team to the database.
//
// Note: Always provide a pointer to the complete Team object to avoid
// creating a copy. The Save method will update the provided Team object
// with any changes made by the database.
//
// Example:
//
//	team, err := auth.TeamIO().Save(ctx, team)
func (a *teamio) Save(ctx context.Context, team *Team) (*Team, error) {
	return team, db.Save(team)
}
