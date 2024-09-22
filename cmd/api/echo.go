package main

import (
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

func configure_metrics(metrics *echo.Echo) {
	metrics.HideBanner = true
	metrics.HidePort = true

	metrics.GET("/metrics", echoprometheus.NewHandler())
}

func configure_web(server *echo.Echo) {
	server.HideBanner = true
	server.HidePort = true
	server.HTTPErrorHandler = shared.EchoAPIErrorHandler
	server.Validator = &shared.EchoValidator{Validator: shared.Validator()}

	server.Use(middleware.CORS())
	server.Use(otelecho.Middleware(shared.Service().GetName()))
	server.Use(logger.NewRequestLoggerMiddleware())
	server.Use(echoprometheus.NewMiddleware(shared.Service().GetName()))
	server.Use(middleware.Recover())

	server.GET("/healthx", healthz)

	auth.RegisterHandlers(server, auth.NewServerHandler(auth.Middleware))
	web.RegisterHandlers(server, web.NewServerHandler(auth.Middleware))
	github.RegisterHandlers(server, github.NewServerHandler(auth.Middleware))
	slack.RegisterHandlers(server, slack.NewServerHandler(auth.Middleware))
}
