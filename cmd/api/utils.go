// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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
