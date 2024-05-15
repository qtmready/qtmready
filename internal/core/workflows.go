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
	"strings"
	"time"

	"github.com/gocql/gocql"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/shared"
)

const (
	unLockTimeOutStackMutex             time.Duration = time.Minute * 30 // TODO: adjust this
	OnPullRequestWorkflowPRSignalsLimit               = 1000             // TODO: adjust this
)

type (
	Workflows        struct{}
	GetAssetsPayload struct {
		StackID       string
		RepoID        gocql.UUID
		ChangeSetID   gocql.UUID
		Image         string
		ImageRegistry string
		Digest        string
	}
)

// ChangesetController controls the rollout lifecycle for one changeset.
func (w *Workflows) ChangesetController(id string) error {
	return nil
}

// StackController runs indefinitely and controls and synchronizes all actions on stack.
// This workflow will start when createStack call is received. it will be the master workflow for all child stack workflows
// for tasks like creating infrastructure, doing deployment, apperture controller etc.
//
// The workflow waits for the signals from the git provider. It consumes events for PR created, updated, merged etc.
func (w *Workflows) StackController(ctx workflow.Context, stackID string) error {
	logger := workflow.GetLogger(ctx)
	// wait for merge complete signal
	ch := workflow.GetSignalChannel(ctx, shared.WorkflowSignalCreateChangeset.String())
	payload := &shared.CreateChangesetSignal{}
	ch.Receive(ctx, payload)

	logger.Info("Stack controller", "signal payload", payload)

	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	// get repos for stack
	repos := SlicedResult[Repo]{}
	if err := workflow.ExecuteActivity(actx, activities.GetRepos, stackID).Get(ctx, &repos); err != nil {
		logger.Error("Get repos activity", "error", err)
		return err
	}

	logger.Info("Stack controller: going to create repomarkers")

	providerActivityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
	}
	pctx := workflow.WithActivityOptions(ctx, providerActivityOpts)

	// get commits against the repos
	repoMarkers := make([]ChangeSetRepoMarker, len(repos.Data))
	for idx, repo := range repos.Data {
		marker := &repoMarkers[idx]
		p := Instance().RepoProvider(repo.Provider) // get the specific provider

		commit := LatestCommit{}
		if err := workflow.ExecuteActivity(pctx, p.GetLatestCommit, repo.ProviderID, repo.DefaultBranch).Get(ctx, &commit); err != nil {
			logger.Error("Repo provider activities: Get latest commit activity", "error", err)
			return err
		}

		marker.CommitID = commit.SHA
		marker.Provider = repo.Provider.String()
		marker.RepoID = repo.ID.String()
		logger.Debug("Debug only", "Commit ID updated for repo ", marker.RepoID)

		// update commit id for the recently changed repo
		if marker.RepoID == payload.RepoID {
			marker.CommitID = payload.CommitID
			marker.HasChanged = true // the repo in which commit was made
		}

		logger.Debug("Debug only", "Repo", repo, "Repo marker", marker)
	}

	// create changeset before deploying the updated changeset
	changesetID, _ := gocql.RandomUUID()
	stackUUID, _ := gocql.ParseUUID(stackID)
	changeset := &ChangeSet{
		RepoMarkers: repoMarkers,
		ID:          changesetID,
		StackID:     stackUUID,
	}

	if err := workflow.ExecuteActivity(actx, activities.CreateChangeset, changeset, changeset.ID).Get(ctx, nil); err != nil {
		logger.Error("Create changeset activity", "error", err)
	}

	logger.Info("Stack controller", "changeset created", changeset)

	for idx, repo := range repos.Data {
		provider := repo.ProviderID
		commitID := repoMarkers[idx].CommitID

		p := Instance().RepoProvider(repo.Provider) // get the specific provider

		if err := workflow.ExecuteActivity(pctx, p.TagCommit, provider, commitID, changesetID.String(), "Tagged by quantm"); err != nil {
			logger.Error("Repo provider activities: Tag commit activity", "error", err)
		}

		if err := workflow.ExecuteActivity(pctx, p.DeployChangeset, repo.ProviderID, changeset.ID).Get(ctx, nil); err != nil {
			logger.Error("Repo provider activities: Deploy changeset activity", "error", err)
		}
	}

	logger.Info("deployment done........")

	// // deployment map is designed to be used in OnPullRequestWorkflow only
	// logger := workflow.GetLogger(ctx)
	// lockID := "stack." + stackID // stack.<stack id>
	// deployments := make(Deployments)

	// // the idea is to save active infra which will be serving all the traffic and use this active infra as reference for next deployment
	// // this is not being used that as active infra for cloud run is being fetched from the cloud which is not an efficient approach
	// activeInfra := make(Infra)

	// // create and initialize mutex, initializing mutex will start a mutex workflow
	// logger.Info("creating mutex for stack", "stack", stackID)

	// lock := mutex.New(
	// 	mutex.WithCallerContext(ctx),
	// 	mutex.WithID(lockID),
	// )

	// if err := lock.Start(ctx); err != nil {
	// 	logger.Debug("unable to start mutex workflow", "error", err)
	// }

	// triggerChannel := workflow.GetSignalChannel(ctx, shared.WorkflowSignalDeploymentStarted.String())
	// assetsChannel := workflow.GetSignalChannel(ctx, WorkflowSignalAssetsRetrieved.String())
	// infraChannel := workflow.GetSignalChannel(ctx, WorkflowSignalInfraProvisioned.String())
	// deploymentChannel := workflow.GetSignalChannel(ctx, WorkflowSignalDeploymentCompleted.String())
	// manualOverrideChannel := workflow.GetSignalChannel(ctx, WorkflowSignalManualOverride.String())

	// selector := workflow.NewSelector(ctx)
	// selector.AddReceive(triggerChannel, onDeploymentStartedSignal(ctx, stackID, deployments))
	// selector.AddReceive(assetsChannel, onAssetsRetrievedSignal(ctx, stackID, deployments))
	// selector.AddReceive(infraChannel, onInfraProvisionedSignal(ctx, stackID, lock, deployments, activeInfra))
	// selector.AddReceive(deploymentChannel, onDeploymentCompletedSignal(ctx, stackID, deployments))
	// selector.AddReceive(manualOverrideChannel, onManualOverrideSignal(ctx, stackID, deployments))

	// // var prSignalsCounter int = 0
	// // return continue as new if this workflow has processed signals upto a limit
	// // if prSignalsCounter >= OnPullRequestWorkflowPRSignalsLimit {
	// // 	return workflow.NewContinueAsNewError(ctx, w.OnPullRequestWorkflow, stackID)
	// // }
	// for {
	// 	logger.Info("waiting for signals ....")
	// 	selector.Select(ctx)
	// }

	return nil
}

func CheckEarlyWarning(
	ctx workflow.Context, rpa RepoProviderActivities, mpa MessageProviderActivities, pushEvent *shared.PushEventSignal,
) error {
	logger := workflow.GetLogger(ctx)
	branchName := pushEvent.RefBranch
	installationID := pushEvent.InstallationID
	repoID := strconv.FormatInt(pushEvent.RepoID, 10)
	repoName := pushEvent.RepoName
	repoOwner := pushEvent.RepoOwner
	defaultBranch := pushEvent.DefaultBranch

	providerActOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	pctx := workflow.WithActivityOptions(ctx, providerActOpts)

	// check merge conflicts
	// create a temporary copy of default branch for the target branch (under inspection)
	// if the rebase with the target branch returns error, raise warning
	logger.Info("Check early warning", "push event", pushEvent)

	commit := &LatestCommit{}
	if err := workflow.ExecuteActivity(pctx, rpa.GetLatestCommit, repoID, defaultBranch).Get(ctx, commit); err != nil {
		logger.Error("Repo provider activities: Get latest commit activity", "error", err)
		return err
	}

	// create a temp branch/ref
	temp := defaultBranch + "-tempcopy-for-target-" + branchName

	// delete the branch if it is present already
	if err := workflow.ExecuteActivity(pctx, rpa.DeleteBranch, installationID, repoName, repoOwner, temp).
		Get(ctx, nil); err != nil {
		logger.Error("Repo provider activities: Delete branch activity", "error", err)
		return err
	}

	// create new ref
	if err := workflow.
		ExecuteActivity(pctx, rpa.CreateBranch, installationID, repoID, repoName, repoOwner, commit.SHA, temp).
		Get(ctx, nil); err != nil {
		logger.Error("Repo provider activities: Create branch activity", "error", err)
		return err
	}

	// get the teamID from repo table
	teamID := ""
	if err := workflow.ExecuteActivity(pctx, rpa.GetRepoTeamID, repoID).Get(ctx, &teamID); err != nil {
		logger.Error("Repo provider activities: Get repo teamID activity", "error", err)
		return err
	}

	if err := workflow.ExecuteActivity(pctx, rpa.MergeBranch, installationID, repoName, repoOwner, temp, branchName).
		Get(ctx, nil); err != nil {
		// dont want to retry this workflow so not returning error, just log and return
		logger.Error("Repo provider activities: Merge branch activity", "error", err)

		// send slack notification
		if err = workflow.ExecuteActivity(pctx, mpa.SendMergeConflictsMessage, teamID, commit).Get(ctx, nil); err != nil {
			logger.Error("Message provider activities: Send merge conflicts message activity", "error", err)
			return err
		}

		return nil
	}

	logger.Info("Merge conflicts NOT detected")

	// detect 200+ changes
	// calculate all changes between default branch (e.g. main) with the target branch
	// raise warning if the changes are more than 200 lines
	logger.Info("Going to detect 200+ changes")

	branchChnages := &BranchChanges{}

	if err := workflow.ExecuteActivity(pctx, rpa.ChangesInBranch, installationID, repoName, repoOwner, defaultBranch, branchName).
		Get(ctx, branchChnages); err != nil {
		logger.Error("Repo provider activities: Changes in branch  activity", "error", err)
		return err
	}

	threshold := 200
	if branchChnages.Changes > threshold {
		if err := workflow.
			ExecuteActivity(pctx, mpa.SendNumberOfLinesExceedMessage, teamID, repoName, branchName, threshold, branchChnages).
			Get(ctx, nil); err != nil {
			logger.Error("Message provider activities: Send number of lines exceed message activity", "error", err)
			return err
		}
	}

	logger.Info("200+ changes NOT detected")

	return nil
}

// when a push event is received by quantm, branch controller gets active.
// if the push event occurred on the default branch (e.g. main) quantm,
// rebases all available branches with the default one.
// otherwise it runs early detection algorithm to see if the branch
// could be problematic when a PR is opened on it.
func (w *Workflows) BranchController(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Branch controller", "waiting for signal", shared.WorkflowPushEvent.String())

	// get push event data via workflow signal
	ch := workflow.GetSignalChannel(ctx, shared.WorkflowPushEvent.String())

	payload := &shared.PushEventSignal{}

	// receive signal payload
	ch.Receive(ctx, payload)

	timeout := 100 * time.Second
	id := fmt.Sprintf("repo.%s.branch.%s", payload.RepoName, payload.RefBranch)
	lock := mutex.New(
		mutex.WithResourceID(id),
		mutex.WithTimeout(timeout+(10*time.Second)),
		mutex.WithHandler(ctx),
	)

	if err := lock.Prepare(ctx); err != nil {
		return err
	}

	if err := lock.Acquire(ctx); err != nil {
		return err
	}

	logger.Debug("Branch controller", "signal payload", payload)

	rpa := Instance().RepoProvider(RepoProvider(payload.RepoProvider))
	mpa := Instance().MessageProvider(MessageProviderSlack) // TODO - maybe not hardcode to slack and get from payload

	providerActOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	actx := workflow.WithActivityOptions(ctx, providerActOpts)

	commit := &LatestCommit{}
	if err := workflow.ExecuteActivity(actx, rpa.GetLatestCommit, (strconv.FormatInt(payload.RepoID, 10)), payload.RefBranch).
		Get(ctx, commit); err != nil {
		logger.Error("Repo provider activities: Get latest commit activity", "error", err)
		return err
	}

	// if the push comes at the default branch i.e. main rebase all branches with main
	if payload.RefBranch == payload.DefaultBranch {
		var branchNames []string
		if err := workflow.ExecuteActivity(actx, rpa.GetAllBranches, payload.InstallationID, payload.RepoName, payload.RepoOwner).
			Get(ctx, &branchNames); err != nil {
			logger.Error("Repo provider activities: Get all branches activity", "error", err)
			return err
		}

		logger.Debug("Branch controller", "Total branches", len(branchNames))

		for _, branch := range branchNames {
			if strings.Contains(branch, "-tempcopy-for-target-") || branch == payload.DefaultBranch {
				// no need to do rebase with quantm created temp branches
				continue
			}

			logger.Debug("Branch controller", "Testing conflicts with branch", branch)

			if err := workflow.ExecuteActivity(
				actx, rpa.MergeBranch, payload.InstallationID, payload.RepoName, payload.RepoOwner, payload.DefaultBranch, branch,
			).
				Get(ctx, nil); err != nil {
				logger.Error("Repo provider activities: Merge branch activity", "error", err)

				// get the teamID from repo table
				teamID := ""
				if err := workflow.ExecuteActivity(actx, rpa.GetRepoTeamID, strconv.FormatInt(payload.RepoID, 10)).Get(ctx, &teamID); err != nil {
					logger.Error("Repo provider activities: Get repo TeamID activity", "error", err)
					return err
				}

				if err = workflow.ExecuteActivity(actx, mpa.SendMergeConflictsMessage, teamID, commit).Get(ctx, nil); err != nil {
					logger.Error("Message provider activities: Send merge conflicts message activity", "error", err)
					return err
				}
			}
		}

		_ = lock.Release(ctx)
		_ = lock.Cleanup(ctx)

		return nil
	}

	// check if the target branch would have merge conflicts with the default branch or it has too much changes
	if err := CheckEarlyWarning(ctx, rpa, mpa, payload); err != nil {
		return err
	}

	// execute child workflow for stale detection
	// if a branch is stale for a long time (5 days in this case) raise warning
	logger.Debug("going to detect stale branch")

	wf := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("repo"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.RepoID, 10)),
			shared.WithWorkflowElement("branch"),
			shared.WithWorkflowElementID(payload.RefBranch),
			shared.WithWorkflowProp("type", "stale_detection"),
		)
	opts.ParentClosePolicy = enums.PARENT_CLOSE_POLICY_ABANDON

	var execution workflow.Execution

	cctx := workflow.WithChildOptions(ctx, opts)
	err := workflow.ExecuteChildWorkflow(
		cctx,
		wf.StaleBranchDetection,
		payload,
		payload.RefBranch,
		commit.SHA,
	).
		GetChildWorkflowExecution().
		Get(cctx, &execution)

	if err != nil {
		// dont want to retry this workflow so not returning error, just log and return
		logger.Error("BranchController", "error executing child workflow", err)
		return nil
	}

	return nil
}

func (w *Workflows) StaleBranchDetection(
	ctx workflow.Context, event *shared.PushEventSignal, branchName string, lastBranchCommit string,
) error {
	logger := workflow.GetLogger(ctx)
	repoID := strconv.FormatInt(event.RepoID, 10)
	// Sleep for 5 days before raising stale detection
	_ = workflow.Sleep(ctx, 5*24*time.Hour)
	// _ = workflow.Sleep(ctx, 30*time.Second)

	logger.Info("Stale branch detection", "woke up from sleep", "checking for stale branch")

	rpa := Instance().RepoProvider(RepoProvider(event.RepoProvider))
	mpa := Instance().MessageProvider(MessageProviderSlack) // TODO - maybe not hardcode to slack and get from payload

	providerActOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	pctx := workflow.WithActivityOptions(ctx, providerActOpts)

	commit := &LatestCommit{}
	if err := workflow.ExecuteActivity(pctx, rpa.GetLatestCommit, repoID, branchName).Get(ctx, &commit); err != nil {
		logger.Error("Repo provider activities: Get latest commit activity", "error", err)
		return err
	}

	// check if the branchName branch has the lastBranchCommit as the latest commit
	if lastBranchCommit == commit.SHA {
		// get the teamID from repo table
		teamID := ""
		if err := workflow.ExecuteActivity(pctx, rpa.GetRepoTeamID, repoID).Get(ctx, &teamID); err != nil {
			logger.Error("Repo provider activities: Get repo TeamID activity", "error", err)
			return err
		}

		if err := workflow.ExecuteActivity(pctx, mpa.SendStaleBranchMessage, teamID, commit).Get(ctx, nil); err != nil {
			logger.Error("Message provider activities: Send stale branch message activity", "error", err)
			return err
		}

		return nil
	}

	// at this point, the branch is not stale so just return
	logger.Info("stale branch NOT detected")

	return nil
}

func (w *Workflows) PollMergeQueue(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("PollMergeQueue", "entry", "workflow started")

	// wait for github action to return success status
	ch := workflow.GetSignalChannel(ctx, shared.MergeQueueStarted.String())
	element := &shared.MergeQueueSignal{}
	ch.Receive(ctx, &element)

	logger.Debug("PollMergeQueue first signal received")
	logger.Info("PollMergeQueue", "data recvd", element)

	// actually merge now
	rpa := Instance().RepoProvider(RepoProvider(element.RepoProvider))
	providerActOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	pctx := workflow.WithActivityOptions(ctx, providerActOpts)

	// get list of all available github workflow actions/files
	if err := workflow.ExecuteActivity(pctx, rpa.GetAllRelevantActions, element.InstallationID, element.RepoName,
		element.RepoOwner).Get(ctx, nil); err != nil {
		logger.Error("error getting all labeled actions", "error", err)
		return err
	}

	logger.Debug("waiting on second signal now.")

	mergeSig := workflow.GetSignalChannel(ctx, shared.MergeTriggered.String())
	mergeSig.Receive(ctx, nil)

	logger.Debug("PollMergeQueue second signal received")

	if err := workflow.ExecuteActivity(pctx, rpa.RebaseAndMerge, element.RepoOwner, element.RepoName, element.Branch,
		element.InstallationID).Get(ctx, nil); err != nil {
		logger.Error("error rebasing & merging activity", "error", err)
		return err
	}

	logger.Info("github action triggered")

	return nil
}

// DeProvisionInfra de-provisions the infrastructure created for stack deployment.
func (w *Workflows) DeProvisionInfra(ctx workflow.Context, stackID string, resourceData *ResourceConfig) error {
	return nil
}

// GetAssets gets assets for stack including resources, workloads and blueprint.
func (w *Workflows) GetAssets(ctx workflow.Context, payload *GetAssetsPayload) error {
	var (
		future workflow.Future
		err    error = nil
	)

	logger := workflow.GetLogger(ctx)
	assets := NewAssets()
	workloads := SlicedResult[Workload]{}
	resources := SlicedResult[Resource]{}
	repos := SlicedResult[Repo]{}
	blueprint := Blueprint{}

	shared.Logger().Info("Get assets workflow")

	selector := workflow.NewSelector(ctx)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	providerActivityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
	}
	pctx := workflow.WithActivityOptions(ctx, providerActivityOpts)

	// get resources for stack
	future = workflow.ExecuteActivity(actx, activities.GetResources, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &resources); err != nil {
			logger.Error("GetResources providers failed", "error", err)
			return
		}
	})

	// get workloads for stack
	future = workflow.ExecuteActivity(actx, activities.GetWorkloads, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &workloads); err != nil {
			logger.Error("GetWorkloads providers failed", "error", err)
			return
		}
	})

	// get repos for stack
	future = workflow.ExecuteActivity(actx, activities.GetRepos, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &repos); err != nil {
			logger.Error("GetRepos providers failed", "error", err)
			return
		}
	})

	// get blueprint for stack
	future = workflow.ExecuteActivity(actx, activities.GetBluePrint, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &blueprint); err != nil {
			logger.Error("GetBluePrint providers failed", "error", err)
			return
		}
	})

	// TODO: come up with a better logic for this
	for i := 0; i < 4; i++ {
		selector.Select(ctx)
		// return if providers failed. TODO: handle race conditions as the 'err' variable is shared among all providers
		if err != nil {
			logger.Error("Exiting due to providers failure")
			return err
		}
	}

	// Tag the build image with changeset ID
	// TODO: update it to tag one or multiple images against every workload
	switch payload.ImageRegistry {
	case "GCPArtifactRegistry":
		err := workflow.ExecuteActivity(actx, activities.TagGcpImage, payload.Image, payload.Digest, payload.ChangeSetID).Get(ctx, nil)
		if err != nil {
			return err
		}
	default:
		shared.Logger().Error("This image registry is not supported in quantm yet", "registry", payload.ImageRegistry)
	}

	// get commits against the repos
	repoMarker := make([]ChangeSetRepoMarker, len(repos.Data))

	for idx, repo := range repos.Data {
		marker := &repoMarker[idx]
		// p := Instance().Provider(repo.Provider) // get the specific provider
		p := Instance().RepoProvider(repo.Provider) // get the specific provider

		commit := LatestCommit{}
		if err := workflow.
			ExecuteActivity(pctx, p.GetLatestCommit, repo.ProviderID, repo.DefaultBranch).
			Get(ctx, &commit); err != nil {
			logger.Error("Error in getting latest commit ID", "repo", repo.Name, "provider", repo.Provider)
			return fmt.Errorf("Error in getting latest commit ID repo:%s, provider:%s", repo.Name, repo.Provider.String())
		}

		marker.CommitID = commit.SHA
		marker.HasChanged = repo.ID == payload.RepoID
		marker.Provider = repo.Provider.String()
		marker.RepoID = repo.ID.String()
		logger.Debug("Repo", "Name", repo.Name, "Repo marker", marker)
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

	assets.ChangesetID = payload.ChangeSetID
	assets.Blueprint = blueprint
	assets.Repos = append(assets.Repos, repos.Data...)
	assets.Resources = append(assets.Resources, resources.Data...)
	assets.Workloads = append(assets.Workloads, workloads.Data...)
	logger.Debug("Assets retrieved", "Assets", assets)

	// signal parent workflow
	parent := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.
		SignalExternalWorkflow(ctx, parent, "", WorkflowSignalAssetsRetrieved.String(), assets).
		Get(ctx, nil)

	return nil
}

// ProvisionInfra provisions the infrastructure required for stack deployment.
func (w *Workflows) ProvisionInfra(ctx workflow.Context, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	futures := make([]workflow.Future, 0)

	shared.Logger().Debug("provision infra", "assets", assets)

	for _, rsc := range assets.Resources {
		logger.Info("Creating resource", "Name", rsc.Name)

		// get the resource constructor specific to the driver e.g gke, cloudrun for GCP, sns, fargate for AWS
		resconstr := Instance().ResourceConstructor(rsc.Provider, rsc.Driver)

		if rsc.IsImmutable {
			// assuming a single region for now
			region := getRegion(rsc.Provider, &assets.Blueprint)
			providerConfig := getProviderConfig(rsc.Provider, &assets.Blueprint)
			r, err := resconstr.Create(rsc.Name, region, rsc.Config, providerConfig)

			if err != nil {
				logger.Error("could not create resource object", "ID", rsc.ID, "name", rsc.Name, "Error", err)
				return err
			}

			// resource is an interface and cannot be sent as a parameter in workflow because workflow cannot unmarshal an interface.
			// So we need to send the marshalled value in this workflow and then unmarshal and resconstruct the resource again in the
			// Deploy workflow
			ser, err := r.Marshal()
			if err != nil {
				logger.Error("could not marshal resource", "ID", rsc.ID, "name", rsc.Name, "Error", err)
				return err
			}

			assets.Infra[rsc.ID] = ser

			// TODO: initiate all resource provisions in parallel in child workflows and wait for all child workflows
			// completion before sending infra provisioned signal
			if f, err := r.Provision(ctx); err != nil {
				logger.Error("could not start resource provisioning", "ID", rsc.ID, "name", rsc.Name, "Error", err)

				return err
			} else if f != nil {
				futures = append(futures, f)
			}
		}
	}

	// not tested yet as provision workflow for cloudrun doesn't return any future.
	// the idea is to run all provision workflows asynchronously and in parallel and wait here for their completion before moving forward
	// also need to test if we can call child.Get with parent workflow's context
	for _, child := range futures {
		_ = child.Get(ctx, nil)
	}

	shared.Logger().Info("Signaling infra provisioned")

	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalInfraProvisioned.String(), assets).Get(ctx, nil)

	return nil
}

// Deploy deploys the stack.
func (w *Workflows) Deploy(ctx workflow.Context, stackID string, lock *mutex.Handler, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Deployment initiated", "changeset", assets.ChangesetID, "infra", assets.Infra)
	infra := make(Infra)

	// Acquire lock
	err := lock.Acquire(ctx)
	if err != nil {
		logger.Error("Error in acquiring lock", "Error", err)
		return err
	}

	// create deployable, map of one or more workloads against each resource
	deployables := make(map[gocql.UUID][]Workload) // map of resource id and workloads
	for _, w := range assets.Workloads {
		_, ok := deployables[w.ResourceID]
		if !ok {
			deployables[w.ResourceID] = make([]Workload, 0)
		}

		deployables[w.ResourceID] = append(deployables[w.ResourceID], w)
	}

	// create the resource object again from marshaled data and deploy workload on resource
	for _, rsc := range assets.Resources {
		resconstr := Instance().ResourceConstructor(rsc.Provider, rsc.Driver)
		inf := assets.Infra[rsc.ID] // get marshaled resource from ID
		r := resconstr.CreateFromJson(inf)
		infra[rsc.ID] = r
		_ = r.Deploy(ctx, deployables[rsc.ID], assets.ChangesetID)
	}

	// update traffic on resource from 50 to 100
	var i int32
	for i = 50; i <= 100; i += 25 {
		for id, r := range infra {
			logger.Info("updating traffic", id, r)
			_ = r.UpdateTraffic(ctx, i)
		}
	}

	_ = lock.Release(ctx)

	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalDeploymentCompleted.String(), assets)

	return nil
}

func onManualOverrideSignal(ctx workflow.Context, _ string, _ Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	triggerID := int64(0)

	return func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &triggerID)
		logger.Info("manual override for", "Trigger ID", triggerID)
	}
}

// onDeploymentStartedSignal is the channel handler for trigger channel
// It will execute GetAssets and update PR deployment state to "GettingAssets".
func onDeploymentStartedSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		// Receive signal data
		payload := &shared.PullRequestSignal{}
		channel.Receive(ctx, payload)
		logger.Info("received deployment request", "Trigger ID", payload.TriggerID)

		// We want to filter workflows with changeset ID, so create changeset ID here and use it for creating workflow ID
		changesetID, _ := gocql.RandomUUID()

		// Set child workflow options
		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowElement("get_assets"),
				shared.WithWorkflowMod("trigger"),
				shared.WithWorkflowModID(strconv.FormatInt(payload.TriggerID, 10)),
			)

		getAssetsPayload := &GetAssetsPayload{
			StackID:       stackID,
			RepoID:        payload.RepoID,
			ChangeSetID:   changesetID,
			Digest:        payload.Digest,
			ImageRegistry: payload.ImageRegistry,
			Image:         payload.Image,
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
		deployment := NewDeployment()
		deployments[changesetID] = deployment
		deployment.state = GettingAssets
		deployment.workflows.GetAssets = execution.ID
	}
}

// onAssetsRetrievedSignal will receive assets sent by GetAssets, update deployment state and execute ProvisionInfra.
func onAssetsRetrievedSignal(ctx workflow.Context, _ string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("received Assets", "changeset", assets.ChangesetID)

		// update deployment state
		deployment := deployments[assets.ChangesetID]
		deployment.state = GotAssets

		// execute provision infra workflow
		logger.Info("Executing provision Infra workflow")

		var execution workflow.Execution

		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowBlock("changeset"), // TODO: shouldn't this be part of the changeset controller?
				shared.WithWorkflowBlockID(assets.ChangesetID.String()),
				shared.WithWorkflowElement("provision_infra"),
			)

		ctx = workflow.WithChildOptions(ctx, opts)

		err := workflow.
			ExecuteChildWorkflow(ctx, w.ProvisionInfra, assets).
			GetChildWorkflowExecution().Get(ctx, &execution)

		if err != nil {
			logger.Error("TODO: Error in executing ProvisionInfra", "Error", err)
		}

		logger.Info("Executed provision Infra workflow")

		deployment.state = ProvisioningInfra
		deployment.workflows.ProvisionInfra = execution.ID
	}
}

// onInfraProvisionedSignal will receive assets by ProvisionInfra, update deployment state and execute Deploy.
func onInfraProvisionedSignal(
	ctx workflow.Context, stackID string, lock mutex.Mutex, deployments Deployments, _ Infra,
) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("Infra provisioned", "changeset", assets.ChangesetID)

		deployment := deployments[assets.ChangesetID]
		deployment.state = InfraProvisioned

		// deployment.OldInfra = activeinfra  // All traffic is currently being routed to this infra
		deployment.NewInfra = assets.Infra // handling zero traffic, no workload is deployed

		var execution workflow.Execution

		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowBlock("changeset"), // TODO: shouldn't this be part of the changeset controller?
				shared.WithWorkflowBlockID(assets.ChangesetID.String()),
				shared.WithWorkflowElement("deploy"),
			)
		ctx = workflow.WithChildOptions(ctx, opts)

		err := workflow.
			ExecuteChildWorkflow(ctx, w.Deploy, stackID, lock.(*mutex.Handler), assets).
			GetChildWorkflowExecution().Get(ctx, &execution)
		if err != nil {
			logger.Error("Error in Executing deployment workflow", "Error", err)
		}

		deployment.state = CreatingDeployment
		deployment.workflows.ProvisionInfra = execution.ID
	}
}

// onDeploymentCompletedSignal will conclude the deployment.
func onDeploymentCompletedSignal(ctx workflow.Context, _ string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("Deployment complete", "changeset", assets.ChangesetID)
		delete(deployments, assets.ChangesetID)

		logger.Info("Deleted deployment data", "changeset", assets.ChangesetID)
	}
}
