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
	"log/slog"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/web"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/logger"
)

// configure_metrics configures the `/metrics` endpoint for exporting metrics using the `echoprometheus` package.
// It can be used to configure the endpoint on any provided Echo instance. However, we recommend using a separate
// Echo instance dedicated to serving metrics on port 9090. This approach provides several benefits:
//
//   - Security: Separating the metrics server from the main application hides the metrics endpoint from the public.
//   - Simplicity:  Using the default Prometheus port (9090) simplifies monitoring, as it's widely recognized and used by
//     various monitoring tools.
func configure_metrics(metrics *echo.Echo) {
	slog.Info("metrics: configuring ...")

	metrics.HideBanner = true
	metrics.HidePort = true

	slog.Info("metrics: enabling prometheus handler ...")

	metrics.GET("/metrics", echoprometheus.NewHandler())
}

// configure_api configures the Echo instance, along with all the middlewares and routes to serve the API.
func configure_api(server *echo.Echo) {
	slog.Info("api: configuring ...")

	server.HideBanner = true
	server.HidePort = true
	server.HTTPErrorHandler = shared.EchoAPIErrorHandler
	server.Validator = &shared.EchoValidator{Validator: shared.Validator()}

	slog.Info("api: registering middlewares ...")

	server.Use(middleware.CORS())
	server.Use(otelecho.Middleware(shared.Service().GetName()))
	server.Use(logger.NewRequestLoggerMiddleware())
	server.Use(echoprometheus.NewMiddleware(shared.Service().GetName()))
	server.Use(middleware.Recover())

	slog.Info("api: registering default service health check ...")

	server.GET("/healthx", healthz)

	slog.Info("api: registering routes ...")

	auth.RegisterHandlers(server, auth.NewServerHandler(auth.Middleware))
	web.RegisterHandlers(server, web.NewServerHandler(auth.Middleware))
	github.RegisterHandlers(server, github.NewServerHandler(auth.Middleware))
	slack.RegisterHandlers(server, slack.NewServerHandler(auth.Middleware))
}
