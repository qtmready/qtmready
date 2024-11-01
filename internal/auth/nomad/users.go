package authnmd

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/cast"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	// UserService implements the UserServiceHandler interface for managing user operations.
	UserService struct {
		authv1connect.UnimplementedUserServiceHandler
	}
)

var (
	NoOrgUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
)

// CreateUser creates a new user on the platform.
// If the organization with the given domain does not exist, it is created.
// The first user of an organization is an administrator, subsequent users are assigned the "member" role.
func (s *UserService) CreateUser(
	ctx context.Context, req *connect.Request[authv1.CreateUserRequest],
) (*connect.Response[authv1.AuthUser], error) { // Default public value for new users.
	role := "member"                                // Default role for new users.
	params := cast.ProtoToCreateUserParams(req.Msg) // Extract user creation parameters (excluding organization ID).
	domain := req.Msg.GetDomain()                   // Extract the domain name to locate the organization.

	// Initiate a database transaction.
	tx, qtx, err := db.Transaction(ctx)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	defer func() { _ = tx.Rollback(ctx) }() // Rollback is deferred to ensure rollback on error.

	// User sign-up without an organization domain.
	if domain == "" {
		params.OrgID = NoOrgUUID
		user, err := qtx.CreateUser(ctx, params)

		if err != nil {
			return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
		}

		return connect.NewResponse(&authv1.AuthUser{User: cast.UserToProto(&user)}), nil
	}

	// Retrieve the organization associated with the provided domain.
	org, err := qtx.GetOrgByDomain(ctx, domain)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
		}

		// Generate a unique slug for the organization based on the domain name.
		slug := db.CreateSlug(domain)

		// Create the organization in the database.
		org, err = qtx.CreateOrg(ctx, entities.CreateOrgParams{Name: domain, Lower: domain, Slug: slug})
		if err != nil {
			return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
		}

		role = "admin" // Assign the "admin" role to the first user of the organization.
	}

	// Update the user creation parameters with the organization ID.
	params.OrgID = org.ID

	// Create the user in the database.
	user, err := qtx.CreateUser(ctx, params)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	// Assign the appropriate role to the user within the organization.
	_, err = qtx.CreateUserRole(ctx, entities.CreateUserRoleParams{
		Name:   role,
		UserID: user.ID,
		OrgID:  org.ID,
	})
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	// Commit the database transaction and then Retrieve the user details for accurate relationships.
	if err := tx.Commit(ctx); err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	// Retrieve the user details from the database.
	details, err := db.Queries().GetAuthUserByID(ctx, user.ID)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	// Convert the retrieved user details to a protobuf structure.
	proto, err := cast.AuthUserQueryResponseToProto(
		details.User, details.Org, details.Roles, details.OauthAccounts, details.Teams,
	)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	// Return a successful response containing the created user information as a protobuf struct.
	return connect.NewResponse(proto), nil
}

// GetUserByProviderAccount retrieves a user based on their provider and provider account ID.
func (s *UserService) GetUserByProviderAccount(
	ctx context.Context, request *connect.Request[authv1.GetUserByProviderAccountRequest],
) (*connect.Response[authv1.AuthUser], error) {
	params := entities.GetUserByProviderAccountParams{
		Provider:          cast.ProtoToAuthProvider(request.Msg.GetProvider()),
		ProviderAccountID: request.Msg.GetProviderAccountId(),
	}

	one, err := db.Queries().GetUserByProviderAccount(ctx, params)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError(
				"entity", "users",
				"provider", request.Msg.GetProvider().String(),
				"provider_account_id", request.Msg.GetProviderAccountId(),
			).ToConnectError()
		}

		slog.Error("unable to get error", "error", err.Error())

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	details, err := db.Queries().GetAuthUserByID(ctx, one.ID)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	proto, err := cast.AuthUserQueryResponseToProto(
		details.User, details.Org, details.Roles, details.OauthAccounts, details.Teams,
	)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(proto), nil
}

// GetUserByEmail retrieves a user based on their email address.
func (s *UserService) GetUserByEmail(
	ctx context.Context, req *connect.Request[authv1.GetUserByEmailRequest],
) (*connect.Response[authv1.AuthUser], error) {
	details, err := db.Queries().GetAuthUserByEmail(ctx, req.Msg.GetEmail())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError("entity", "users", "email", req.Msg.GetEmail()).ToConnectError()
		}

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	proto, err := cast.AuthUserQueryResponseToProto(
		details.User, details.Org, details.Roles, details.OauthAccounts, details.Teams,
	)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(proto), nil
}

// GetUserByID retrieves a user based on their user ID.
func (s *UserService) GetUserByID(
	ctx context.Context, req *connect.Request[authv1.GetUserByIDRequest],
) (*connect.Response[authv1.AuthUser], error) {
	details, err := db.Queries().GetAuthUserByID(ctx, uuid.MustParse(req.Msg.GetId()))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError("entity", "users").ToConnectError()
		}

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	proto, err := cast.AuthUserQueryResponseToProto(
		details.User, details.Org, details.Roles, details.OauthAccounts, details.Teams,
	)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(proto), nil
}

// UpdateUser updates the user details in the database.
func (s *UserService) UpdateUser(
	ctx context.Context, req *connect.Request[authv1.UpdateUserRequest],
) (*connect.Response[authv1.UpdateUserResponse], error) {
	params := cast.ProtoToUpdateUserParams(req.Msg)

	user, err := db.Queries().UpdateUser(ctx, params)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(&authv1.UpdateUserResponse{User: cast.UserToProto(&user)}), nil
}

// NewUserSericeServiceHandler creates a new UserServiceHandler instance and returns the service name and handler.
func NewUserSericeServiceHandler() (string, http.Handler) {
	return authv1connect.NewUserServiceHandler(&UserService{})
}
