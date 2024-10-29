package handlers

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/core/v1/corev1connect"
)

type (
	CoreRepoService struct {
		corev1connect.UnimplementedRepoServiceHandler
	}
)

func (s *CoreRepoService) CreateRepo(
	ctx context.Context, req *connect.Request[corev1.CreateRepoRequest],
) (*connect.Response[corev1.CreateRepoResponse], error) {
	return nil, nil
}

func (s *CoreRepoService) GetOrgReposByOrgID(
	ctx context.Context, req *connect.Request[corev1.GetOrgReposByOrgIDRequest],
) (*connect.Response[corev1.GetOrgReposByOrgIDResponse], error) {
	return nil, nil
}

func (s *CoreRepoService) GetRepoByID(ctx context.Context, req *connect.Request[corev1.GetRepoByIDRequest],
) (*connect.Response[corev1.GetRepoByIDResponse], error) {
	return nil, nil
}

func NewCoreRepoServiceHandler() (string, http.Handler) {
	return corev1connect.NewRepoServiceHandler(&CoreRepoService{})
}
