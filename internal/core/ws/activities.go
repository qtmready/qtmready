package ws

import (
	"context"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	Activities struct{}
)

func NewActivities() *Activities {
	return &Activities{}
}

func (a *Activities) SendMessage(ctx context.Context, userID string, message []byte) error {
	return Instance().Send(ctx, userID, message)
}

func (a *Activities) GetTeamUsers(ctx context.Context, teamID string) ([]string, error) {
	users := make([]auth.User, 0)
	err := db.Filter(&auth.Team{}, &users, db.QueryParams{"team_id": teamID})

	if err != nil {
		return nil, err
	}

	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID.String()
	}

	return userIDs, nil
}

func (a *Activities) BroadcastMessage(ctx context.Context, teamID string, message []byte) error {
	userIDs, err := a.GetTeamUsers(ctx, teamID)
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		if err := a.SendMessage(ctx, userID, message); err != nil {
			shared.Logger().Error("Failed to send message to user", "user_id", userID, "error", err)
		}
	}

	return nil
}
