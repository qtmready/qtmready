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

package ws

import (
	"context"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	Activities struct{}

	TeamUsersReponse struct {
		IDs []string `json:"users"`
	}
)

func (a *Activities) SendMessage(ctx context.Context, id string, message []byte) (bool, error) {
	return instance.send_local(id, message), nil
}

func (a *Activities) RouteMessage(ctx context.Context, id string, message []byte) error {
	return instance.Send(ctx, id, message)
}

func (a *Activities) Signal(ctx context.Context, signal shared.WorkflowSignal, payload any) error {
	return instance.Signal(ctx, signal, payload)
}

func (a *Activities) GetUserQueue(ctx context.Context, id string) (string, error) {
	return instance.query(ctx, id)
}

func (a *Activities) GetTeamUsers(ctx context.Context, team_id string) (*TeamUsersReponse, error) {
	users := make([]auth.User, 0)
	err := db.Filter(&auth.Team{}, &users, db.QueryParams{"team_id": team_id})

	if err != nil {
		return nil, err
	}

	ids := make([]string, len(users))
	for i, user := range users {
		ids[i] = user.ID.String()
	}

	return &TeamUsersReponse{IDs: ids}, nil
}
