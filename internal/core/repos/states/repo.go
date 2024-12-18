package states

import (
	"errors"

	"github.com/google/uuid"
	"go.breu.io/durex/dispatch"
	"go.breu.io/durex/queues"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/pulse"
)

type (
	// Repo defines the state for Repo Workflows. It embeds BaseState to inherit its functionality.
	Repo struct {
		*Base    `json:"base"`  // Base workflow state.
		Triggers BranchTriggers `json:"triggers"` // Branch triggers.

		acts *activities.Repo
	}
)

// - signal handlers -

// OnPush handles the push event on the repository. If the branch is the default branch, the event is forwarded to all
// branches with a rebase instruction. Otherwise, the event is forwarded to the branch.
//
// TODO: Define a new event type for rebase events.
func (state *Repo) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, rx, push)

		branch := fns.BranchNameFromRef(push.Payload.Ref)

		if branch == state.Repo.DefaultBranch {
			state.attempt_rebase(ctx, push)

			return
		}

		state.Triggers.add(branch, push.ID)

		if err := state.forward_to_branch(ctx, defs.SignalPush, branch, push); err != nil {
			state.logger.Warn("push: unable to signal branch", "repo", state.Repo.ID, "branch", branch, "error", err.Error())
		}
	}
}

func (state *Repo) OnRef(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		ref := &events.Event[eventsv1.RepoHook, eventsv1.GitRef]{}
		state.rx(ctx, rx, ref)

		if ref.Payload.Kind == "branch" {
			branch := fns.BranchNameFromRef(ref.Payload.Ref)

			if err := state.forward_to_branch(ctx, defs.SignalRef, branch, ref); err != nil {
				state.logger.Warn("ref: unable to signal branch", "repo", state.Repo.ID, "branch", branch, "error", err.Error())
			}

			if ref.Context.Action == events.ActionCreated {
				state.Triggers.add(branch, ref.ID)
			}

			if ref.Context.Action == events.ActionDeleted {
				state.Triggers.remove(branch)
			}
		}
	}
}

// OnPR handles the pull request event on the repository.
func (state *Repo) OnPR(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		pr := &events.Event[eventsv1.RepoHook, eventsv1.PullRequest]{}
		state.rx(ctx, rx, pr)
	}
}

// OnLabel handles the pull request event with label on the repository.
func (state *Repo) OnLabel(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		label := &events.Event[eventsv1.RepoHook, eventsv1.PullRequestLabel]{}
		state.rx(ctx, rx, label)
	}
}

// OnLabel handles the pull request event with label on the repository.
func (state *Repo) OnPrReview(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		label := &events.Event[eventsv1.RepoHook, eventsv1.PullRequestReview]{}
		state.rx(ctx, rx, label)
	}
}

// - query handlers -

// QueryBranchTrigger queries the parent branch for the specified branch.
func (state *Repo) QueryBranchTrigger(branch string) (uuid.UUID, error) {
	id, ok := state.Triggers.get(branch)
	if ok {
		return id, nil
	}

	return uuid.Nil, errors.New("branch not found")
}

// - local -

// forward_to_branch routes the signal to the appropriate branch.
func (state *Repo) forward_to_branch(ctx workflow.Context, signal queues.Signal, branch string, event any) error {
	ctx = dispatch.WithDefaultActivityContext(ctx)

	next := NewBranch(state.Repo, state.ChatLink, branch)
	payload := &defs.SignalBranchPayload{Signal: signal, Repo: state.Repo, Branch: branch}

	return workflow.ExecuteActivity(ctx, state.acts.ForwardToBranch, payload, event, next).Get(ctx, nil)
}

// forward_to_trunk routes the signal to the trunk.
func (state *Repo) forward_to_trunk(ctx workflow.Context, signal queues.Signal, event any) error {
	ctx = dispatch.WithDefaultActivityContext(ctx)

	next := NewTrunk(state.Repo, state.ChatLink)
	payload := &defs.SignalTrunkPayload{Signal: signal, Repo: state.Repo}

	return workflow.ExecuteActivity(ctx, state.acts.ForwardToTrunk, payload, event, next).Get(ctx, nil)
}

// attempt_rebase rebases all branches with a trigger on the default branch.
func (state *Repo) attempt_rebase(ctx workflow.Context, push *events.Event[eventsv1.RepoHook, eventsv1.Push]) {
	for branch := range state.Triggers {
		workflow.Go(ctx, func(ctx workflow.Context) {
			rebase := events.
				Next[eventsv1.RepoHook, eventsv1.Push, eventsv1.Rebase](push, events.ScopeRebase, events.ActionRequested).
				SetPayload(&eventsv1.Rebase{Base: branch, Head: push.Payload.After, Repository: push.Payload.Repository})

			if err := pulse.Persist(ctx, rebase); err != nil {
				state.logger.Warn(
					"attempt_rebase: unable to persist rebase event",
					"repo", state.Repo.ID, "branch", branch, "error", err.Error(),
				)
			}

			if err := state.forward_to_branch(ctx, defs.SignalRebase, branch, rebase); err != nil {
				state.logger.Warn("atemmpt_rebase: unable to signal branch", "repo", state.Repo.ID, "branch", branch, "error", err.Error())
			}
		})
	}
}

// - state managers -

func (state *Repo) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
}

// NewRepo creates a new RepoState instance. It initializes BaseState using the provided context and
// hydrated repository data.
func NewRepo(repo *entities.Repo, chat *entities.ChatLink) *Repo {
	base := &Base{Repo: repo, ChatLink: chat}
	triggers := make(BranchTriggers)

	return &Repo{base, triggers, &activities.Repo{}} // Return new RepoState instance.
}
