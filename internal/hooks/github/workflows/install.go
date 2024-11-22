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
	StatusInstall struct {
		webhook bool
		request bool
	}

	InstallWorkflowState struct {
		do      *activities.Install
		status  StatusInstall
		entity  *entities.GithubInstallation
		request *defs.RequestInstall
		webhook *defs.WebhookInstall

		log log.Logger
	}
)

// Install handles the installation, and associate org with installation id based on incoming signals
// (webhook, complete installation request).
// Critically, it synchronizes the GitHub installation with the internal system.
//
// NOTE: This workflow must not be executed directly, rather always use SignalWithStartWorkflow.
func Install(ctx workflow.Context) error {
	state := NewInstallWorkflowState(ctx)
	selector := workflow.NewSelector(ctx)

	rqst := workflow.GetSignalChannel(ctx, defs.SignalRequestInstall.String())
	selector.AddReceive(rqst, state.on_request(ctx))

	wb := workflow.GetSignalChannel(ctx, defs.SignalWebhookInstall.String())
	selector.AddReceive(wb, state.on_webhook(ctx))

	for !state.done() {
		selector.Select(ctx)
	}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	if err := workflow.ExecuteActivity(ctx, state.do.GetOrCreateInstallation, state.entity).Get(ctx, state.entity); err != nil {
		return err
	}

	for _, repo := range state.webhook.Repositories {
		payload := &defs.SyncRepo{InstallationID: state.entity.ID, Repo: repo, OrgID: state.entity.OrgID}

		selector.AddFuture(workflow.ExecuteActivity(ctx, state.do.AddRepoToInstall, payload), func(f workflow.Future) {})
	}

	for range state.webhook.Repositories {
		selector.Select(ctx)
	}

	return nil
}

func (s *InstallWorkflowState) on_request(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		rx.Receive(ctx, s.request)
		s.status.request = true

		s.entity.OrgID = s.request.OrgID
	}
}

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
