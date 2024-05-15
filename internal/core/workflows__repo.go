package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/shared"
)

type (
	RepoWorkflows struct{}
)

func CheckEarlyWarning(
	ctx workflow.Context, rpa RepoIO, mpa MessageIO, pushEvent *shared.PushEventSignal,
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

	if err := workflow.ExecuteActivity(pctx, rpa.DetectChange, installationID, repoName, repoOwner, defaultBranch, branchName).
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
func (w *RepoWorkflows) BranchController(ctx workflow.Context) error {
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

	wf := &RepoWorkflows{}
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

func (w *RepoWorkflows) StaleBranchDetection(
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

func (w *RepoWorkflows) PollMergeQueue(ctx workflow.Context) error {
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
