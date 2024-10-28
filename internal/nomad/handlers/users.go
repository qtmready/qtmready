package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/nomad/convert"
	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	UserService struct {
		authv1connect.UnimplementedUserServiceHandler
	}
)

func (s *UserService) CreateUser(
	ctx context.Context, req *connect.Request[authv1.CreateUserRequest],
) (*connect.Response[authv1.CreateUserResponse], error) {
	params := convert.ProtoToCreateUserParams(req.Msg) // protobuf to create user params (without org id).
	domain := req.Msg.GetDomain()                      // extract domain to lookup org.

	// Begin a database transaction.
	tx, qtx, err := db.Transaction(ctx)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	defer func() { _ = tx.Rollback(ctx) }() // rollback is deferred to ensure that we rollback on error.

	// Get the organization associated with the given domain.
	org, err := qtx.GetOrgByDomain(ctx, domain)
	if err != nil {
		if err != pgx.ErrNoRows {
			return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
		}

		// Generate a slug for the organization.
		slug := db.CreateSlug(domain)

		// Create the organization in the database.
		org, err = qtx.CreateOrg(ctx, entities.CreateOrgParams{Name: domain, Lower: domain, Slug: slug})
		if err != nil {
			return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
		}
	}

	// Set the organization ID in the CreateUserParams struct.
	params.OrgID = org.ID

	// Create the user in the database.
	user, err := qtx.CreateUser(ctx, params)
	if err != nil {
		// Return an internal server error if there's an error creating the user.
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	// Commit the database transaction.
	if err := tx.Commit(ctx); err != nil {
		// Return an internal server error if there's an error committing the transaction.
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	// Convert the created user to a protobuf struct and return a successful response.
	return connect.NewResponse(&authv1.CreateUserResponse{User: convert.UserToProto(&user)}), nil
}

func (s *UserService) GetUserByProviderAccount(
	ctx context.Context, request *connect.Request[authv1.GetUserByProviderAccountRequest],
) (*connect.Response[authv1.GetUserByProviderAccountResponse], error) {
	params := entities.GetUserByProviderAccountParams{
		Provider:          convert.ProtoToAuthProvider(request.Msg.GetProvider()),
		ProviderAccountID: request.Msg.GetProviderAccountId(),
	}

	user, err := db.Queries().GetUserByProviderAccount(ctx, params)
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

	return connect.NewResponse(&authv1.GetUserByProviderAccountResponse{User: convert.UserToProto(&user)}), nil
}

func (s *UserService) GetUserByEmail(
	ctx context.Context, req *connect.Request[authv1.GetUserByEmailRequest],
) (*connect.Response[authv1.GetUserByEmailResponse], error) {
	user, err := db.Queries().GetUserByEmail(ctx, req.Msg.GetEmail())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError("entity", "users", "email", req.Msg.GetEmail()).ToConnectError()
		}

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(&authv1.GetUserByEmailResponse{User: convert.UserToProto(&user)}), nil
}

func (s *UserService) GetUserByID(
	ctx context.Context, req *connect.Request[authv1.GetUserByIDRequest],
) (*connect.Response[authv1.GetUserByIDResponse], error) {
	user, err := db.Queries().GetUserByID(ctx, uuid.MustParse(req.Msg.GetId().GetValue()))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError("entity", "users", "id", req.Msg.GetId().GetValue()).ToConnectError()
		}

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(&authv1.GetUserByIDResponse{User: convert.UserToProto(&user)}), nil
}

func (s *UserService) UpdateUser(
	ctx context.Context, req *connect.Request[authv1.UpdateUserRequest],
) (*connect.Response[authv1.UpdateUserResponse], error) {
	params := convert.ProtoToUpdateUserParams(req.Msg)

	user, err := db.Queries().UpdateUser(ctx, params)
	if err != nil {
		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(&authv1.UpdateUserResponse{User: convert.UserToProto(&user)}), nil
}

func NewUserSericeServiceHandler() (string, http.Handler) {
	return authv1connect.NewUserServiceHandler(&UserService{})
}
