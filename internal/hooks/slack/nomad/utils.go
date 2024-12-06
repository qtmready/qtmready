package nomad

import (
	"context"
	"encoding/base64"

	"connectrpc.com/connect"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/hooks/slack/config"
	"go.breu.io/quantm/internal/hooks/slack/defs"
	"go.breu.io/quantm/internal/hooks/slack/fns"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	slackv1 "go.breu.io/quantm/internal/proto/hooks/slack/v1"
	"go.breu.io/quantm/internal/utils"
)

func _user(
	ctx context.Context, reqst *connect.Request[slackv1.OauthRequest], response *slack.OAuthV2Response,
) error {
	client, _ := config.GetSlackClient(response.AuthedUser.AccessToken)

	identity, err := client.GetUserIdentity()
	if err != nil {
		return err
	}

	// Generate a key for AES-256.
	key := fns.Generate(response.Team.ID)

	// Encrypt the user access token.
	user_token, err := fns.Encrypt([]byte(response.AuthedUser.AccessToken), key)
	if err != nil {
		return err
	}

	// Encrypt the bot access token.
	bot_token, err := fns.Encrypt([]byte(response.AccessToken), key)
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
	link_to, err := utils.ParseUUID(reqst.Msg.GetLinkTo())
	if err != nil {
		return err
	}

	// save messaging
	m := entities.CreateChatLinkParams{
		Hook:   int32(eventsv1.ChatHook_CHAT_HOOK_SLACK),
		Kind:   defs.KindUser,
		LinkTo: link_to,
		Data:   data,
	}

	_, err = db.Queries().CreateChatLink(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func _bot(
	ctx context.Context, reqst *connect.Request[slackv1.OauthRequest], response *slack.OAuthV2Response,
) error {
	// Generate a key for AES-256.
	key := fns.Generate(response.Team.ID)

	// Encrypt the bot access token.
	bot_token, err := fns.Encrypt([]byte(response.AccessToken), key)
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
	link_to, err := utils.ParseUUID(reqst.Msg.GetLinkTo())
	if err != nil {
		return err
	}

	// save messaging
	m := entities.CreateChatLinkParams{
		Hook:   int32(eventsv1.ChatHook_CHAT_HOOK_SLACK),
		Kind:   defs.KindBot,
		LinkTo: link_to,
		Data:   data,
	}

	_, err = db.Queries().CreateChatLink(ctx, m)
	if err != nil {
		return err
	}

	return nil
}
