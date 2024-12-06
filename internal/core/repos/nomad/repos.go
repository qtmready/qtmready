package reposnmd

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/auth"
	reposcast "go.breu.io/quantm/internal/core/repos/cast"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/erratic"
	corev1 "go.breu.io/quantm/internal/proto/ctrlplane/core/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/core/v1/corev1connect"
)

type (
	RepoService struct {
		corev1connect.UnimplementedRepoServiceHandler
	}
)

func (s *RepoService) ListRepos(
	ctx context.Context, req *connect.Request[emptypb.Empty],
) (*connect.Response[corev1.ListReposResponse], error) {
	_, org_id := auth.NomadAuthContext(ctx)

	rows, err := db.Queries().ListRepos(ctx, org_id)
	if err != nil {
		return nil, erratic.NewDatabaseError(erratic.CoreModule).Wrap(err)
	}

	protos := reposcast.RepoExtendedListToProto(rows)

	return connect.NewResponse(&corev1.ListReposResponse{Repos: protos}), nil
}

func NewRepoServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return corev1connect.NewRepoServiceHandler(&RepoService{}, opts...)
}
