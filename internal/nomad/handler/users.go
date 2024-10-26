package handler

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	authv1 "go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/nomad/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	UserService struct {
		authv1connect.UnimplementedUserServiceHandler
	}
)

func (s *UserService) CreateUser(
	ctx context.Context,
	req *connect.Request[authv1.CreateUserRequest],
) (*connect.Response[authv1.CreateUserResponse], error) {
	return nil, nil
}

func NewUserSericeServiceHandler() (string, http.Handler) {
	return authv1connect.NewUserServiceHandler(&UserService{})
}
