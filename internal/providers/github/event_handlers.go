// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package github

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/shared"
)

func handleInstallationEvent(ctx echo.Context) error {
	payload := &InstallationEventPayload{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	shared.Logger.Info("installation event received ...")

	workflows := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.ProvidersQueue].
		GetWorkflowOptions("github", strconv.FormatInt(payload.Installation.ID, 10), string(InstallationEvent))

	exe, err := shared.Temporal.Client.SignalWithStartWorkflow(
		ctx.Request().Context(),
		opts.ID,
		InstallationEventSignal.String(),
		payload,
		opts,
		workflows.OnInstall,
	)
	if err != nil {
		shared.Logger.Error("unable to signal ...", "options", opts, "error", err)
		return nil
	}

	shared.Logger.Info("installation event handled ...", "options", opts, "execution", exe.GetRunID())

	return ctx.JSON(http.StatusCreated, exe.GetRunID())
}

// handles GitHub push event.
func handlePushEvent(ctx echo.Context) error {
	payload := PushEventPayload{}
	if err := ctx.Bind(&payload); err != nil {
		return err
	}

	w := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.ProvidersQueue].
		GetWorkflowOptions("github", strconv.FormatInt(payload.Installation.ID, 10), PushEvent.String(), "ref", payload.After)

	exe, err := shared.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPush, payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, exe.GetRunID())
}
