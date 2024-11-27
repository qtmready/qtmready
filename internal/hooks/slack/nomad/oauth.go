package nomad

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/hooks/slack/config"
	"go.breu.io/quantm/internal/hooks/slack/errors"
	"go.breu.io/quantm/internal/hooks/slack/fns"
	slackv1 "go.breu.io/quantm/internal/proto/hooks/slack/v1"
	"go.breu.io/quantm/internal/proto/hooks/slack/v1/slackv1connect"
)

type (
	SlackService struct {
		slackv1connect.UnimplementedSlackServiceHandler
	}
)

func (s *SlackService) Oauth(
	ctx context.Context, reqst *connect.Request[slackv1.OauthRequest],
) (*connect.Response[emptypb.Empty], error) {
	var c fns.HTTPClient

	// check the already exist record against the link_to
	// if exist return the error already exit
	link_to, err := uuid.Parse(reqst.Msg.GetLinkTo())
	if err != nil {
		return nil, err
	}

	message, err := db.Queries().GetMessagesByLinkTo(ctx, link_to)
	if message.ID != uuid.Nil {
		return nil, erratic.NewInternalServerError().AddHint("reason", errors.ErrRecordExist.Error()).ToConnectError()
	}

	if reqst.Msg.GetCode() == "" {
		return nil, erratic.NewInternalServerError().AddHint("reason", errors.ErrCodeEmpty.Error()).ToConnectError()
	}

	response, err := slack.
		GetOAuthV2Response(&c, config.ClientID(), config.ClientSecret(), reqst.Msg.GetCode(), config.ClientRedirectURL())
	if err != nil {
		return nil, erratic.NewInternalServerError().AddHint("reason", err.Error()).ToConnectError()
	}

	if response.AuthedUser.AccessToken != "" {
		// linked message provider user (slack) to quantm user
		if err := _user(ctx, reqst, response); err != nil {
			return nil, erratic.NewInternalServerError().AddHint("reason", err.Error()).ToConnectError()
		}

		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	// connect slack bot user with channel info
	if err := _bot(ctx, reqst, response); err != nil {
		return nil, erratic.NewInternalServerError().AddHint("reason", err.Error()).ToConnectError()
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewSlackServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return slackv1connect.NewSlackServiceHandler(&SlackService{}, opts...)
}
