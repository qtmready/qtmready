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
	"time"

	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/shared"
)

const (
	unLockTimeOutStackMutex             time.Duration = time.Minute * 30 //TODO: adjust this
	OnPullRequestWorkflowPRSignalsLimit               = 1000             // TODO: adjust this
)

type (
	Workflows    struct{}
	ResourceData struct{}

	SlicedResult[T any] struct {
		Data []T `json:"data"`
	}

	// Assets contains all the assets fetched from DB against a stack.
	Assets struct {
		Repos           []Repo
		Resources       []Resource
		Workloads       []Workload
		Blueprint       Blueprint
		ResourcesConfig []ResourceData
		PRID            int64
		ChangesetID     gocql.UUID
	}

	ChildWorkflowIDs struct {
		GetAssets      string
		ProvisionInfra string
		Deployment     string
	}

	State int64

	DeploymentData struct {
		State       State
		WorkflowIDs ChildWorkflowIDs
	}

	DeploymentDataMap map[int64]*DeploymentData
	AssetsMap         map[int64]*Assets
)

const (
	GettingAssets State = iota
	GotAssets
	ProvisioningInfra
	InfraProvisioned
	CreatingDeployment
)

var (
	activities *Activities
)

// ChangesetController controls the rollout lifecycle for one changeset.
func (w *Workflows) ChangesetController(id string) error {
	return nil
}

// DeProvisionInfraWorkflow de-provisions the infrastructure created for stack deployment.
func (w *Workflows) DeProvisionInfraWorkflow(ctx workflow.Context, stackID string, resourceData *ResourceData) error {
	return nil
}

// OnPullRequestWorkflow runs indefinitely and controls and synchronizes all actions on stack
// This workflow will start when createStack call is received. it will be the master workflow for all child stack workflows
// like for tasks like creating infrastructure, doing deployment, apperture controller etc.
//
// The workflow waits for the signals from github workflows for pull requests. It consumes events for PR created, updated, merged etc.
func (w *Workflows) OnPullRequestWorkflow(ctx workflow.Context, stackID string) error {
	// deployment map is designed to be used in OnPullRequestWorkflow only
	deploymentDataMap := make(DeploymentDataMap)
	logger := workflow.GetLogger(ctx)
	currentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	resourceID := "stack." + stackID // stack.<stack id>

	// create and initialize mutex, initializing mutex will start a mutex workflow
	logger.Info("Creating mutex workflow")

	mutex := NewMutex(currentWorkflowID, resourceID, unLockTimeOutStackMutex)
	_ = mutex.Init(ctx)

	// var prSignalsCounter int = 0

	prChannel := workflow.GetSignalChannel(ctx, shared.WorkflowSignalPullRequest.String())
	assetsChannel := workflow.GetSignalChannel(ctx, WorkflowSignalAssetsRetrieved.String())
	infrachannel := workflow.GetSignalChannel(ctx, WorkflowSignalInfraProvisioned.String())
	deploymentchannel := workflow.GetSignalChannel(ctx, WorkflowSignalDeploymentCompleted.String())
	manualOverrideChannel := workflow.GetSignalChannel(ctx, WorkflowSignalManaulOverride.String())

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(prChannel, onPRSignal(ctx, stackID, deploymentDataMap))
	selector.AddReceive(assetsChannel, onAssetsRetreivedSignal(ctx, stackID, deploymentDataMap))
	selector.AddReceive(infrachannel, onInfraProvisionedSignal(ctx, stackID, deploymentDataMap))
	selector.AddReceive(deploymentchannel, onDeploymentCompletedSignal(ctx, stackID, deploymentDataMap))
	selector.AddReceive(manualOverrideChannel, onManualOverrideSignal(ctx, stackID, deploymentDataMap))

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

func onManualOverrideSignal(ctx workflow.Context, stackID string, deploymentMap DeploymentDataMap) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	prID := int64(0)

	return func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &prID)
		logger.Info("manual override for", "PR ID", prID)
	}
}

// onPRSignal is the channel handler for PR channel
// It will execute getAssetsWorkflow and update PR deployment state to "GettingAssets".
func onPRSignal(ctx workflow.Context, stackID string, deploymentMap DeploymentDataMap) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		// Receive signal data
		payload := &shared.PullRequestSignal{}
		channel.Receive(ctx, payload)
		logger.Info("received PR event", "PR ID", payload.PullRequestID)

		// execute child workflow and wait for it to spawn
		var execution workflow.Execution

		opts := shared.Temporal.Queues[shared.CoreQueue].GetChildWorkflowOptions("get_assets", "stack", stackID)
		ctx = workflow.WithChildOptions(ctx, opts)
		err := workflow.
			ExecuteChildWorkflow(ctx, w.GetAssetsWorkflow, stackID).
			GetChildWorkflowExecution().
			Get(ctx, &execution)

		if err != nil {
			logger.Info("TODO: Handle error", "error", err)
		}

		// create and save deployment data
		deploymentData := &DeploymentData{}
		deploymentMap[payload.PullRequestID] = deploymentData
		deploymentData.State = GettingAssets
		deploymentData.WorkflowIDs.GetAssets = execution.ID
	}
}

// GetAssetsWorkflow gets assests for stack including resources, workloads and blueprint.
func (w *Workflows) GetAssetsWorkflow(ctx workflow.Context, stackID string, prID int64) error {
	var future workflow.Future

	logger := workflow.GetLogger(ctx)
	assets := new(Assets)
	resources := SlicedResult[Resource]{}
	workloads := SlicedResult[Workload]{}
	repos := SlicedResult[Repo]{}
	blueprint := new(Blueprint)
	var err error = nil

	selector := workflow.NewSelector(ctx)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	act := workflow.WithActivityOptions(ctx, activityOpts)

	// get resources for stack
	future = workflow.ExecuteActivity(act, activities.GetResources, stackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = future.Get(ctx, &resources); err != nil {
			logger.Error("GetResources activity failed", "error", err)
			return
		}
		assets.resources = resources.Data
		logger.Debug("GetResources activity complete", "resources", resources)
	})

	// get workloads for stack
	future = workflow.ExecuteActivity(act, activities.GetWorkloads, stackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = future.Get(ctx, &workloads); err != nil {
			logger.Error("GetWorkloads activity failed", "error", err)
			return
		}

		assets.workloads = workloads.Data
		logger.Debug("GetWorkloads activity complete", "workloads", assets.workloads)
	})

	// get repos for stack
	future = workflow.ExecuteActivity(act, activities.GetRepos, stackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = future.Get(ctx, &repos); err != nil {
			logger.Error("GetRepos activity failed", "error", err)
			return
		}
		assets.repos = repos.Data
		logger.Debug("GetRepos activity complete", "repos", assets.repos)
	})

	// get blueprint for stack
	future = workflow.ExecuteActivity(act, activities.GetBluePrint, stackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = future.Get(ctx, &blueprint); err != nil {
			logger.Error("GetBluePrint activity failed", "error", err)
			return
		}
		assets.blueprint = *blueprint
		logger.Debug("GetBluePrint activity complete", "blueprint", blueprint)
	})

	logger.Info("Assets retreived", "Assets", assets)

	// TODO: come up with a better logic for this
	for i := 0; i < 1; i++ {
		selector.Select(ctx)
		// return if activity failed. TODO: handle race conditions as the variable is shared among all activities
		if err != nil {
			return err
		}
	}

	// TODO: create changeset id, for now making changeset id = pr id
	stackUUID, _ := gocql.ParseUUID(stackID)
	changeset := &ChangeSet{
		StackID: stackUUID,
	}

	// get commits against the repos
	err = workflow.ExecuteActivity(act, activities.CreateChangeset, changeset).Get(ctx, nil)
	if err != nil {
		logger.Error("Error in creating changeset")
	}

	// TODO: Create Assets here!
	assets.PRID = prID

	// signal parent workflow
	// PRWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	// workflow.
	// 	SignalExternalWorkflow(ctx, PRWorkflowID, "", WorkflowSignalAssetsRetrieved.String(), assets).
	// 	Get(ctx, nil)

	return nil
}

// onAssetsRetreivedSignal will receive assets sent by GetAssetsWorkflow, update deployment state and execute provisionInfraWorkflow.
func onAssetsRetreivedSignal(ctx workflow.Context, stackID string, deploymentMap DeploymentDataMap) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		logger.Info("received Assets")

		assets := &Assets{}
		channel.Receive(ctx, assets)

		// update deployment state
		deploymentData := deploymentMap[assets.PRID]
		deploymentData.State = GotAssets

		// execute provision infra workflow
		logger.Info("Executing provision Infra workflow")

		var execution workflow.Execution

		opts := shared.Temporal.Queues[shared.CoreQueue].GetChildWorkflowOptions("provisionInfra", "stack", stackID)
		ctx = workflow.WithChildOptions(ctx, opts)
		err := workflow.
			ExecuteChildWorkflow(ctx, w.ProvisionInfraWorkflow, assets).
			GetChildWorkflowExecution().Get(ctx, execution)

		if err != nil {
			logger.Info("TODO: Handle error")
		}

		logger.Info("Executed provision Infra workflow")

		deploymentData.State = ProvisioningInfra
		deploymentData.WorkflowIDs.ProvisionInfra = execution.ID
	}
}

// ProvisionInfraWorkflow provisions the infrastructure required for stack deployment.
func (w *Workflows) ProvisionInfraWorkflow(ctx workflow.Context, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	for _, resource := range assets.Resources {
		logger.Info("Creating resource", resource.Name)
	}

	PRWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	workflow.SignalExternalWorkflow(ctx, PRWorkflowID, "", WorkflowSignalInfraProvisioned.String(), assets)

	return nil
}

func onInfraProvisionedSignal(ctx workflow.Context, stackID string, deploymentMap DeploymentDataMap) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := &Assets{}
		channel.Receive(ctx, assets)
		logger.Info("Infra provisioned")

		deploymentData := deploymentMap[assets.PRID]
		deploymentData.State = InfraProvisioned

		var execution workflow.Execution

		opts := shared.Temporal.Queues[shared.CoreQueue].GetChildWorkflowOptions("Deployment", "stack", stackID)
		ctx = workflow.WithChildOptions(ctx, opts)
		err := workflow.
			ExecuteChildWorkflow(ctx, w.DeploymentWorkflow, assets).
			GetChildWorkflowExecution().Get(ctx, execution)

		if err != nil {
			logger.Info("TODO: Handle error")
		}

		deploymentData.State = CreatingDeployment
		deploymentData.WorkflowIDs.ProvisionInfra = execution.ID
	}
}

// DeploymentWorkflow deploys the stack.
func (w *Workflows) DeploymentWorkflow(ctx workflow.Context, stackID string, mutex *Mutex, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	// Acquire lock
	logger.Info("Deployment initiated", "Stack ID", stackID)

	unlockFunc, err := mutex.Lock(ctx)
	if err != nil {
		logger.Info("Error in acquiring lock", err)
	}

	// simulate critical section
	_ = workflow.Sleep(ctx, 5*time.Second)

	// release lock
	_ = unlockFunc()

	PRWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	workflow.SignalExternalWorkflow(ctx, PRWorkflowID, "", WorkflowSignalDeploymentCompleted.String(), assets)

	return nil
}

func onDeploymentCompletedSignal(ctx workflow.Context, stackID string, deploymentMap DeploymentDataMap) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := &Assets{}
		channel.Receive(ctx, assets)
		logger.Info("Deployment complete for", "PR ID", assets.PRID)
		delete(deploymentMap, assets.PRID)

		logger.Info("Deleted deployment data for", "PR ID", assets.PRID)
	}
}
