package web

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/hooks/github/config"
	"go.breu.io/quantm/internal/hooks/github/defs"
	"go.breu.io/quantm/internal/hooks/github/workflows"
	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
)

type (
	// Webhook is a Github Webhook event receiver responsible for scheduling transient workflows.
	//
	// Transient workflows gather the necessary context to formulate QuantmEvents, package them,
	// and then dispatch them to the appropriate workflow within the Quantm core for processing.
	Webhook struct{}

	// WebhookEventHandler is a function that handles Github Webhook events.
	WebhookEventHandler func(ctx echo.Context, event defs.WebhookEvent, id string) error

	// WebhookEventHandlers is a map of Github Webhook event names to their handlers.
	WebhookEventHandlers map[defs.WebhookEvent]WebhookEventHandler
)

// Handler handles Github Webhook events.
func (h *Webhook) Handler(ctx echo.Context) error {
	// Get the signature from the request header. If the signature is missing, return an unauthorized error.
	signature := ctx.Request().Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return erratic.NewUnauthorizedError().AddHint("reason", "missing signature")
	}

	// Read the request body and then reset it for subsequent use.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	// Verify the signature. Return an unauthorized error if the signature is invalid.
	if err := config.Instance().VerifyWebhookSignature(body, signature); err != nil {
		return erratic.NewUnauthorizedError().AddHint("reason", "invalid signature")
	}

	// Get the event type from the request header.
	event := defs.WebhookEvent(ctx.Request().Header.Get("X-GitHub-Event"))
	if event == defs.WebhookEventUnspecified {
		return ctx.NoContent(http.StatusNoContent)
	}

	// Get the event handler for the event type. If the event handler is not found, ignore the event.
	fn, found := h.on(event)
	if !found {
		return ctx.NoContent(http.StatusNoContent)
	}

	id := ctx.Request().Header.Get("X-GitHub-Delivery")

	// Execute the event handler.
	return fn(ctx, event, id)
}

// on returns the event handler for the given event type.
func (h *Webhook) on(event defs.WebhookEvent) (WebhookEventHandler, bool) {
	handlers := WebhookEventHandlers{
		defs.WebhookEventInstallation:             h.install,
		defs.WebhookEventInstallationRepositories: h.install_repos,
		defs.WebhookEventCreate:                   h.ref,
		defs.WebhookEventDelete:                   h.ref,
		defs.WebhookEventPush:                     h.push,
		defs.WebhookEventPullRequest:              h.pr,
	}

	fn, ok := handlers[event]

	return fn, ok
}

// install handles the installation event.
func (h *Webhook) install(ctx echo.Context, event defs.WebhookEvent, id string) error {
	payload := &defs.WebhookInstall{}
	if err := ctx.Bind(payload); err != nil {
		slog.Info("failed to bind payload", "error", err.Error())
		return erratic.NewBadRequestError("reason", "invalid payload")
	}

	action := githubv1.SetupAction_UNSPECIFIED

	switch payload.Action {
	case "created":
		action = githubv1.SetupAction_INSTALL
	case "updated":
		action = githubv1.SetupAction_UPDATE
	case "deleted":
		action = githubv1.SetupAction_DELETE
	case "new_permissions_accepted":
		action = githubv1.SetupAction_NEW_PERMISSIONS_ACCEPTED
	case "suspend":
		action = githubv1.SetupAction_SUSPEND
	case "unsuspend":
		action = githubv1.SetupAction_UNSUSPEND
	}

	if action == githubv1.SetupAction_UNSPECIFIED {
		slog.Warn("unsupported action during github install", "action", payload.Action)

		return ctx.NoContent(http.StatusNoContent)
	}

	opts := defs.NewInstallWorkflowOptions(payload.Installation.ID, action)

	_, err := durable.
		OnHooks().
		SignalWithStartWorkflow(ctx.Request().Context(), opts, defs.SignalWebhookInstall, payload, workflows.Install)
	if err != nil {
		return erratic.NewInternalServerError("reason", "failed to signal workflow", "error", err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *Webhook) install_repos(ctx echo.Context, _ defs.WebhookEvent, id string) error {
	payload := &defs.WebhookInstallRepos{}
	if err := ctx.Bind(payload); err != nil {
		return erratic.NewBadRequestError("reason", "invalid payload")
	}

	opts := defs.NewSyncReposWorkflows(payload.Installation.ID, payload.Action, id)

	_, err := durable.OnHooks().ExecuteWorkflow(ctx.Request().Context(), opts, workflows.SyncRepos, payload)
	if err != nil {
		return erratic.NewInternalServerError("reason", "failed to signal workflow", "error", err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *Webhook) ref(ctx echo.Context, event defs.WebhookEvent, id string) error {
	payload := &defs.WebhookRef{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("failed to bind payload", "error", err.Error())
		return erratic.NewBadRequestError("reason", "invalid payload")
	}

	opts := defs.NewRefWorkflowOptions(payload.Repository.ID, payload.Ref, payload.RefType, "", event.String(), id)

	if payload.RefType == "branch" {
		_, err := durable.OnHooks().ExecuteWorkflow(ctx.Request().Context(), opts, workflows.Ref, payload, event)
		if err != nil {
			return erratic.NewInternalServerError("reason", "failed to signal workflow", "error", err.Error())
		}
	}

	return ctx.NoContent(http.StatusNoContent)
}

// push handles the push event.
func (h *Webhook) push(ctx echo.Context, _ defs.WebhookEvent, id string) error {
	payload := &defs.Push{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("failed to bind payload", "error", err.Error())
		return erratic.NewBadRequestError("reason", "invalid payload")
	}

	if payload.After == defs.NoCommit {
		return ctx.NoContent(http.StatusNoContent)
	}

	action := "created"

	if payload.Deleted {
		action = "deleted"
	}

	if payload.Forced {
		action = "forced"
	}

	opts := defs.NewRefWorkflowOptions(payload.Repository.ID, payload.Ref, "push", payload.After, action, id)

	_, err := durable.
		OnHooks().
		ExecuteWorkflow(ctx.Request().Context(), opts, workflows.Push, payload)
	if err != nil {
		slog.Error("failed to signal workflow", "error", err.Error())
		return erratic.NewInternalServerError("reason", "failed to signal workflow", "error", err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}

// pr handles the pull request event.
func (h *Webhook) pr(ctx echo.Context, event defs.WebhookEvent, id string) error {
	payload := &defs.PR{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("failed to bind payload", "error", err.Error())
		return erratic.NewBadRequestError("reason", "invalid payload")
	}

	opts := defs.NewRefWorkflowOptions(
		payload.GetRepositoryID(), payload.GetHeadBranch(), "pr", fmt.Sprintf("%d", payload.GetNumber()), payload.GetAction(), id)

	_, err := durable.
		OnHooks().
		ExecuteWorkflow(ctx.Request().Context(), opts, workflows.Pr, payload)
	if err != nil {
		slog.Error("failed to signal workflow", "error", err.Error())
		return erratic.NewInternalServerError("reason", "failed to signal workflow", "error", err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}
