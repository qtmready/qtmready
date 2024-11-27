package nomad

import (
	"context"
	"encoding/base64"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	pkg_slack "go.breu.io/quantm/internal/hooks/slack"
	"go.breu.io/quantm/internal/hooks/slack/config"
	"go.breu.io/quantm/internal/hooks/slack/defs"
	"go.breu.io/quantm/internal/hooks/slack/errors"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	slackv1 "go.breu.io/quantm/internal/proto/hooks/slack/v1"
	"go.breu.io/quantm/internal/proto/hooks/slack/v1/slackv1connect"
)

type (
	SlackService struct {
		slackv1connect.UnimplementedSlackServiceHandler
	}
)

func (s *SlackService) SlackOauth(
	ctx context.Context, reqst *connect.Request[slackv1.SlackOauthRequest],
) (*connect.Response[emptypb.Empty], error) {
	var c pkg_slack.HTTPClient

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

func _user(
	ctx context.Context, reqst *connect.Request[slackv1.SlackOauthRequest], response *slack.OAuthV2Response,
) error {
	client, _ := config.GetSlackClient(response.AuthedUser.AccessToken)

	identity, err := client.GetUserIdentity()
	if err != nil {
		return err
	}

	// Generate a key for AES-256.
	key := pkg_slack.Generate(response.Team.ID)

	// Encrypt the user access token.
	user_token, err := pkg_slack.Encrypt([]byte(response.AuthedUser.AccessToken), key)
	if err != nil {
		return err
	}

	// Encrypt the bot access token.
	bot_token, err := pkg_slack.Encrypt([]byte(response.AccessToken), key)
	if err != nil {
		return err
	}

	slack_user := &defs.MessageProviderSlackUserInfo{
		BotToken:       base64.StdEncoding.EncodeToString(bot_token),
		UserToken:      base64.StdEncoding.EncodeToString(user_token),
		ProviderUserID: identity.User.ID,
		ProviderTeamID: identity.Team.ID,
	}

	data, err := slack_user.Marshal()
	if err != nil {
		return err
	}

	// Convert the string to uuid.UUID
	// TODO - find better approach
	link_to, err := uuid.Parse(reqst.Msg.GetLinkTo())
	if err != nil {
		return err
	}

	// save messaging
	m := entities.CreateMessagingParams{
		Hook:   int32(eventsv1.MessagingHook_MESSAGING_HOOK_SLACK),
		Kind:   "user",
		LinkTo: link_to,
		Data:   data,
	}

	_, err = db.Queries().CreateMessaging(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func _bot(
	ctx context.Context, reqst *connect.Request[slackv1.SlackOauthRequest], response *slack.OAuthV2Response,
) error {
	// Generate a key for AES-256.
	key := pkg_slack.Generate(response.Team.ID)

	// Encrypt the bot access token.
	bot_token, err := pkg_slack.Encrypt([]byte(response.AccessToken), key)
	if err != nil {
		return err
	}

	slack_bot := &defs.MessageProviderSlackData{
		ChannelID:     response.IncomingWebhook.ChannelID,
		ChannelName:   response.IncomingWebhook.Channel,
		WorkspaceName: response.Team.Name,
		WorkspaceID:   response.Team.ID,
		BotToken:      base64.StdEncoding.EncodeToString(bot_token), // Store the base64-encoded encrypted token
	}

	data, err := slack_bot.Marshal()
	if err != nil {
		return err
	}

	// Convert the string to uuid.UUID
	// TODO - find better approach
	link_to, err := uuid.Parse(reqst.Msg.GetLinkTo())
	if err != nil {
		return err
	}

	// save messaging
	m := entities.CreateMessagingParams{
		Hook:   int32(eventsv1.MessagingHook_MESSAGING_HOOK_SLACK),
		Kind:   "bot",
		LinkTo: link_to,
		Data:   data,
	}

	_, err = db.Queries().CreateMessaging(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func NewSlackServiceHandler(opts ...connect.HandlerOption) (string, http.Handler) {
	return slackv1connect.NewSlackServiceHandler(&SlackService{}, opts...)
}
