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

	linkTo, err := uuid.Parse(reqst.Msg.GetLinkTo())
	if err != nil {
		return nil, erratic.NewBadRequestError(erratic.HooksSlackModule).
			WithReason("invalid link_to UUID").Wrap(err)
	}

	message, err := db.Queries().GetChatLink(ctx, linkTo)
	if err != nil {
		return nil, erratic.NewDatabaseError(erratic.HooksSlackModule).
			WithReason("failed to query message by link_to").Wrap(err)
	}

	if message.ID != uuid.Nil {
		return nil, erratic.NewExistsError(erratic.HooksSlackModule).
			WithReason("message with link_to already exists")
	}

	if reqst.Msg.GetCode() == "" {
		return nil, erratic.NewBadRequestError(erratic.HooksSlackModule).WithReason("missing OAuth code")
	}

	response, err := slack.GetOAuthV2Response(&c, config.ClientID(), config.ClientSecret(), reqst.Msg.GetCode(), config.ClientRedirectURL())
	if err != nil {
		return nil, erratic.NewNetworkError(erratic.HooksSlackModule).
			WithReason("failed to get OAuth response from Slack").Wrap(err)
	}

	if response.AuthedUser.AccessToken != "" {
		if err := _user(ctx, reqst, response); err != nil {
			return nil, erratic.NewSystemError(erratic.HooksSlackModule).
				WithReason("failed to process user OAuth").Wrap(err) // More specific reason if possible
		}

		return connect.NewResponse(&emptypb.Empty{}), nil
	}

	if err := _bot(ctx, reqst, response); err != nil {
		return nil, erratic.NewSystemError(erratic.HooksSlackModule).
			WithReason("failed to process bot OAuth").Wrap(err) // More specific reason if possible
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func NewSlackServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return slackv1connect.NewSlackServiceHandler(&SlackService{}, opts...)
}
