package github

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	g.POST("/webhook", webhook)

	// protected routes
	g.Use(middlewares...)
	g.POST("/complete-installation", completeInstallation)
	g.GET("/repos", repos)
}

func webhook(ctx echo.Context) error {
	shared.Logger.Debug("webhook received", "headers", ctx.Request().Header)
	signature := ctx.Request().Header.Get("X-Hub-Signature-256")

	if signature == "" {
		return ctx.JSON(http.StatusUnauthorized, ErrMissingHeaderGithubSignature)
	}

	// NOTE: We are reading the request body twice. This is not ideal.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		return ctx.JSON(http.StatusUnauthorized, err)
	}

	headerEvent := ctx.Request().Header.Get("X-GitHub-Event")
	if headerEvent == "" {
		return ctx.JSON(http.StatusBadRequest, ErrMissingHeaderGithubEvent)
	}

	event := WebhookEvent(headerEvent)
	// A Map of event types to their respective handlers
	handlers := WebhookEventHandlers{
		InstallationEvent: handleInstallationEvent,
		PushEvent:         handlePushEvent,
	}

	if handle, exists := handlers[event]; exists {
		return handle(ctx)
	}

	return ctx.JSON(http.StatusBadRequest, ErrInvalidEvent)
}

// @Summary     Completes the installation of a GitHub app.
// @Description Completes the installation of a GitHub app.
// @Tags        providers, github
// @Accept      json
// @Produce     json
// @Param       body body     CompleteInstallationRequest true "CompleteInstallationRequest"
// @Success     200  {object} WorkflowRunResponse
// @Failure     400  {object} echo.HTTPError
// @Router      /provders/github/complete-installation [post]
//
// completeInstallation completes the installation of a GitHub app.
func completeInstallation(ctx echo.Context) error {
	request := &CompleteInstallationRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, err := gocql.ParseUUID(ctx.Get("team_id").(string))

	if err != nil {
		return err
	}

	payload := &CompleteInstallationSignalPayload{request.InstallationID, request.SetupAction, teamID}

	workflows := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.ProvidersQueue].
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

	return ctx.JSON(http.StatusOK, &WorkflowRunResponse{ID: run.GetID(), RunID: run.GetRunID()})
}

// @Summary     Get GitHub repositories.
// @Description Get GitHub repositories.
// @Tags        providers, github
// @Accept      json
// @Produce     json
// @Success     200 {array}  entities.GithubRepo
// @Failure     400 {object} echo.HTTPError
// @Router      /provders/github/repos [get]
//
// repos get all the github repos for a team.
func repos(ctx echo.Context) error {
	result := make([]entities.GithubRepo, 0)
	if err := db.Filter(
		&entities.GithubRepo{},
		&result,
		db.QueryParams{"team_id": ctx.Get("team_id").(string)},
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}
