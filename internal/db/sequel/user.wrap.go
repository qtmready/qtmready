package sequel

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
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/shared"
)

type (
	// CreateUserRequest represents the request body for creating a new user.
	CreateUserRequest struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password"`
	}

	// GetUserRequest represents the request body for retrieving a user by ID.
	GetUserRequest struct {
		ID uuid.UUID `json:"id"`
	}

	// GetUserByEmailRequest represents the request body for retrieving a user by email.
	GetUserByEmailRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	// GetUserByIDRequest represents the request body for retrieving a user by ID and OrgID.
	GetUserByIDRequest struct {
		ID    uuid.UUID `json:"id"`
		OrgID uuid.UUID `json:"org_id"`
	}

	// GetUserByProviderAccountRequest represents the request body for retrieving a user by provider and provider account
	// ID.
	GetUserByProviderAccountRequest = entities.GetUserByProviderAccountParams

	// User represents a user in the system.
	User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		OrgID     uuid.UUID `json:"org_id"`
	}

	// GetUserByEmailFullResponse represents the response body for retrieving a user by email with full information.
	GetUserByEmailFullResponse struct {
		ID        uuid.UUID               `json:"id"`
		CreatedAt time.Time               `json:"created_at"`
		UpdatedAt time.Time               `json:"updated_at"`
		FirstName string                  `json:"first_name"`
		LastName  string                  `json:"last_name"`
		Email     string                  `json:"email"`
		OrgID     uuid.UUID               `json:"org_id"`
		Teams     []entities.Team         `json:"teams"`
		Accounts  []entities.OauthAccount `json:"accounts"`
		Orgs      []entities.Org          `json:"orgs"`
	}

	// UpdateUserRequest represents the request body for updating a user.
	UpdateUserRequest struct {
		ID        uuid.UUID `json:"id"`
		FirstName string    `json:"first_name" validate:"required"`
		LastName  string    `json:"last_name" validate:"required"`
		Email     string    `json:"email" validate:"required,email"`
		OrgID     uuid.UUID `json:"org_id"`
	}

	// UpdateUserPasswordRequest represents the request body for updating a user's password.
	UpdateUserPasswordRequest struct {
		ID       uuid.UUID `json:"id"`
		Password string    `json:"password" validate:"required"`
	}
)

// CreateUser creates a new user.
func CreateUser(ctx context.Context, req CreateUserRequest) (entities.CreateUserRow, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return entities.CreateUserRow{}, erratic.NewBadRequestError().SetVaidationErrors(err)
	}

	var password string

	if req.Password == "" {
		var err error
		password, err = auth.GeneratePassword(12)

		if err != nil {
			return entities.CreateUserRow{}, erratic.NewInternalServerError().SetInternal(err)
		}
	} else {
		password = req.Password
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entities.CreateUserRow{}, erratic.NewInternalServerError().SetInternal(err)
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
	// Check if authenticated
	val := ctx.Value("org_id")
	if val == nil {
		return entities.User{}, erratic.NewUnauthorizedError()
	}

	user, err := db.Queries().GetUser(ctx, req.ID)
	if err != nil {
		return entities.User{}, erratic.NewNotFoundError().AddHint("user_id", req.ID.String())
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email.
func GetUserByEmail(ctx context.Context, req GetUserByEmailRequest) (User, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return User{}, erratic.NewBadRequestError().SetVaidationErrors(err)
	}

	// Check if authenticated
	val := ctx.Value("org_id")
	if val == nil {
		return User{}, erratic.NewUnauthorizedError()
	}

	row, err := db.Queries().GetUserByEmail(ctx, req.Email)
	if err != nil {
		return User{}, erratic.NewNotFoundError()
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
	// Check if authenticated
	val := ctx.Value("org_id")
	if val == nil {
		return User{}, erratic.NewUnauthorizedError()
	}

	org_id, ok := val.(uuid.UUID)
	if !ok {
		return User{}, erratic.NewInternalServerError()
	}

	row, err := db.Queries().GetUserByID(ctx, req.ID)
	if err != nil {
		return User{}, erratic.NewNotFoundError()
	}

	if row.OrgID != org_id {
		return User{}, erratic.NewForbiddenError().AddHint("org_id", org_id.String())
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
	// Check if authenticated
	val := ctx.Value("org_id")
	if val == nil {
		return User{}, erratic.NewUnauthorizedError()
	}

	row, err := db.Queries().GetUserByProviderAccount(ctx, req)
	if err != nil {
		return User{}, erratic.NewNotFoundError()
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
		return GetUserByEmailFullResponse{}, erratic.NewBadRequestError().SetVaidationErrors(err)
	}

	// Check if authenticated
	val := ctx.Value("org_id")
	if val == nil {
		return GetUserByEmailFullResponse{}, erratic.NewUnauthorizedError()
	}

	row, err := db.Queries().GetUserByEmailFull(ctx, req.Email)
	if err != nil {
		return GetUserByEmailFullResponse{}, erratic.NewNotFoundError(err.Error())
	}

	var teams []entities.Team

	err = json.Unmarshal(row.Teams.([]byte), &teams)
	if err != nil {
		return GetUserByEmailFullResponse{}, erratic.NewInternalServerError().SetInternal(err)
	}

	var accounts []entities.OauthAccount

	err = json.Unmarshal(row.OauthAccounts.([]byte), &accounts)
	if err != nil {
		return GetUserByEmailFullResponse{}, erratic.NewInternalServerError().SetInternal(err)
	}

	var orgs []entities.Org

	err = json.Unmarshal(row.Orgs.([]byte), &orgs)
	if err != nil {
		return GetUserByEmailFullResponse{}, erratic.NewInternalServerError().SetInternal(err)
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
	// Check if authenticated
	val := ctx.Value("org_id")
	if val == nil {
		return entities.UpdateUserRow{}, erratic.NewUnauthorizedError()
	}

	_, ok := val.(uuid.UUID)
	if !ok {
		return entities.UpdateUserRow{}, erratic.NewInternalServerError()
	}

	if err := shared.Validator().Struct(req); err != nil {
		return entities.UpdateUserRow{}, erratic.NewBadRequestError().SetVaidationErrors(err)
	}

	row, err := db.Queries().UpdateUser(ctx, entities.UpdateUserParams{
		ID:        req.ID,
		FirstName: pgtype.Text{String: req.FirstName, Valid: true},
		LastName:  pgtype.Text{String: req.LastName, Valid: true},
		Lower:     req.Email,
		OrgID:     req.OrgID,
	})
	if err != nil {
		return entities.UpdateUserRow{}, erratic.NewInternalServerError().DataBaseError(err)
	}

	return row, nil
}

// UpdateUserPassword updates a user's password.
func UpdateUserPassword(ctx context.Context, req UpdateUserPasswordRequest) (User, error) {
	// Check if authenticated
	val := ctx.Value("org_id")
	if val == nil {
		return User{}, erratic.NewUnauthorizedError()
	}

	_, ok := val.(uuid.UUID)
	if !ok {
		return User{}, erratic.NewInternalServerError()
	}

	if err := shared.Validator().Struct(req); err != nil {
		return User{}, erratic.NewBadRequestError().SetVaidationErrors(err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, erratic.NewInternalServerError().DataBaseError(err)
	}

	row, err := db.Queries().UpdateUserPassword(ctx, entities.UpdateUserPasswordParams{
		ID:       req.ID,
		Password: pgtype.Text{String: string(hashed), Valid: true},
	})
	if err != nil {
		return User{}, erratic.NewInternalServerError().DataBaseError(err)
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
