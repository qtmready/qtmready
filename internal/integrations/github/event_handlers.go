package github

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.breu.io/ctrlplane/internal/cmn"
	"go.uber.org/zap"
)

type eventHandler func(ctx echo.Context) error

var w *Workflows

// A Map of event types to their respective handlers
var eventHandlers = map[WebhookEvent]eventHandler{
	InstallationEvent: handleInstallationEvent,
	PushEvent:         handlePushEvent,
}

// handles GitHub installation event
func handleInstallationEvent(ctx echo.Context) error {
	payload := InstallationEventPayload{}
	if err := ctx.Bind(&payload); err != nil {
		return err
	}

	opts := cmn.Temporal.
		Queues[cmn.GithubIntegrationQueue].
		CreateWorkflowOptions(strconv.Itoa(int(payload.Installation.ID)), string(InstallationEvent))

	exe, err := cmn.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnInstall, payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, exe.GetRunID())
}

// handles GitHub push event
func handlePushEvent(ctx echo.Context) error {
	payload := PushEventPayload{}
	if err := ctx.Bind(&payload); err != nil {
		cmn.Log.Info("Error: ", zap.Any("body", ctx.Request().Body), zap.String("error", err.Error()))
		return err
	}

	opts := cmn.Temporal.
		Queues[cmn.GithubIntegrationQueue].
		CreateWorkflowOptions(strconv.Itoa(int(payload.Installation.ID)), string(PushEvent), "ref", payload.Ref)

	exe, err := cmn.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPush, payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, exe.GetRunID())
}
