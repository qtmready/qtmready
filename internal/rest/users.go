package rest

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/shared"
)

type (
	// CreateUserRequest represents the request body for creating a new user.
	CreateUserRequest struct {
		FirstName string `json:"first_name" valid:"required"`  // First name of the user.
		LastName  string `json:"last_name" valid:"required"`   // Last name of the user.
		Email     string `json:"email" valid:"required,email"` // Email address of the user.
		Password  string `json:"password"`                     // Password of the user.
	}

	// GetUserRequest represents the request body for retrieving a user by ID.
	GetUserRequest struct {
		ID uuid.UUID `json:"id"` // ID of the user to retrieve.
	}

	// GetUserByEmailRequest represents the request body for retrieving a user by email.
	GetUserByEmailRequest struct {
		Email string `json:"email" valid:"required,email"` // Email address of the user to retrieve.
	}

	// GetUserByIDRequest represents the request body for retrieving a user by ID and OrgID.
	GetUserByIDRequest struct {
		ID    uuid.UUID `json:"id"`     // ID of the user to retrieve.
		OrgID uuid.UUID `json:"org_id"` // ID of the organization the user belongs to.
	}

	// GetUserByProviderAccountRequest represents the request body for retrieving a user by provider and provider account
	// ID.
	GetUserByProviderAccountRequest = entities.GetUserByProviderAccountParams

	// User represents a user in the system.
	User struct {
		ID        uuid.UUID `json:"id"`         // Unique identifier of the user.
		CreatedAt time.Time `json:"created_at"` // Time when the user was created.
		UpdatedAt time.Time `json:"updated_at"` // Time when the user was last updated.
		FirstName string    `json:"first_name"` // First name of the user.
		LastName  string    `json:"last_name"`  // Last name of the user.
		Email     string    `json:"email"`      // Email address of the user.
		OrgID     uuid.UUID `json:"org_id"`     // ID of the organization the user belongs to.
	}

	// GetUserByEmailFullResponse represents the response body for retrieving a user by email with full information.
	GetUserByEmailFullResponse struct {
		ID        uuid.UUID               `json:"id"`         // Unique identifier of the user.
		CreatedAt time.Time               `json:"created_at"` // Time when the user was created.
		UpdatedAt time.Time               `json:"updated_at"` // Time when the user was last updated.
		FirstName string                  `json:"first_name"` // First name of the user.
		LastName  string                  `json:"last_name"`  // Last name of the user.
		Email     string                  `json:"email"`      // Email address of the user.
		OrgID     uuid.UUID               `json:"org_id"`     // ID of the organization the user belongs to.
		Teams     []entities.Team         `json:"teams"`      // List of teams the user belongs to.
		Accounts  []entities.OauthAccount `json:"accounts"`   // List of OAuth accounts linked to the user.
		Orgs      []entities.Org          `json:"orgs"`       // List of organizations the user belongs to.
	}

	// UpdateUserRequest represents the request body for updating a user.
	UpdateUserRequest struct {
		ID        uuid.UUID `json:"id"`                           // Unique identifier of the user to update.
		FirstName string    `json:"first_name" valid:"required"`  // First name of the user.
		LastName  string    `json:"last_name" valid:"required"`   // Last name of the user.
		Email     string    `json:"email" valid:"required,email"` // Email address of the user.
		OrgID     uuid.UUID `json:"org_id"`                       // ID of the organization the user belongs to.
	}

	// UpdateUserPasswordRequest represents the request body for updating a user's password.
	UpdateUserPasswordRequest struct {
		ID       uuid.UUID `json:"id"`                        // Unique identifier of the user to update.
		Password string    `json:"password" valid:"required"` // New password of the user.
	}
)

// CreateUser creates a new user.
func CreateUser(ctx context.Context, req CreateUserRequest) (entities.CreateUserRow, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return entities.CreateUserRow{}, NewBadRequestError().FormatValidationError(err)
	}

	var password string

	if req.Password == "" {
		var err error
		password, err = auth.GeneratePassword(12)
		if err != nil {
			return entities.CreateUserRow{}, NewInternalServerError("reason", err.Error())
		}
	} else {
		password = req.Password
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entities.CreateUserRow{}, NewInternalServerError("reason", err.Error())
	}

	return db.Queries().CreateUser(ctx, entities.CreateUserParams{
		FirstName: pgtype.Text{String: req.FirstName, Valid: true},
		LastName:  pgtype.Text{String: req.LastName, Valid: true},
		Email:     req.Email,
		Password:  pgtype.Text{String: string(hashed), Valid: true},
	})
}

// GetUser retrieves a user by ID.
func GetUser(ctx context.Context, req GetUserRequest) (entities.User, error) {
	user, err := db.Queries().GetUser(ctx, req.ID)
	if err != nil {
		return entities.User{}, NewNotFoundError(err.Error())
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email.
func GetUserByEmail(ctx context.Context, req GetUserByEmailRequest) (User, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return User{}, NewBadRequestError().FormatValidationError(err)
	}

	row, err := db.Queries().GetUserByEmail(ctx, req.Email)
	if err != nil {
		return User{}, NewNotFoundError(err.Error())
	}

	return User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		FirstName: row.FirstName.String,
		LastName:  row.LastName.String,
		Email:     row.Email,
		OrgID:     row.OrgID,
	}, nil
}

// GetUserByID retrieves a user by ID and OrgID.
func GetUserByID(ctx context.Context, req GetUserByIDRequest) (User, error) {
	row, err := db.Queries().GetUserByID(ctx, req.ID)
	if err != nil {
		return User{}, NewNotFoundError(err.Error())
	}

	if row.OrgID != req.OrgID {
		return User{}, NewUnauthorizedError("user_id", req.ID.String())
	}

	return User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		FirstName: row.FirstName.String,
		LastName:  row.LastName.String,
		Email:     row.Email,
		OrgID:     row.OrgID,
	}, nil
}

// GetUserByProviderAccount retrieves a user by provider and provider account ID.
func GetUserByProviderAccount(ctx context.Context, req GetUserByProviderAccountRequest) (User, error) {
	row, err := db.Queries().GetUserByProviderAccount(ctx, req)
	if err != nil {
		return User{}, NewNotFoundError(err.Error())
	}

	return User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		FirstName: row.FirstName.String,
		LastName:  row.LastName.String,
		Email:     row.Email,
		OrgID:     row.OrgID,
	}, nil
}

// GetUserByEmailFull retrieves a user by email with full information.
func GetUserByEmailFull(ctx context.Context, req GetUserByEmailRequest) (GetUserByEmailFullResponse, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return GetUserByEmailFullResponse{}, NewBadRequestError().FormatValidationError(err)
	}

	row, err := db.Queries().GetUserByEmailFull(ctx, req.Email)
	if err != nil {
		return GetUserByEmailFullResponse{}, NewNotFoundError(err.Error())
	}

	var teams []entities.Team
	err = json.Unmarshal(row.Teams.([]byte), &teams)
	if err != nil {
		return GetUserByEmailFullResponse{}, NewInternalServerError("reason", err.Error())
	}

	var accounts []entities.OauthAccount
	err = json.Unmarshal(row.OauthAccounts.([]byte), &accounts)
	if err != nil {
		return GetUserByEmailFullResponse{}, NewInternalServerError("reason", err.Error())
	}

	var orgs []entities.Org
	err = json.Unmarshal(row.Orgs.([]byte), &orgs)
	if err != nil {
		return GetUserByEmailFullResponse{}, NewInternalServerError("reason", err.Error())
	}

	return GetUserByEmailFullResponse{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		FirstName: row.FirstName.String,
		LastName:  row.LastName.String,
		Email:     row.Email,
		OrgID:     row.OrgID,
		Teams:     teams,
		Accounts:  accounts,
		Orgs:      orgs,
	}, nil
}

// UpdateUser updates a user.
func UpdateUser(ctx context.Context, req UpdateUserRequest) (entities.UpdateUserRow, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return entities.UpdateUserRow{}, NewBadRequestError().FormatValidationError(err)
	}

	row, err := db.Queries().UpdateUser(ctx, entities.UpdateUserParams{
		ID:        req.ID,
		FirstName: pgtype.Text{String: req.FirstName, Valid: true},
		LastName:  pgtype.Text{String: req.LastName, Valid: true},
		Lower:     req.Email,
		OrgID:     req.OrgID,
	})
	if err != nil {
		return entities.UpdateUserRow{}, NewInternalServerError("reason", err.Error())
	}

	return row, nil
}

// UpdateUserPassword updates a user's password.
func UpdateUserPassword(ctx context.Context, req UpdateUserPasswordRequest) (User, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return User{}, NewBadRequestError().FormatValidationError(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, NewInternalServerError("reason", err.Error())
	}

	row, err := db.Queries().UpdateUserPassword(ctx, entities.UpdateUserPasswordParams{
		ID:       req.ID,
		Password: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return User{}, NewInternalServerError("reason", err.Error())
	}

	return User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		FirstName: row.FirstName.String,
		LastName:  row.LastName.String,
		Email:     row.Email,
		OrgID:     row.OrgID,
	}, nil
}
