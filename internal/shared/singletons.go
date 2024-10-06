// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2023, 2024.
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

package shared

import (
	"log/slog"
	"os"
	"reflect"
	"strings"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"github.com/go-playground/validator/v10"
	"github.com/sethvargo/go-password/password"

	"go.breu.io/quantm/internal/shared/cli"
	"go.breu.io/quantm/internal/shared/logger"
	"go.breu.io/quantm/internal/shared/service"
	"go.breu.io/quantm/internal/shared/temporal"
)

var (
	svc     service.Service // Global service instance.
	svcOnce sync.Once       // Global service instance initializer.

	lgr     *slog.Logger // Global logger instance.
	lgrOnce sync.Once    // Global logger instance initializer

	vld     *validator.Validate // Global validator instance.
	vldOnce sync.Once           // Global validator instance initializer

	tmprl     temporal.Temporal // Global temporal instance.
	tmprlOnce sync.Once         // Global temporal instance initializer

	ci      cli.Cli   // Global cli instance
	cliOnce sync.Once // Global cli instance initializer
)

// Service returns the global service instance.
//
// If the global service instance has not been initialized, it will be initialized with default values. The benefit is,
// you don't need to initialize the service instance in `main.go` if you don't need to override the default values.
func Service() service.Service {
	svcOnce.Do(func() {
		svc = service.New(
			service.FromEnvironment(),
			service.WithVersionFromBuildInfo(),
		)
	})

	return svc
}

// Logger returns the global structured logger.
//
// If the global structured logger has not been initialized, it will be initialized with default values. This is required
// if we need to pass the logger to other packages during initialization.
//
// Deprecated: Do not use this for logging. Use `slog.Info`, `slog.Warn`, `slog.Error`, etc. instead.
func Logger() *slog.Logger {
	var handler slog.Handler

	debug := Service().GetDebug()
	level := slog.LevelInfo
	opts := &slog.HandlerOptions{AddSource: !debug}
	gcp := metadata.OnGCE()

	lgrOnce.Do(func() {
		switch {
		case debug:
			level = slog.LevelDebug
			opts.Level = level
			handler = slog.NewTextHandler(os.Stdout, opts)

		case gcp:
			handler = logger.NewGoogleCloudHandler(os.Stdout, opts)

		default:
			handler = slog.NewJSONHandler(os.Stdout, opts)
		}

		lgr = slog.New(handler)
	})

	slog.SetDefault(lgr)

	return lgr
}

// Validator returns the global validator instance.
func Validator() *validator.Validate {
	vldOnce.Do(func() {
		vld = validator.New()
		vld.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}

			return name
		})
	})

	return vld
}

// Temporal returns the global temporal instance.
func Temporal() temporal.Temporal {
	tmprlOnce.Do(func() {
		tmprl = temporal.New(
			temporal.FromEnvironment(),
			temporal.WithLogger(Logger()),
			temporal.WithClientCreation(),
		)
	})

	return tmprl
}

// InitServiceForTest initializes the global service instance for testing.
func InitServiceForTest() {
	svcOnce.Do(func() {
		secret := password.MustGenerate(32, 8, 0, false, false)[0:32]
		svc = service.New(
			service.WithName("test"),
			service.WithDebug(true),
			service.WithSecret(secret),
		)
	})
}

// CLI returns the global CLI instance.
func CLI() cli.Cli {
	cliOnce.Do(func() {
		ci = cli.New(
			cli.FromEnvironment(),
		)
	})

	return ci
}
