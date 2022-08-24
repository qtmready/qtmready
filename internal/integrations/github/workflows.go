package github

import (
	"time"

	"go.breu.io/ctrlplane/internal/entities"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

type Workflows struct{}

var activities *Activities

// OnInstall workflow is executed when we initiate the installation of github app.
//
// In an ideal world, the complete installation request would hit the API after the installation event has hit the
// webhook, however, there can be number of things that can go wrong and we can recieve the complete installation
// request before the push event. To handle this, we use temporal.io signal API to provide two possible entry points
// for the system. See the README.md for a detailed explaination on how this workflow works.
//
// NOTE: This workflow is only meant to be started with `SignalWithStartWorkflow`
func (w *Workflows) OnInstall(ctx workflow.Context) error {
	// prelude
	log := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)
	webhook := &InstallationEventPayload{}
	request := &CompleteInstallationPayload{}
	webhookDone := false
	requestDone := false

	// setting up channels to recieve signals
	webhookChannel := workflow.GetSignalChannel(ctx, InstallationEventSignal.String())
	requestChannel := workflow.GetSignalChannel(ctx, CompleteInstallationSignal.String())

	// push event entry point
	selector.AddReceive(
		webhookChannel,
		func(channel workflow.ReceiveChannel, more bool) {
			log.Info("webhook: ", zap.Any("payload", webhook))
			channel.Receive(ctx, webhook)
			webhookDone = true

			switch webhook.Action {
			case "deleted", "suspend", "unsuspend":
				log.Info("delete/suspend/unsuspend event.")
				requestDone = true
			default:
				log.Info("create event.")
			}
		},
	)

	// complete installation entry point
	selector.AddReceive(
		requestChannel,
		func(channel workflow.ReceiveChannel, more bool) {
			log.Info("request: ", zap.Any("payload", request))
			channel.Receive(ctx, request)
			requestDone = true
		},
	)

	// keep listening for signals until we have received both the installation id and the team id
	for !webhookDone && !requestDone {
		log.Info("selecting...")
		selector.Select(ctx)
	}

	// Finalizing the installation
	installation := &entities.GithubInstallation{
		TeamID:            request.TeamID,
		InstallationID:    webhook.Installation.ID,
		InstallationLogin: webhook.Installation.Account.Login,
		InstallationType:  webhook.Installation.Account.Type,
		SenderID:          webhook.Sender.ID,
		SenderLogin:       webhook.Sender.Login,
		Status:            webhook.Action,
	}

	opt := workflow.ActivityOptions{StartToCloseTimeout: 30 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opt)
	err := workflow.
		ExecuteActivity(ctx, activities.GetOrCreateInstallation, installation).
		Get(ctx, installation)

	if err != nil {
		log.Error("error saving installation", zap.Error(err))
		return err
	}

	// TODO: save the repository related data.

	return nil
}

func (w *Workflows) OnPush(ctx workflow.Context, payload PushEventPayload) error {
	return nil
}

func (w *Workflows) OnPullRequest(ctx workflow.Context, payload PullRequestEventPayload) error {
	return nil
}
