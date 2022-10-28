// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package github

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/entities"
)

var (
	activities *Activities
)

type (
	Workflows struct{}
)

// OnInstall workflow is executed when we initiate the installation of GitHub core.
//
// In an ideal world, the complete installation request would hit the API after the installation event has hit the
// webhook, however, there can be number of things that can go wrong, and we can receive the complete installation
// request before the push event. To handle this, we use temporal's signal API to provide two possible entry points
// for the system. See the README.md for a detailed explanation on how this workflow works.
//
// NOTE: This workflow is only meant to be started with `SignalWithStartWorkflow`.
func (w *Workflows) OnInstall(ctx workflow.Context) error {
	// prelude
	log := workflow.GetLogger(ctx)
	log.Info("received installation event ...")

	selector := workflow.NewSelector(ctx)
	webhook := &InstallationEventPayload{}
	request := &CompleteInstallationSignalPayload{}
	webhookDone := false
	requestDone := false

	// setting up channels to receive signals
	webhookChannel := workflow.GetSignalChannel(ctx, InstallationEventSignal.String())
	requestChannel := workflow.GetSignalChannel(ctx, CompleteInstallationSignal.String())

	// webhook entry point
	selector.AddReceive(
		webhookChannel,
		func(channel workflow.ReceiveChannel, more bool) {
			log.Info("received installation webhook ...")
			channel.Receive(ctx, webhook)
			webhookDone = true

			switch webhook.Action {
			case "deleted", "suspend", "unsuspend":
				log.Info("installation removed, skipping complete installation request ...")
				requestDone = true
			default:
				log.Info("installation created, waiting for complete installation request ...")
			}
		},
	)

	// complete installation entry point
	selector.AddReceive(
		requestChannel,
		func(channel workflow.ReceiveChannel, more bool) {
			log.Info("received complete installation request ...")
			channel.Receive(ctx, request)
			requestDone = true
		},
	)

	// keep listening for signals until we have received both the installation id and the team id
	for !(webhookDone && requestDone) {
		log.Info("waiting for signals ....")
		selector.Select(ctx)
	}

	log.Info("all signals received, processing ...")

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

	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)
	err := workflow.
		ExecuteActivity(ctx, activities.CreateOrUpdateInstallation, installation).
		Get(ctx, installation)

	if err != nil {
		log.Error("error saving installation", "error", err)
		return err
	}

	// If webhook.Action == "created", save the repository information to the database.
	if webhook.Action == "created" {
		log.Info("saving associated repositories ...")

		// asynchronously save the repos
		for _, repository := range webhook.Repositories {
			log.Info("saving repository ...", "repository", repository.ID)

			repo := &entities.GithubRepo{
				GithubID: repository.ID,
				Name:     repository.Name,
				FullName: repository.FullName,
				TeamID:   installation.TeamID,
			}

			future := workflow.ExecuteActivity(ctx, activities.CreateOrUpdateRepo, repo)
			selector.AddFuture(future, func(f workflow.Future) { log.Info("repository saved ...", repo, repo.GithubID) })
		}

		// wait for all repositories to be saved.
		for range webhook.Repositories {
			selector.Select(ctx)
		}
	}

	log.Info("installation complete", "installation", installation)

	return nil
}

func (w *Workflows) OnPush(ctx workflow.Context, payload PushEventPayload) error {
	log := workflow.GetLogger(ctx)
	log.Debug("received push event ...")

	return nil
}

func (w *Workflows) OnPullRequest(ctx workflow.Context, payload PullRequestEventPayload) error {
	return nil
}
