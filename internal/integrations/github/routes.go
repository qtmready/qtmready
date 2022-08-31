package github

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.breu.io/ctrlplane/internal/shared"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	g.POST("/webhook", webhook)

	// protected routes
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
	handlers := WebhookEventHandlers{
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

	// TODO: abstract this out to a common function
	claims := ctx.Get("user").(*jwt.Token).Claims.(*shared.JWTClaims)
	shared.Logger.Info("claims", "claims", claims)
	teamID, err := gocql.ParseUUID(claims.TeamID)
	if err != nil {
		shared.Logger.Error("error", "error", err)
		return err
	}

	payload := &CompleteInstallationSignalPayload{request.InstallationID, request.SetupAction, teamID}

	workflows := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.IntegrationsQueue].
		GetWorkflowOptions("github", strconv.Itoa(int(payload.InstallationID)), string(InstallationEvent))

	run, err := shared.Temporal.Client.
		SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			CompleteInstallationSignal.String(),
			payload,
			opts,
			workflows.OnInstall,
		)

	if err != nil {
		shared.Logger.Error("error", "error", err)
		return err
	}

	return ctx.JSON(http.StatusOK, run.GetRunID())
}
