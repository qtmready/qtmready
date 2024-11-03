package main

import (
	"context"

	"github.com/labstack/echo/v4"

	githubweb "go.breu.io/quantm/internal/hooks/github/web"
)

type (
	WebhookService struct {
		*echo.Echo
	}
)

func (w *WebhookService) Start(ctx context.Context) error {
	return w.Echo.Start(":8000")
}

func (w *WebhookService) Stop(ctx context.Context) error {
	return w.Echo.Shutdown(ctx)
}

func NewWebhookServer() *WebhookService {
	webhook := echo.New()
	github := &githubweb.Webhook{}

	webhook.POST("/webhooks/github", github.Handler)

	return &WebhookService{webhook}
}
