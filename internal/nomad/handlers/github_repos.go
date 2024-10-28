package handlers

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
	"go.breu.io/quantm/internal/proto/hooks/github/v1/githubv1connect"
)

type (
	GithubRepoService struct {
		githubv1connect.UnimplementedRepoServiceHandler
	}
)

func (s GithubRepoService) CreateRepo(
	ctx context.Context, req *connect.Request[githubv1.CreateGithubRepoRequest],
) (*connect.Response[githubv1.CreateGithubRepoResponse], error) {
	return nil, nil
}

func (s GithubRepoService) GetGithubRepoByID(
	ctx context.Context, req *connect.Request[githubv1.GetGithubRepoByIDRequest],
) (*connect.Response[githubv1.GetGithubRepoByIDResponse], error) {
	return nil, nil
}

func (s GithubRepoService) GetGithubRepoByName(
	ctx context.Context, req *connect.Request[githubv1.GetGithubRepoByNameRequest],
) (*connect.Response[githubv1.GetGithubRepoByNameResponse], error) {
	return nil, nil
}

func NewGithubRepoServiceHandler() (string, http.Handler) {
	return githubv1connect.NewRepoServiceHandler(&GithubRepoService{})
}
