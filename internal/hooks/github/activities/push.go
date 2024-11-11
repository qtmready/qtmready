package githubacts

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go.breu.io/quantm/internal/events"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Push groups all the activities required for the Github Push.
	Push struct{}
)

func (p *Push) ConvertToPushEvent(
	ctx context.Context, payload *githubdefs.Push,
) (*events.Event[commonv1.RepoHook, eventsv1.Push], error) {
	commits := make([]*eventsv1.Commit, len(payload.Commits))
	for i, c := range payload.Commits {
		commits[i] = githubdefs.NormalizeCommit(c)
	}

	pl := eventsv1.Push{
		Ref:        payload.Ref,
		Before:     payload.Before,
		After:      payload.After,
		Repository: payload.Repository.Name,
		SenderId:   payload.SenderID(),
		Timestamp:  timestamppb.New(time.Now()),
		Commits:    commits,
	}

	// popluate the quantum event
	params := &githubdefs.RepoEventPayload{
		InstallationID: payload.InstallationID(),
		RepoID:         payload.RepoID(),
		Action:         events.EventActionCreated,
		Scope:          events.EventScopePush,
	}

	event, err := PopulateRepoEvent[commonv1.RepoHook, eventsv1.Push](ctx, params)
	if err != nil {
		return nil, err
	}

	event.Payload = pl

	return event, nil
}
