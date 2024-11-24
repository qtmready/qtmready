package workflows

import (
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	// StatusInstall represents the status of the installation workflow.
	StatusInstall struct {
		webhook bool // indicates whether the webhook signal has been received
		request bool // indicates whether the request signal has been received
	}

	// InstallWorkflowState represents the state of the GitHub installation workflow.
	InstallWorkflowState struct {
		do      *activities.Install          // Install activities
		status  StatusInstall                // Status of the installation workflow
		entity  *entities.GithubInstallation // GitHub installation entity
		request *defs.RequestInstall         // Request installation data
		webhook *defs.WebhookInstall         // Webhook installation data

		log log.Logger // Workflow logger
	}
)

// Install handles GitHub installation synchronization.  Install uses signals to coordinate installation.
// It synchronizes GitHub installations with the internal system.  Install should not be invoked directly.  Use
// SignalWithStartWorkflow instead.
func Install(ctx workflow.Context) error {
	state := NewInstallWorkflowState(ctx)
	selector := workflow.NewSelector(ctx)

	rqst := workflow.GetSignalChannel(ctx, defs.SignalRequestInstall.String())
	selector.AddReceive(rqst, state.on_request(ctx))

	webhook := workflow.GetSignalChannel(ctx, defs.SignalWebhookInstall.String())
	selector.AddReceive(webhook, state.on_webhook(ctx))

	for !state.done() {
		selector.Select(ctx)
	}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	if err := workflow.ExecuteActivity(ctx, state.do.GetOrCreateInstallation, state.entity).Get(ctx, state.entity); err != nil {
		return err
	}

	for _, repo := range state.webhook.Repositories {
		payload := &defs.SyncRepoPayload{InstallationID: state.entity.ID, Repo: repo, OrgID: state.entity.OrgID}
		selector.AddFuture(workflow.ExecuteActivity(ctx, state.do.AddRepoToInstall, payload), func(f workflow.Future) {})
	}

	for range state.webhook.Repositories {
		selector.Select(ctx)
	}

	return nil
}

// on_request is a channel handler for the request signal. Request handler is used to set the OrgID.
func (s *InstallWorkflowState) on_request(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		rx.Receive(ctx, s.request)
		s.status.request = true

		s.entity.OrgID = s.request.OrgID
	}
}

// on_webhook is a channel handler for the webhook signal. The webhook contains the installation information,
// e.g. installation id, account, and the repos that are part of the installation.
func (s *InstallWorkflowState) on_webhook(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		rx.Receive(ctx, s.webhook)
		s.status.webhook = true

		s.entity.InstallationID = s.webhook.Installation.ID
		s.entity.InstallationLogin = s.webhook.Installation.Account.Login
		s.entity.InstallationLoginID = s.webhook.Installation.Account.ID
		s.entity.InstallationType = s.webhook.Installation.Account.Type
		s.entity.SenderID = s.webhook.Sender.ID
		s.entity.SenderLogin = s.webhook.Sender.Login

		if s.webhook.Action != "created" && s.webhook.Action != "updated" {
			s.status.request = true
		}
	}
}

// done returns true if the installation is complete.
func (s *InstallWorkflowState) done() bool {
	return s.status.request && s.status.webhook
}

func NewInstallWorkflowState(ctx workflow.Context) *InstallWorkflowState {
	return &InstallWorkflowState{
		log:     workflow.GetLogger(ctx),
		status:  StatusInstall{webhook: false, request: false},
		entity:  &entities.GithubInstallation{},
		request: &defs.RequestInstall{},
		webhook: &defs.WebhookInstall{},
	}
}
