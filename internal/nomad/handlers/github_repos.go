package handlers

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
	"go.breu.io/quantm/internal/proto/hooks/github/v1/githubv1connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	GithubRepoService struct {
		githubv1connect.UnimplementedGithubRepoServiceHandler
	}
)

func (s *GithubRepoService) GithubInstall(
	ctx context.Context, req *connect.Request[githubv1.GithubInstallRequest],
) (*connect.Response[emptypb.Empty], error) {
	return nil, nil
}

func NewGithubRepoServiceHandler() (string, http.Handler) {
	return githubv1connect.NewGithubRepoServiceHandler(&GithubRepoService{})
}
