package authnmd

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	authv1 "go.breu.io/quantm/internal/proto/ctrlplane/auth/v1"
	"go.breu.io/quantm/internal/proto/ctrlplane/auth/v1/authv1connect"
)

type (
	OrgService struct {
		authv1connect.UnimplementedOrgServiceHandler
	}
)

func (s *OrgService) SetOrgHooks(
	ctx context.Context, req *connect.Request[authv1.SetOrgHooksRequest],
) (*connect.Response[emptypb.Empty], error) {
	hooks, err := json.Marshal(req.Msg.Hooks)
	if err != nil {
		slog.Info(err.Error(), "hooks", req.Msg.Hooks)

		return nil, erratic.NewInternalServerError("error", err.Error()).ToConnectError()
	}

	params := entities.SetOrgHooksParams{ID: uuid.MustParse(req.Msg.GetOrgId()), Hooks: hooks}

	err = db.Queries().SetOrgHooks(ctx, params)
	if err != nil {
		return nil, erratic.NewInternalServerError().Wrap(err).ToConnectError()
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewOrgServiceServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return authv1connect.NewOrgServiceHandler(
		&OrgService{},
		opts...,
	)
}
