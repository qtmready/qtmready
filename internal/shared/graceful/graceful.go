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

// Package graceful provides a robust mechanism for gracefully shutting down goroutines and handling potential errors
// during their initialization. This ensures smooth program termination, guaranteeing all in-flight requests are
// processed and ongoing tasks are completed before exit.
//
// The package offers two key functions, working in concert to achieve graceful shutdown:
//
//   - graceful.Go: Launches a goroutine and forwards any encountered initialization errors to the provided error
//     channel.
//   - graceful.Shutdown: Executes user-defined cleanup functions and manages the shutdown process. Optionally signals
//     other programs to gracefully terminate and includes a timeout parameter to handle processes that might get stuck
//     during shutdown, allowing for forceful termination if necessary.
//
// The package also provides two helper functions to make it easier to use graceful.Go with different types of
// functions:
//
//   - GrabAndGo:  Creates a function that can be launched using graceful.Go, accepting a parameter. It simplifies
//     starting functions that accept a single parameter.
//   - StopAndDrop: Creates a function that can be launched using graceful.Go, designed for programs like Temporal that
//     utilize an interrupt channel for graceful shutdown.
//
// Example Usage:
//
//	import (
//	  "context"
//	  "os"
//	  "os/signal"
//	  "syscall"
//	  "time"
//
//	  "github.com/labstack/echo/v4"
//	  "go.breu.io/quantm/internal/graceful"
//	  "go.breu.io/quantm/internal/shared"
//	)
//
//	// Define cleanup functions
//	func shutdownDatabase(ctx context.Context) error {
//	  // ... perform database shutdown actions ...
//	  return nil
//	}
//
//	func closeConnections(ctx context.Context) error {
//	  // ... close network connections ...
//	  return nil
//	}
//
//	func main() {
//	  ctx, cancel := context.WithCancel(context.Background())
//	  defer cancel()
//
//	  quit := make(chan error)
//	  interrupt := make(chan any)
//
//	  // Handle termination signals (SIGINT, SIGTERM, SIGQUIT)
//	  sigterm := make(chan os.Signal, 1)
//	  signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
//
//	  // Start the Echo server:
//	  graceful.Go(ctx, graceful.GrabAndGo(ctx, echo.New().Start, ":8080"), quit)
//
//	  // Run a Temporal worker:
//	  graceful.Go(ctx, graceful.StopAndDrop(ctx, worker.Run, interrupt), quit)
//
//	  // Wait for a signal or an error
//	  select {
//	  case <-sigterm:
//	    slog.Info("shutdown signal received, gracefully shutting down all connections")
//	  case err := <-errs:
//	    if err != nil {
//	      slog.Error("error received from goroutine", "error", err)
//	    }
//	  }
//
//	  // Gracefully shutdown components
//	  cleanups := []graceful.Cleanup{
//	    shutdownDatabase,
//	    closeConnections,
//	  }
//	  code := graceful.Shutdown(ctx, cleanups, quit, 10*time.Second, 0)
//	  os.Exit(code)
//	}
package graceful

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type (
	// Cleanup represents a function that performs cleanup actions during a graceful shutdown.
	Cleanup func(ctx context.Context) error

	// Parameterized represents a function that can be started, typically accepting a param.
	Parameterized[T any] func(param T) error

	// Interruptable represents a function that can be gracefully interrupted.
	Interruptable func(interrupt <-chan any) error
)

// GrabAndGo simplifies the use of graceful.Go with functions that require an argument.
func GrabAndGo[T any](fn Parameterized[T], arg T) func() error {
	return func() error {
		return fn(arg)
	}
}

// StopAndDrop simplifies the use of graceful.Go with functions that accept an interrupt channel for graceful shutdown.
func StopAndDrop(fn Interruptable, interrupt <-chan any) func() error {
	return func() error {
		return fn(interrupt)
	}
}

// Go runs a function in a goroutine and sends any errors to the quit channel.
//
// The Go function takes a context, a function to execute, and a channel to send errors to. It runs the function in a
// goroutine and sends any errors to the quit channel.
//
// It is intended to be used in conjunction with the Shutdown function to handle errors from goroutines and ensure a
// graceful shutdown.
func Go(ctx context.Context, fn func() error, quit chan error) {
	go func() {
		if err := fn(); err != nil {
			quit <- err
		}
	}()
}

// Shutdown handles the graceful shutdown process for the given components. enabling components to drain inflight
// requests and complete ongoing tasks before exiting.
//
// The Shutdown function gracefully shuts down components by:
//
//  1. Sending a shutdown signal to the quit channel.
//  2. Calling each shutdown handler in the handlers slice in a separate goroutine.
//  3. Waiting for all handlers to complete before exiting.
//
// It is intended to be used in conjunction with the Go function to handle errors from goroutines and ensure a graceful
// shutdown.
func Shutdown(ctx context.Context, cleanups []Cleanup, quit chan bool, timeout time.Duration, code int) int {
	quit <- true

	var wg sync.WaitGroup

	wg.Add(len(cleanups))

	for _, cleanup := range cleanups {
		go func() {
			defer wg.Done()

			if err := cleanup(ctx); err != nil {
				slog.Error("unable to shutdown gracefully", "error", err)

				code = 1
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All cleanups completed within the timeout
	case <-time.After(timeout):
		slog.Warn("shutdown timeout reached, some cleanups may not have completed")
	}

	return code
}
