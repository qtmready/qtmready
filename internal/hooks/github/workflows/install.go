package githubwfs

import (
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db/entities"
	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	StatusInstall struct {
		wehook  bool
		request bool
	}

	InstallWorkflowState struct {
		do      *githubacts.Install
		status  StatusInstall
		entity  *entities.GithubInstallation
		request *githubdefs.RequestInstall
		webhook *githubdefs.WebhookInstall

		log log.Logger
	}
)

// Install installs the Github Integration.
func Install(ctx workflow.Context) error {
	state := NewInstallWorkflowState(ctx)
	selector := workflow.NewSelector(ctx)

	rqst := workflow.GetSignalChannel(ctx, githubdefs.SignalRequestInstall.String())
	selector.AddReceive(rqst, state.on_request(ctx))

	wb := workflow.GetSignalChannel(ctx, githubdefs.SignalWebhookInstall.String())
	selector.AddReceive(wb, state.on_webhook)

	for !state.done() {
		selector.Select(ctx)
	}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	// Get or create the installation.
	err := workflow.ExecuteActivity(ctx, state.do.GetOrCreateInstallation, state.entity).Get(ctx, state.entity)
	if err != nil {
		return err
	}

	state.sync(ctx)

	return nil
}

func (s *InstallWorkflowState) on_request(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		rx.Receive(ctx, s.request)
		s.status.request = true
	}
}

func (s *InstallWorkflowState) on_webhook(rx workflow.ReceiveChannel, more bool) {}

func (s *InstallWorkflowState) done() bool {
	return s.status.request && s.status.wehook
}

func (s *InstallWorkflowState) sync(ctx workflow.Context) {}

func NewInstallWorkflowState(ctx workflow.Context) *InstallWorkflowState {
	return &InstallWorkflowState{
		log:     workflow.GetLogger(ctx),
		request: &githubdefs.RequestInstall{},
		webhook: &githubdefs.WebhookInstall{},
	}
}
