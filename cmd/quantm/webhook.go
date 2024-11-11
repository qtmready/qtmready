package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/hooks/github"
)

type (
	// WebhookService is a webserver to manage webhooks for all the hooks. It conforms to the graceful.Service
	// interface, allowing for graceful start and shutdown. It wraps echo.Echo to provide this functionality.
	WebhookService struct {
		*echo.Echo
	}
)

func (w *WebhookService) Start(ctx context.Context) error {
	slog.Info("webhook: starting ...", "port", 8000)

	err := w.Echo.Start(":8000")
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (w *WebhookService) Stop(ctx context.Context) error {
	return w.Echo.Shutdown(ctx)
}

func NewWebhookServer() *WebhookService {
	webhook := echo.New()
	webhook.HideBanner = true
	webhook.HidePort = true

	github := &github.Webhook{}

	webhook.POST("/webhooks/github", github.Handler)

	return &WebhookService{webhook}
}
