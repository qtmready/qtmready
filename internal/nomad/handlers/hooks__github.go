package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
	"go.breu.io/quantm/internal/proto/hooks/github/v1/githubv1connect"
)

type (
	GithubRepoService struct {
		githubv1connect.UnimplementedGithubRepoServiceHandler
	}
)

func (s *GithubRepoService) GithubInstall(
	ctx context.Context, req *connect.Request[githubv1.GithubInstallRequest],
) (*connect.Response[emptypb.Empty], error) {
	slog.Info("completing github installation", "request", req)
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewGithubRepoServiceHandler() (string, http.Handler) {
	return githubv1connect.NewGithubRepoServiceHandler(&GithubRepoService{})
}
