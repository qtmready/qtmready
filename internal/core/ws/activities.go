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
