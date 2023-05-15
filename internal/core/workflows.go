// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
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

package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/shared"
)

const (
	unLockTimeOutStackMutex             time.Duration = time.Minute * 30 // TODO: adjust this
	OnPullRequestWorkflowPRSignalsLimit               = 1000             // TODO: adjust this
)

type (
	Workflows        struct{}
	GetAssetsPayload struct {
		StackID     string
		RepoID      gocql.UUID
		ChangeSetID gocql.UUID
	}
)

// ChangesetController controls the rollout lifecycle for one changeset.
func (w *Workflows) ChangesetController(id string) error {
	return nil
}

// DeProvisionInfra de-provisions the infrastructure created for stack deployment.
func (w *Workflows) DeProvisionInfra(ctx workflow.Context, stackID string, resourceData *ResourceData) error {
	return nil
}

// OnPullRequestWorkflow runs indefinitely and controls and synchronizes all actions on stack
// This workflow will start when createStack call is received. it will be the master workflow for all child stack workflows
// like for tasks like creating infrastructure, doing deployment, apperture controller etc.
//
// The workflow waits for the signals from github workflows for pull requests. It consumes events for PR created, updated, merged etc.
func (w *Workflows) StackController(ctx workflow.Context, stackID string) error {
	// deployment map is designed to be used in OnPullRequestWorkflow only
	deployments := make(DeploymentsData)
	logger := workflow.GetLogger(ctx)
	resourceID := "stack." + stackID // stack.<stack id>

	// create and initialize mutex, initializing mutex will start a mutex workflow
	logger.Info("Creating mutex workflow")

	mutex := NewMutex(resourceID, unLockTimeOutStackMutex)
	err := mutex.Init(ctx)
	if err != nil {
		logger.Error("Error in creating mutex for stack", "stack ID", stackID, "Error", err)
	}

	// var prSignalsCounter int = 0

	triggerChannel := workflow.GetSignalChannel(ctx, shared.WorkflowSignalTriggerDeployment.String())
	assetsChannel := workflow.GetSignalChannel(ctx, WorkflowSignalAssetsRetrieved.String())
	infrachannel := workflow.GetSignalChannel(ctx, WorkflowSignalInfraProvisioned.String())
	deploymentchannel := workflow.GetSignalChannel(ctx, WorkflowSignalDeploymentCompleted.String())
	manualOverrideChannel := workflow.GetSignalChannel(ctx, WorkflowSignalManaulOverride.String())

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(triggerChannel, onTriggerSignal(ctx, stackID, deployments))
	selector.AddReceive(assetsChannel, onAssetsRetreivedSignal(ctx, stackID, deployments))
	selector.AddReceive(infrachannel, onInfraProvisionedSignal(ctx, stackID, mutex, deployments))
	selector.AddReceive(deploymentchannel, onDeploymentCompletedSignal(ctx, stackID, deployments))
	selector.AddReceive(manualOverrideChannel, onManualOverrideSignal(ctx, stackID, deployments))

	for {
		// return continue as new if this workflow has processed signals upto a limit
		// if prSignalsCounter >= OnPullRequestWorkflowPRSignalsLimit {
		// 	return workflow.NewContinueAsNewError(ctx, w.OnPullRequestWorkflow, stackID)
		// }
		for {
			logger.Info("waiting for signals ....")
			selector.Select(ctx)
		}
	}
}

func onManualOverrideSignal(ctx workflow.Context, stackID string, deployments DeploymentsData) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	triggerID := int64(0)

	return func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &triggerID)
		logger.Info("manual override for", "Trigger ID", triggerID)
	}
}

// onTriggerSignal is the channel handler for trigger channel
// It will execute GetAssets and update PR deployment state to "GettingAssets".
func onTriggerSignal(ctx workflow.Context, stackID string, deployments DeploymentsData) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		// Receive signal data
		payload := &shared.PullRequestSignal{}
		channel.Receive(ctx, payload)
		logger.Info("received deployment request", "Trigger ID", payload.TriggerID)

		// We want to filter workflows with changeset ID, so create changeset ID here and use it for creating workflow ID
		changesetID, _ := gocql.RandomUUID()

		// Set childworkflow options
		opts := shared.Temporal.Queues[shared.CoreQueue].
			GetChildWorkflowOptions("get_assets", "stack", stackID, "changeset", changesetID.String(),
				"trigger", strconv.FormatInt(payload.TriggerID, 10))

		getAssetsPayload := &GetAssetsPayload{
			StackID:     stackID,
			RepoID:      payload.RepoID,
			ChangeSetID: changesetID,
		}

		// execute GetAssets and wait until spawned
		var execution workflow.Execution
		cctx := workflow.WithChildOptions(ctx, opts)
		err := workflow.
			ExecuteChildWorkflow(cctx, w.GetAssets, getAssetsPayload).
			GetChildWorkflowExecution().
			Get(cctx, &execution)

		if err != nil {
			logger.Error("TODO: Error in executing GetAssets", "Error", err)
		}

		// create and save deployment data against a changeset
		deploymentData := &DeploymentData{}
		deployments[changesetID] = deploymentData
		deploymentData.State = GettingAssets
		deploymentData.WorkflowIDs.GetAssets = execution.ID
	}
}

// GetAssets gets assests for stack including resources, workloads and blueprint.
func (w *Workflows) GetAssets(ctx workflow.Context, payload *GetAssetsPayload) error {
	var (
		future workflow.Future
		err    error = nil
	)

	logger := workflow.GetLogger(ctx)
	assets := new(Assets)
	workloads := SlicedResult[Workload]{}
	resources := SlicedResult[Resource]{}
	repos := SlicedResult[Repo]{}
	blueprint := new(Blueprint)

	selector := workflow.NewSelector(ctx)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)
	providerActivityOpts := workflow.
		ActivityOptions{StartToCloseTimeout: 60 * time.Second, TaskQueue: shared.Temporal.Queues[shared.ProvidersQueue].GetName()}
	pctx := workflow.WithActivityOptions(ctx, providerActivityOpts)

	// get resources for stack
	future = workflow.ExecuteActivity(actx, activities.GetResources, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &resources); err != nil {
			logger.Error("GetResources activity failed", "error", err)
			return
		}
	})

	// get workloads for stack
	future = workflow.ExecuteActivity(actx, activities.GetWorkloads, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &workloads); err != nil {
			logger.Error("GetWorkloads activity failed", "error", err)
			return
		}
	})

	// get repos for stack
	future = workflow.ExecuteActivity(actx, activities.GetRepos, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &repos); err != nil {
			logger.Error("GetRepos activity failed", "error", err)
			return
		}
	})

	// get blueprint for stack
	future = workflow.ExecuteActivity(actx, activities.GetBluePrint, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &blueprint); err != nil {
			logger.Error("GetBluePrint activity failed", "error", err)
			return
		}
	})

	// TODO: come up with a better logic for this
	for i := 0; i < 4; i++ {
		selector.Select(ctx)
		// return if activity failed. TODO: handle race conditions as the 'err' variable is shared among all activities
		if err != nil {
			logger.Error("Exiting due to activity failure")
			return err
		}
	}

	// get commits against the repos
	repoMarker := make([]ChangeSetRepoMarker, len(repos.Data))

	for i := 0; i < len(repos.Data); i++ {
		rep := &repos.Data[i]
		marker := &repoMarker[i]
		p := Core.Providers[rep.Provider] // get the specific provider
		var commitID string
		if err := workflow.ExecuteActivity(pctx, p.GetLatestCommitforRepo, rep.ProviderID, rep.DefaultBranch).
			Get(ctx, &commitID); err != nil {
			logger.Error("Error in getting latest commit ID", "repo", rep.Name, "provider", rep.Provider)
			return fmt.Errorf("Error in getting latest commit ID repo:%s, provider:%s", rep.Name, rep.Provider.String())
		}

		marker.CommitID = commitID
		marker.HasChanged = rep.ID == payload.RepoID
		marker.Provider = rep.Provider.String()
		marker.RepoID = rep.ID.String()
		logger.Debug("Repo", "Name", rep.Name, "Repo marker", marker)
	}

	// save changeset
	stackID, _ := gocql.ParseUUID(payload.StackID)
	changeset := &ChangeSet{
		RepoMarkers: repoMarker,
		ID:          payload.ChangeSetID,
		StackID:     stackID,
	}

	err = workflow.ExecuteActivity(actx, activities.CreateChangeset, changeset, payload.ChangeSetID).Get(ctx, nil)
	if err != nil {
		logger.Error("Error in creating changeset")
	}

	// create assets
	assets.Create()
	assets.ChangesetID = payload.ChangeSetID
	assets.Blueprint = *blueprint
	assets.Repos = append(assets.Repos, repos.Data...)
	assets.Resources = append(assets.Resources, resources.Data...)
	assets.Workloads = append(assets.Workloads, workloads.Data...)
	logger.Debug("Assets retreived", "Assets", assets)

	// signal parent workflow
	PRWorkflow := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.
		SignalExternalWorkflow(ctx, PRWorkflow, "", WorkflowSignalAssetsRetrieved.String(), assets).
		Get(ctx, nil)

	return nil
}

// onAssetsRetreivedSignal will receive assets sent by GetAssets, update deployment state and execute ProvisionInfra.
func onAssetsRetreivedSignal(ctx workflow.Context, stackID string, deployments DeploymentsData) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := &Assets{}
		channel.Receive(ctx, assets)
		logger.Info("received Assets", "changeset", assets.ChangesetID)

		// update deployment state
		deploymentData := deployments[assets.ChangesetID]
		deploymentData.State = GotAssets

		// execute provision infra workflow
		logger.Info("Executing provision Infra workflow")

		var execution workflow.Execution
		opts := shared.Temporal.Queues[shared.CoreQueue].
			GetChildWorkflowOptions("provisionInfra", "stack", stackID, "changeset", assets.ChangesetID.String())

		cctx := workflow.WithChildOptions(ctx, opts)
		err := workflow.
			ExecuteChildWorkflow(cctx, w.ProvisionInfra, assets).
			GetChildWorkflowExecution().Get(cctx, &execution)

		if err != nil {
			logger.Error("TODO: Error in executing ProvisionInfra", "Error", err)
		}

		logger.Info("Executed provision Infra workflow")

		deploymentData.State = ProvisioningInfra
		deploymentData.WorkflowIDs.ProvisionInfra = execution.ID
	}
}

// ProvisionInfra provisions the infrastructure required for stack deployment.
func (w *Workflows) ProvisionInfra(ctx workflow.Context, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	for _, resource := range assets.Resources {
		logger.Info("Creating resource", "Name", resource.Name)
	}

	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalInfraProvisioned.String(), assets).Get(ctx, nil)

	return nil
}

// onInfraProvisionedSignal will receive assets by ProvisionInfra, update deployment state and execute Deploy.
func onInfraProvisionedSignal(ctx workflow.Context, stackID string, mutex *Mutex, deployments DeploymentsData) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := &Assets{}
		channel.Receive(ctx, assets)
		logger.Info("Infra provisioned", "changeset", assets.ChangesetID)

		deploymentData := deployments[assets.ChangesetID]
		deploymentData.State = InfraProvisioned

		var execution workflow.Execution

		opts := shared.Temporal.Queues[shared.CoreQueue].
			GetChildWorkflowOptions("Deployment", "stack", stackID, "changeset", assets.ChangesetID.String())
		cctx := workflow.WithChildOptions(ctx, opts)

		err := workflow.
			ExecuteChildWorkflow(cctx, w.Deploy, stackID, mutex, assets).
			GetChildWorkflowExecution().Get(cctx, &execution)
		if err != nil {
			logger.Error("Error in Executing deployment workflow", "Error", err)
		}

		deploymentData.State = CreatingDeployment
		deploymentData.WorkflowIDs.ProvisionInfra = execution.ID
	}
}

// Deploy deploys the stack.
func (w *Workflows) Deploy(ctx workflow.Context, stackID string, mutex *Mutex, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	// Acquire lock
	logger.Info("Deployment initiated", "changeset", assets.ChangesetID)

	unlockFunc, err := mutex.Lock(ctx)
	if err != nil {
		logger.Error("Error in acquiring lock", "Error", err)
		return err
	}

	// simulate critical section
	_ = workflow.Sleep(ctx, 60*time.Second)

	// release lock
	_ = unlockFunc()

	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalDeploymentCompleted.String(), assets)

	return nil
}

// onDeploymentCompletedSignal will conclude the deployment.
func onDeploymentCompletedSignal(ctx workflow.Context, stackID string, deployments DeploymentsData) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := &Assets{}
		channel.Receive(ctx, assets)
		logger.Info("Deployment complete", "changeset", assets.ChangesetID)
		delete(deployments, assets.ChangesetID)

		logger.Info("Deleted deployment data", "changeset", assets.ChangesetID)
	}
}
