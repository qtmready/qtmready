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

type (
	ServerHandler struct{}
)

func (s *ServerHandler) GithubCompleteInstallation(ctx echo.Context) error {
	request := &CompleteInstallationRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, err := gocql.ParseUUID(ctx.Get("team_id").(string))

	if err != nil {
		return err
	}

	payload := &CompleteInstallationSignalPayload{request.InstallationId, request.SetupAction, teamID}

	workflows := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.ProvidersQueue].
		GetWorkflowOptions("github", strconv.Itoa(int(payload.InstallationID)), WebhookEventInstallation.String())

	exe, err := shared.Temporal.Client.SignalWithStartWorkflow(
		ctx.Request().Context(),
		opts.ID,
		WorkflowSignalCompleteInstallation.String(),
		payload,
		opts,
		workflows.OnInstall,
	)
	if err != nil {
		shared.Logger.Error("error", "error", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunId: exe.GetID(), Status: WorkflowStatusQueued})
}

func (s *ServerHandler) GithubGetRepos(ctx echo.Context) error {
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

func (s *ServerHandler) GithubWebhook(ctx echo.Context) error {
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
	handlers := WebhookEventHandlers{
		WebhookEventInstallation: handleInstallationEvent,
		WebhookEventPush:         handlePushEvent,
		WebhookEventPullRequest:  handlePullRequestEvent,
	}

	if handle, exists := handlers[event]; exists {
		return handle(ctx)
	}

	return ctx.JSON(http.StatusBadRequest, ErrInvalidEvent)
}
