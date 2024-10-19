package rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	User struct {
		ID         uuid.UUID   `json:"id"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
		OrgID      uuid.UUID   `json:"org_id"`
		Email      string      `json:"email"`
		FirstName  pgtype.Text `json:"first_name"`
		LastName   pgtype.Text `json:"last_name"`
		Password   pgtype.Text `json:"password"`
		IsActive   bool        `json:"is_active"`
		IsVerified bool        `json:"is_verified"`
	}
)

// import (
// 	"context"
// 	"time"

// 	"github.com/google/uuid"
// )

// type (
// 	Empty struct{}

// 	GetUserByEmailWithInfoRow struct {
// 		User
// 		Teams    []Team         `json:"teams"`
// 		Accounts []OauthAccount `json:"accounts"`
// 		Orgs     []Org          `json:"orgs"`
// 	}
// )

// // GetUserByEmailWithInfo returns a user by email with additional information. It is a wrapper around
// // GetUserByEmailFull but returns proper structs.
// func (q *Queries) GetUserByEmailWithInfo(ctx context.Context, email string) (GetUserByEmailWithInfoRow, error) {
// 	result, err := q.GetUserByEmailFull(ctx, email)
// 	if err != nil {
// 		return GetUserByEmailWithInfoRow{}, err
// 	}

// 	tr, ok := result.Teams.([]map[string]interface{})
// 	if !ok {
// 		return GetUserByEmailWithInfoRow{}, err
// 	}

// 	var teams []Team
// 	for _, teamData := range tr {
// 		teams = append(teams, Team{
// 			ID:        teamData["id"].(uuid.UUID),
// 			CreatedAt: teamData["created_at"].(time.Time),
// 			UpdatedAt: teamData["updated_at"].(time.Time),
// 			OrgID:     teamData["org_id"].(uuid.UUID),
// 			Name:      teamData["name"].(string),
// 			Slug:      teamData["slug"].(string),
// 		})
// 	}

// 	oar, ok := result.OauthAccounts.([]map[string]interface{})
// 	if !ok {
// 		return GetUserByEmailWithInfoRow{}, err
// 	}

// 	var accounts []OauthAccount
// 	for _, data := range oar {
// 		accounts = append(accounts, OauthAccount{
// 			ID:                data["id"].(uuid.UUID),
// 			CreatedAt:         data["created_at"].(time.Time),
// 			UpdatedAt:         data["updated_at"].(time.Time),
// 			UserID:            data["user_id"].(uuid.UUID),
// 			Provider:          data["provider"].(string),
// 			ProviderAccountID: data["provider_account_id"].(string),
// 			ExpiresAt:         data["expires_at"].(pgtype.Timestamptz),
// 			Type:              data["type"].(pgtype.Text),
// 		})
// 	}

// 	or, ok := result.Orgs.([]map[string]interface{})
// 	if !ok {
// 		return GetUserByEmailWithInfoRow{}, err
// 	}

// 	var orgs []Org
// 	for _, data := range or {
// 		orgs = append(orgs, Org{
// 			ID:        data["id"].(uuid.UUID),
// 			CreatedAt: data["created_at"].(time.Time),
// 			UpdatedAt: data["updated_at"].(time.Time),
// 			Name:      data["name"].(string),
// 			Slug:      data["slug"].(string),
// 		})
// 	}

// 	return GetUserByEmailWithInfoRow{
// 		User: User{
// 			ID:         result.ID,
// 			CreatedAt:  result.CreatedAt,
// 			UpdatedAt:  result.UpdatedAt,
// 			OrgID:      result.OrgID,
// 			Email:      result.Email,
// 			FirstName:  result.FirstName,
// 			LastName:   result.LastName,
// 			Password:   result.Password,
// 			IsActive:   result.IsActive,
// 			IsVerified: result.IsVerified,
// 		},
// 		Teams:    teams,
// 		Accounts: accounts,
// 		Orgs:     orgs,
// 	}, nil
// }
