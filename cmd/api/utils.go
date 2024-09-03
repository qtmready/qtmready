// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
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
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/core/ws"
)

// _run runs a function in a goroutine.
func _run(fn func() error, ch chan error) {
	if err := fn(); err != nil {
		ch <- err
	}
}

// _serve starts the echo server in a goroutine.
func _serve(e *echo.Echo, port string) func() error {
	return func() error { return e.Start(":" + port) }
}

func _hub() error {
	worker := ws.ConnectionsHubWorker()

	return worker.Start()
}

// _graceful shuts down each goroutine gracefully.
func _graceful(ctx context.Context, fns []shutdownfn, signals []chan any, code int) {
	for _, signal := range signals {
		signal <- true
	}

	for _, fn := range fns {
		if err := fn(ctx); err != nil {
			code = 1
		}
	}

	slog.Info("shutdown complete, exiting.")

	os.Exit(code)
}
