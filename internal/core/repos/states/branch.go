package states

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/cast"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/durable/periodic"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/pulse"
)

type (
	BranchIntervals struct {
		pr    periodic.Interval // used to send a notifcation if a pr is not opened within a certain time.
		stale periodic.Interval // used to send a notification if a branch is stale.
	}

	Branch struct {
		*Base `json:"base"` // Base workflow state.

		Branch       string           `json:"branch"`
		LatestCommit *eventsv1.Commit `json:"latest_commit"`

		intervals BranchIntervals
		acts      *activities.Branch
		notify    *activities.Notify
		done      bool
	}
)

// PullRequestMonitor is a goroutine that monitors the branch for pull requests. If a pull request is not opened
// within a certain time, a notification is to the hook associated with the branch.
//
// TODO: implement the logic for sending a notification if a pr is not opened within a certain time.
func (state *Branch) PullRequestMonitor(ctx workflow.Context) {
	workflow.Go(ctx, func(ctx_ workflow.Context) {
		for {
			state.intervals.pr.Tick(ctx_)
			_ = state.notify_user(ctx)
		}
	})
}

// StaleMonitor is a goroutine that monitors the branch for staleness. If the branch is stale, a notification is
// sent to the hook associated with the branch.
//
// TODO: implement the logic for sending a notification if the branch is stale.
func (state *Branch) StaleMonitor(ctx workflow.Context) {
	workflow.Go(ctx, func(ctx_ workflow.Context) {
		for {
			state.intervals.stale.Tick(ctx_)
			_ = state.notify_user(ctx)
		}
	})
}

// OnPush resets the stale timer and processes the push event. The repo is cloned, the diff calculated, and
// notifications sent if change complexity warrants. Author notification is prioritized, falling back to
// the repo's chat hook.
func (state *Branch) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(ch workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, ch, event)

		state.intervals.stale.Reset(ctx)

		opts := &workflow.SessionOptions{ExecutionTimeout: time.Minute * 30, CreationTimeout: time.Second * 30}

		session, err := workflow.CreateSession(ctx, opts)
		if err != nil {
			state.logger.Error("clone: unable to create session", "push", event.Payload.After, "error", err.Error())
			return
		}

		defer workflow.CompleteSession(session)

		state.LatestCommit = fns.GetLatestCommit(event.Payload)

		clone := &defs.ClonePayload{Repo: state.Repo, Hook: event.Context.Hook, Branch: state.Branch, SHA: event.Payload.After}
		path := state.clone(session, clone)
		diff := state.diff(session, path, state.Repo.DefaultBranch, event.Payload.After)
		state.remove_dir(ctx, path)

		// compare the diff
		state.compare_diff(session, event, diff)
	}
}

// OnRebase handles the rebase event for the branch. It creates a session, clones the repository at given branch,
// attempts to rebase the branch with given sha, and removes the cloned repository.
func (state *Branch) OnRebase(ctx workflow.Context) durable.ChannelHandler {
	return func(ch workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Rebase]{}
		state.rx(ctx, ch, event)

		opts := &workflow.SessionOptions{ExecutionTimeout: time.Minute * 30, CreationTimeout: time.Second * 30}

		session, err := workflow.CreateSession(ctx, opts)
		if err != nil {
			state.logger.Error("clone: unable to create session", "rebase", event.Payload.Head, "error", err.Error())
			return
		}

		defer workflow.CompleteSession(session)

		clone := &defs.ClonePayload{Repo: state.Repo, Hook: event.Context.Hook, Branch: state.Branch, SHA: event.Payload.Head}
		path := state.clone(session, clone)

		rebase := &defs.RebaseResult{}
		_ = state.run(ctx, "rebase", state.acts.Rebase, &defs.RebasePayload{Rebase: event.Payload, Path: path}, rebase)

		state.check_merge_conflict(session, event, rebase)

		state.remove_dir(ctx, path)
	}
}

// OnLabel handles pull request label events.
func (state *Branch) OnLabel(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.PullRequestLabel]{}
		state.rx(ctx, rx, event)

		switch event.Payload.Name {
		case "qmerge":
			fmt.Println("push to the queue based in level of priority")
		case "priority-qmerge":
			fmt.Println("push to the queue based in level of priority")
		default:
			return
		}
	}
}

// OnPrReview handles pull request review events.
func (state *Branch) OnPrReview(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.PullRequestReview]{}
		state.rx(ctx, rx, event)
	}
}

// OnPRReviewComment handles pull request review comment events.
func (state *Branch) OnPRReviewComment(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.PullRequestReview]{}
		state.rx(ctx, rx, event)
	}
}

// ExitLoop returns true if the branch should exit the event loop.
func (state *Branch) ExitLoop(ctx workflow.Context) bool {
	return state.done || workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// Init initializes the branch state.
func (state *Branch) Init(ctx workflow.Context) {
	state.Base.Init(ctx)

	pr := periodic.New(ctx, time.Minute*60*24)
	stale := periodic.New(ctx, time.Minute*60*24)

	state.intervals = BranchIntervals{pr: pr, stale: stale}
}

// clone clones the repository at the given SHA using a Temporal activity.  A UUID is generated for the clone path via SideEffect
// to ensure idempotency. Returns the clone path.
func (state *Branch) clone(ctx workflow.Context, payload *defs.ClonePayload) string {
	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any { return uuid.New().String() }).Get(&payload.Path)

	if err := state.run(ctx, "clone", state.acts.Clone, payload, &payload.Path); err != nil {
		state.logger.Error("clone: unable to clone", "error", err.Error())
	}

	return payload.Path
}

func (state *Branch) remove_dir(ctx workflow.Context, path string) {
	if err := state.run(ctx, "remove", state.acts.RemoveDir, path, nil); err != nil {
		state.logger.Error("remove: unable to remove directory", "error", err.Error())
	}
}

// diff calculates the diff between the given base and SHA using a Temporal activity.  Returns the diff result.
func (state *Branch) diff(ctx workflow.Context, path, base, sha string) *eventsv1.Diff {
	payload := &defs.DiffPayload{Path: path, Base: base, SHA: sha}
	result := &eventsv1.Diff{}

	if err := state.run(ctx, "diff", state.acts.Diff, payload, result); err != nil {
		state.logger.Error("diff: unable to calculate diff", "error", err.Error())
	}

	return result
}

// check the change diff and if it exceed from the threshold sends message to user other wise message to repo connected group.
func (state *Branch) compare_diff(
	ctx workflow.Context, push *events.Event[eventsv1.RepoHook, eventsv1.Push], diff *eventsv1.Diff,
) {
	dlt := diff.GetLines().GetAdded() + diff.GetLines().GetRemoved()

	if dlt > state.Repo.Threshold {
		// check the repo's connected chat or user's connected chat.
		hook := int32(eventsv1.ChatHook_CHAT_HOOK_SLACK)
		event := cast.PushEventToDiffEvent(push, hook, diff)

		// persist chat event
		if err := pulse.Persist(ctx, event); err != nil {
			state.logger.Warn(
				"attempt_diff: unable to persist diff event",
				"repo", state.Repo.ID, "branch", fns.BranchNameFromRef(push.Payload.Ref), "error", err.Error(),
			)
		}

		if err := state.run(ctx, "line_exceed", state.notify.LinesExceeded, event, nil); err != nil {
			state.logger.Error("lines_exceed: unable to to send", "error", err.Error())
		}
	}
}

// check_merge_conflict check the merge conflict and send chat message otherwise nothing.
func (state *Branch) check_merge_conflict(
	ctx workflow.Context, rebase *events.Event[eventsv1.RepoHook, eventsv1.Rebase], res *defs.RebaseResult,
) {
	if len(res.Conflicts) > 0 {
		// check the repo's connected chat or user's connected chat.
		hook := int32(eventsv1.ChatHook_CHAT_HOOK_SLACK)

		// TODO - head and base commits
		payload := &eventsv1.Merge{
			HeadBranch: rebase.Payload.Head,
			BaseBranch: rebase.Payload.Base,
			Files:      res.Conflicts,
		}

		event := cast.RebaseEventToMergeConflictEvent(rebase, hook, payload)

		// persist chat event
		if err := pulse.Persist(ctx, event); err != nil {
			state.logger.Warn(
				"attempt_merge: unable to persist merge event",
				"repo", state.Repo.ID, "branch", payload.HeadBranch, "error", err.Error(),
			)
		}

		if err := state.run(ctx, "merge_conflict", state.notify.MergeConflict, event, nil); err != nil {
			state.logger.Error("merge_conflict: unable to to send", "error", err.Error())
		}
	}
}

func (state *Branch) notify_user(_ workflow.Context) error { return nil }

// NewBranch constructs a new Branch state.
func NewBranch(repo *entities.Repo, chat *entities.ChatLink, branch string) *Branch {
	base := &Base{Repo: repo, ChatLink: chat}

	return &Branch{Base: base, Branch: branch, acts: &activities.Branch{}, notify: &activities.Notify{}}
}
