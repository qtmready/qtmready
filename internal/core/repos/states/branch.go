package states

import (
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Branch struct {
		*Base        `json:"base"`    // Base workflow state.
		Branch       string           `json:"branch"`
		LatestCommit *eventsv1.Commit `json:"latest_commit"`

		acts *activities.Branch
		done bool
	}
)

// OnPush clones the repo at the given SHA, calculates the diff, and sends a notification to desired messaging service
// if the complexity of the push is above a certain threshold.
func (state *Branch) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(ch workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, ch, event)

		opts := &workflow.SessionOptions{ExecutionTimeout: time.Minute * 30, CreationTimeout: time.Second * 30}

		session, err := workflow.CreateSession(ctx, opts)
		if err != nil {
			state.logger.Error("clone: unable to create session", "push", event.Payload.After, "error", err.Error())
			return
		}

		defer workflow.CompleteSession(session)

		path := state.clone(session, event)
		_ = state.diff(session, path, state.Repo.DefaultBranch, event.Payload.After)
	}
}

func (state *Branch) ExitLoop(ctx workflow.Context) bool {
	return state.done || workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// Init initializes the branch state.
func (state *Branch) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
}

// clone clones the repository at the given SHA using a Temporal activity.  A UUID is generated for the clone path via SideEffect
// to ensure idempotency. Returns the clone path.
func (state *Branch) clone(ctx workflow.Context, event *events.Event[eventsv1.RepoHook, eventsv1.Push]) string {
	state.LatestCommit = fns.GetLatestCommit(event.Payload)

	payload := &defs.ClonePayload{
		Repo:   state.Repo,
		Branch: state.Branch,
		Hook:   event.Context.Hook,
		SHA:    event.Payload.After,
	}

	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any { return uuid.New().String() }).Get(&payload.Path)

	if err := state.run(ctx, "clone", state.acts.Clone, payload, &payload.Path); err != nil {
		state.logger.Error("clone: unable to clone", "error", err.Error())
	}

	return payload.Path
}

// diff calculates the diff between the given base and SHA using a Temporal activity.  Returns the diff result.
func (state *Branch) diff(ctx workflow.Context, path, base, sha string) *defs.DiffResult {
	payload := &defs.DiffPayload{Path: path, Base: base, SHA: sha}
	result := &defs.DiffResult{}

	if err := state.run(ctx, "diff", state.acts.Diff, payload, result); err != nil {
		state.logger.Error("diff: unable to calculate diff", "error", err.Error())
	}

	return result
}

// NewBranch constructs a new Branch state.
func NewBranch(repo *entities.Repo, msg *entities.Messaging, branch string) *Branch {
	base := &Base{Repo: repo, Messaging: msg}

	return &Branch{Base: base, Branch: branch, acts: &activities.Branch{}}
}
