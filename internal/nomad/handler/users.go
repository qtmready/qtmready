package handler

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
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
	return nil, nil
}

func (s *UserService) GetUserByEmail(
	ctx context.Context, req *connect.Request[authv1.GetUserByEmailRequest],
) (*connect.Response[authv1.GetUserByEmailResponse], error) {
	user, err := db.Queries().GetUserByEmail(ctx, req.Msg.GetEmail())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, erratic.NewNotFoundError(
				"entity", "users",
				"email", req.Msg.GetEmail(),
			).ToConnectError()
		}

		return nil, erratic.NewInternalServerError().DataBaseError(err).ToConnectError()
	}

	return connect.NewResponse(&authv1.GetUserByEmailResponse{User: convert.UserToProto(&user)}), nil
}

func NewUserSericeServiceHandler() (string, http.Handler) {
	return authv1connect.NewUserServiceHandler(&UserService{})
}
