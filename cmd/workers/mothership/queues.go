// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package main

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	queue "go.breu.io/quantm/internal/shared/queue"
)

func configure_core() {
	queue.Core().CreateWorker(
		queues.WithWorkerOptionEnableSessionWorker(true),
	)

	queue.Core().RegisterWorkflow(code.RepoCtrl)
	queue.Core().RegisterWorkflow(code.TrunkCtrl)
	queue.Core().RegisterWorkflow(code.BranchCtrl)
	queue.Core().RegisterWorkflow(code.QueueCtrl)

	queue.Core().RegisterActivity(&code.Activities{})

	queue.Core().RegisterActivity(&github.RepoIO{})
	queue.Core().RegisterActivity(&slack.Activities{})
	queue.Core().RegisterActivity(mutex.PrepareMutexActivity)
}

func configure_providers() {
	queue.Providers().CreateWorker()

	github_workflows := &github.Workflows{}

	queue.Providers().RegisterWorkflow(github_workflows.OnInstallationEvent)
	queue.Providers().RegisterWorkflow(github_workflows.OnInstallationRepositoriesEvent)
	queue.Providers().RegisterWorkflow(github_workflows.PostInstall)
	queue.Providers().RegisterWorkflow(github_workflows.OnPushEvent)
	queue.Providers().RegisterWorkflow(github_workflows.OnCreateOrDeleteEvent)
	queue.Providers().RegisterWorkflow(github_workflows.OnPullRequestEvent)
	queue.Providers().RegisterWorkflow(github_workflows.CollectRepoEventMetadata)
	queue.Providers().RegisterWorkflow(github_workflows.OnWorkflowRunEvent)

	queue.Providers().RegisterActivity(&github.Activities{})
	queue.Providers().RegisterActivity(&slack.Activities{})
}

func configure_mutex() {
	queue.Mutex().CreateWorker()
	queue.Mutex().RegisterWorkflow(mutex.MutexWorkflow)
}
