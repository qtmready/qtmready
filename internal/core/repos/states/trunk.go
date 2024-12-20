package states

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	MergeRequest struct {
		Number int64
		Branch string
	}

	Trunk struct {
		*Base
		MergeQueue *Sequencer[int64, MergeRequest]
		IsEmpty    bool
	}
)

// OnLabel handles the pull request event with label on the repository.
func (state *Trunk) OnLabel(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		label := &events.Event[eventsv1.RepoHook, eventsv1.PullRequestLabel]{}
		state.rx(ctx, rx, label)

		request := &MergeRequest{label.Payload.GetNumber(), label.Payload.GetBranch()}

		if label.Context.Action == events.ActionCreated {
			if label.Payload.GetName() == "quantm-merge" {
				state.MergeQueue.Push(ctx, label.Payload.GetNumber(), request)
			}

			if label.Payload.GetName() == "quantm-priority" {
				state.MergeQueue.Priority(ctx, label.Payload.GetNumber(), request)
			}
		} else {
			state.MergeQueue.Remove(ctx, label.Payload.GetNumber())
		}
	}
}

// - queue process -.
func (state Trunk) ProcessItem()

func (state *Trunk) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
	state.MergeQueue.Init(ctx)
}

func NewTrunk(repo *entities.Repo, chat *entities.ChatLink) *Trunk {
	return &Trunk{
		&Base{Repo: repo, ChatLink: chat},
		NewSequencer[int64, MergeRequest](),
		true,
	}
}
