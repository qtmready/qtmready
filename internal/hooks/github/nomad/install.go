package nomad

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/hooks/github/defs"
	"go.breu.io/quantm/internal/hooks/github/workflows"
	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
	"go.breu.io/quantm/internal/proto/hooks/github/v1/githubv1connect"
)

type (
	GithubService struct {
		githubv1connect.UnimplementedGithubServiceHandler
	}
)

func (s *GithubService) Install(
	ctx context.Context, req *connect.Request[githubv1.InstallRequest],
) (*connect.Response[emptypb.Empty], error) {
	if req.Msg.Action != githubv1.SetupAction_INSTALL {
		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	opts := defs.NewInstallWorkflowOptions(req.Msg.InstallationId, req.Msg.Action)
	args := defs.RequestInstall{
		InstallationID: req.Msg.InstallationId,
		SetupAction:    req.Msg.Action,
		OrgID:          uuid.MustParse(req.Msg.OrgId),
	}

	_, err := durable.OnHooks().SignalWithStartWorkflow(ctx, opts, defs.SignalRequestInstall, args, workflows.Install)
	if err != nil {
		return nil, erratic.NewInternalServerError().AddHint("reason", "unable to schedule workflow").ToConnectError()
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewGithubServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return githubv1connect.NewGithubServiceHandler(&GithubService{}, opts...)
}
