// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  APPLICABLE LAW.

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
