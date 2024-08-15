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

// Service returns the global service instance. If the global service instance has not been initialized, it will be initialized with
// default values. The benefit is, you don't need to initialize the service instance in main.go if you don't need to override the
// default values.
func Service() service.Service {
	svcOnce.Do(func() {
		svc = service.New(
			service.FromEnvironment(),
			service.WithVersionFromBuildInfo(),
		)
	})

	return svc
}

// Logger returns the global structured logger. If the global structured logger has not been initialized, it will be initialized with
// default values. This is required if we need to pass the logger to other packages during initialization.
// NOTE: Do not use this for logging. Use slog.Info, slog.Warn, slog.Error, etc. instead.
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
			temporal.WithQueue(CoreQueue),
			temporal.WithQueue(ProvidersQueue),
			temporal.WithQueue(MutexQueue), // FIXME: WithClientCreation needs to come before the queue.
			temporal.WithQueue(WebSocketQueue),
		)
	})

	return tmprl
}

// InitServiceForTest initializes the global service instance for testing.
func InitServiceForTest() {
	svc = service.New(
		service.WithName("test"),
		service.WithDebug(true),
		service.WithSecret(password.MustGenerate(32, 8, 0, false, false)),
	)
}

func CLI() cli.Cli {
	cliOnce.Do(func() {
		ci = cli.New(
			cli.FromEnvironment(),
		)
	})

	return ci
}
