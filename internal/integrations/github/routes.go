package github

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.breu.io/ctrlplane/internal/cmn"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	g.POST("/webhook", webhook)
	g.Use(middlewares...)
	g.POST("/complete-installation", completeInstallation)
}

func webhook(ctx echo.Context) error {
	signature := ctx.Request().Header.Get("X-Hub-Signature")
	if signature == "" {
		return ctx.JSON(http.StatusUnauthorized, ErrorMissingHeaderGithubSignature)
	}

	// NOTE: We are reading the request body twice. This is not ideal.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		return ctx.JSON(http.StatusUnauthorized, err)
	}

	headerEvent := ctx.Request().Header.Get("X-GitHub-Event")
	if headerEvent == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorMissingHeaderGithubEvent)
	}

	event := WebhookEvent(headerEvent)
	// A Map of event types to their respective handlers
	handlers := EventHandlers{
		InstallationEvent: handleInstallationEvent,
		PushEvent:         handlePushEvent,
	}

	if handle, exists := handlers[event]; exists {
		return handle(ctx)
	} else {
		return ctx.JSON(http.StatusBadRequest, ErrorInvalidEvent)
	}
}

func completeInstallation(ctx echo.Context) error {
	request := &CompleteInstallationRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, err := gocql.ParseUUID(ctx.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["team_id"].(string))
	if err != nil {
		return err
	}

	payload := &CompleteInstallationPayload{request.InstallationID, teamID}
	workflows := &Workflows{}
	opts := cmn.Temporal.
		Queues[cmn.GithubIntegrationQueue].
		GetWorkflowOptions(strconv.Itoa(int(payload.InstallationID)), string(InstallationEvent))

	run, err := cmn.Temporal.Client.
		SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			CompleteInstallationSignal.String(),
			payload,
			opts,
			workflows.OnInstall,
			payload,
		)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, run.GetRunID())
}
