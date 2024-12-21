package states

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Trunk struct {
		*Base
		MergeQueue *Sequencer[int64, eventsv1.MergeQueue] `json:"merge_queue"`

		done bool
	}
)

// - queue process -.
func (state Trunk) OnMergeQueue(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		mq := &events.Event[eventsv1.RepoHook, eventsv1.MergeQueue]{}
		state.rx(ctx, rx, mq)

		if mq.Context.Action == events.EventActionRemoved {
			state.MergeQueue.Remove(ctx, mq.Payload.GetNumber())

			return
		}

		if mq.Payload.IsPriority {
			state.MergeQueue.Priority(ctx, mq.Payload.GetNumber(), mq.Payload)

			return
		}

		state.MergeQueue.Push(ctx, mq.Payload.GetNumber(), mq.Payload)
	}
}

func (state *Trunk) StartQueueProcessor(ctx workflow.Context) {}

func (state *Trunk) Continue() bool {
	return !state.done
}

func (state *Trunk) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
	state.MergeQueue.Init(ctx)
}

func NewTrunk(repo *entities.Repo, chat *entities.ChatLink) *Trunk {
	return &Trunk{
		&Base{Repo: repo, ChatLink: chat},
		NewSequencer[int64, eventsv1.MergeQueue](),
		false,
	}
}
