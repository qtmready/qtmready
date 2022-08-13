package github

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	g.POST("/webhook", webhook)
}

func webhook(ctx echo.Context) error {
	signature := ctx.Request().Header.Get("X-Hub-Signature")
	if signature == "" {
		return ctx.JSON(http.StatusUnauthorized, ErrorMissingHeaderGithubSignature)
	}

	// NOTE: Multiple concerns are involved here.
	// 1. We are reading the request twice (once to get the body, once to verify the signature).
	// 2. In order to do this .. we use NopCloser to prevent the body from being closed.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		return ctx.JSON(http.StatusUnauthorized, err)
	}

	headerEvent := ctx.Request().Header.Get("X-GitHub-Event")
	if headerEvent == "" {
		return ctx.JSON(http.StatusBadRequest, ErrorMissingHeaderGithubEvent)
	}

	event := WebhookEvent(headerEvent)
	if handle, exists := eventHandlers[event]; exists {
		return handle(ctx)
	} else {
		return ctx.JSON(http.StatusBadRequest, ErrorInvalidEvent)
	}
}
