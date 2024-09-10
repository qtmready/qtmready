package slack

import (
	"encoding/base64"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	ServerHandler struct{ *auth.SecurityHandler }
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(middleware echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &auth.SecurityHandler{Middleware: middleware},
	}
}

func (e *ServerHandler) SlackOauth(ctx echo.Context, params SlackOauthParams) error {
	var c HTTPClient

	code := ctx.QueryParam("code")
	if code == "" {
		return shared.NewAPIError(http.StatusNotFound, ErrCodeEmpty)
	}

	response, err := slack.GetOAuthV2Response(&c, ClientID(), ClientSecret(), code, ClientRedirectURL())
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if response.AuthedUser.AccessToken != "" {
		// linked message provider user (slack) to quantm user
		_, err := _user(ctx, response)
		if err != nil {
			return shared.NewAPIError(http.StatusBadRequest, err)
		}

		return ctx.JSON(http.StatusOK, nil)
	}

	// connect slack bot user with channel info
	resp, err := _bot(response)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func _user(ctx echo.Context, response *slack.OAuthV2Response) (*auth.MessageProviderUserInfo, error) {
	userID, _ := gocql.ParseUUID(ctx.Get("user_id").(string))

	teamuser := &auth.TeamUser{}
	query := db.QueryParams{"user_id": userID.String()}

	if err := db.Get(teamuser, query); err != nil {
		return nil, err
	}

	client, _ := instance.GetSlackClient(response.AuthedUser.AccessToken)

	identity, err := client.GetUserIdentity()
	if err != nil {
		shared.Logger().Error("SlackOauth/identity", "error", err.Error())
		return nil, err
	}

	// Generate a key for AES-256.
	key := generateKey(response.Team.ID)

	// Encrypt the user access token.
	encryptedUserToken, err := encrypt([]byte(response.AuthedUser.AccessToken), key)
	if err != nil {
		return nil, err
	}

	// Encrypt the bot access token.
	encryptedBotToken, err := encrypt([]byte(response.AccessToken), key)
	if err != nil {
		return nil, err
	}

	userinfo := auth.MessageProviderUserInfo{
		Slack: &auth.MessageProviderSlackUserInfo{
			BotToken:       base64.StdEncoding.EncodeToString(encryptedBotToken),
			UserToken:      base64.StdEncoding.EncodeToString(encryptedUserToken),
			ProviderUserID: identity.User.ID,
			ProviderTeamID: identity.Team.ID,
		},
	}

	teamuser.IsMessageProviderLinked = true
	teamuser.MessageProvider = auth.MessageProviderSlack
	teamuser.MessageProviderUserInfo = userinfo

	if err := db.Save(teamuser); err != nil {
		return nil, err
	}

	return &userinfo, nil
}

func _bot(response *slack.OAuthV2Response) (*defs.MessageProviderData, error) {
	// Generate a key for AES-256.
	key := generateKey(response.Team.ID)

	// Encrypt the access token.
	encryptedToken, err := encrypt([]byte(response.AccessToken), key)
	if err != nil {
		return nil, err
	}

	resp := &defs.MessageProviderData{
		Slack: &defs.MessageProviderSlackData{
			ChannelID:     response.IncomingWebhook.ChannelID,
			ChannelName:   response.IncomingWebhook.Channel,
			WorkspaceName: response.Team.Name,
			WorkspaceID:   response.Team.ID,
			BotToken:      base64.StdEncoding.EncodeToString(encryptedToken), // Store the base64-encoded encrypted token
		},
	}

	return resp, nil
}
