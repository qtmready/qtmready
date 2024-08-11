package ws

import (
	"context"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
)

type (
	Activities struct{}
)

func (a *Activities) SendMessage(userID string, message []byte) (bool, error) {
	return instance.send_local(userID, message), nil
}

func (a *Activities) GetTeamUsers(ctx context.Context, teamID string) ([]string, error) {
	users := make([]auth.User, 0)
	err := db.Filter(&auth.Team{}, &users, db.QueryParams{"team_id": teamID})

	if err != nil {
		return nil, err
	}

	ids := make([]string, len(users))
	for i, user := range users {
		ids[i] = user.ID.String()
	}

	return ids, nil
}
