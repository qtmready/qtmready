package github

import (
	"go.temporal.io/sdk/workflow"
)

type Workflows struct{}

var activities *Activities

// OnInstall workflow is executed when we initiate the installation of github app.
// It has two possible entry points
//
// 1. When we receive a webhook event from github. The payload has all the data, except the team_id.
// 2. When we receive the complete installation event from our ui. Github redirects to the ui, which knows the team_id.
//
// In order to cater to latency and hops, we might recieve the complete installation event first. Therefore, we are
// handling for whichever comes first. In order to do that, this workflow is only meant to be started via signals. See
// README.md for more details.
func (w *Workflows) OnInstall(ctx workflow.Context) error {
	// prelude
	webhook := &InstallationEventPayload{}
	request := &CompleteInstallationPayload{}
	webhookDone := false
	requestDone := false
	// result := &entities.GithubInstallation{}
	selector := workflow.NewSelector(ctx)

	// setting up channels to recieve signals
	webhookChannel := workflow.GetSignalChannel(ctx, InstallationEventSignal.String())
	requestChannel := workflow.GetSignalChannel(ctx, CompleteInstallationSignal.String())

	// setting up signal reciever for webhook
	selector.AddReceive(
		webhookChannel,
		func(channel workflow.ReceiveChannel, more bool) {
			channel.Receive(ctx, webhook)
			webhookDone = true
		},
	)

	// setting up signal reciever for http request
	selector.AddReceive(
		requestChannel,
		func(channel workflow.ReceiveChannel, more bool) {
			channel.Receive(ctx, request)
			requestDone = true
		},
	)

	// keep listening for signals until we have received both the installation id and the team id
	for !webhookDone && !requestDone {
		selector.Select(ctx)
	}

	return nil
}

func (w *Workflows) OnPush(ctx workflow.Context, payload PushEventPayload) error {
	return nil
}

func (w *Workflows) OnPullRequest(ctx workflow.Context, payload PullRequestEventPayload) error {
	return nil
}
