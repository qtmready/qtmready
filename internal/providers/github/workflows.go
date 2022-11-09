// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

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
	selector := workflow.NewSelector(ctx)
	webhook := &InstallationEventPayload{}
	request := &CompleteInstallationSignalPayload{}
	webhookDone := false
	requestDone := false

	// setting up channels to receive signals
	webhookChannel := workflow.GetSignalChannel(ctx, WebhookInstallationEventSignal.String())
	requestChannel := workflow.GetSignalChannel(ctx, RequestCompleteInstallationSignal.String())

	// webhook entry point
	selector.AddReceive(
		webhookChannel,
		func(channel workflow.ReceiveChannel, more bool) {
			log.Info("received webhook installation event ...")
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
			log.Info("saving repository ...")
			log.Debug("repository", "repository", repository)

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

	log.Info("installation complete")
	log.Debug("installation", "installation", installation)

	return nil
}

// OnPush checks if the push event is associated with an open pull request.If so, it will get the idempotent key for
// the immutable rollout. Depending upon the target branch, it will either queue the rollout or update the existing
// rollout.
func (w *Workflows) OnPush(ctx workflow.Context, payload PushEventPayload) error {
	log := workflow.GetLogger(ctx)
	log.Debug("received push event ...")

	return nil
}

// OnPullRequest is the core workflow responsible to create an idempotency key for the immutable infrastructre.
// After creating the idempotency key, it will create a new worklfow to create a rollout.
//
// The spawned workflow will contain the mutex lock to ensure that only one rollout is created at a time.
func (w *Workflows) OnPullRequest(ctx workflow.Context, payload PullRequestEventPayload) error {
	log := workflow.GetLogger(ctx)
	log.Debug("pull request opened, creating new changeset ...")

	return nil
}
