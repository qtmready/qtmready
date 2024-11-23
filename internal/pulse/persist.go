package pulse

import (
	"context"
	"fmt"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

const (
	statement__events__persist = `
INSERT INTO events_%s VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
`
)

// Persist persists an event to the database. It dispatches to the correct activity based
// on the event's hook.
func Persist[H events.Hook, P events.Payload](ctx workflow.Context, event *events.Event[H, P]) error {
	ctx = dispatch.WithDefaultActivityContext(ctx)
	flat := event.Flatten()

	var future workflow.Future

	switch any(flat.Hook).(type) {
	case eventsv1.RepoHook:
		future = workflow.ExecuteActivity(ctx, PersistRepoEvent, flat)
	case eventsv1.MessagingHook:
		future = workflow.ExecuteActivity(ctx, PersistMessagingEvent, flat)
	}

	return future.Get(ctx, nil)
}

func PersistRepoEvent(ctx context.Context, event events.Flat[eventsv1.RepoHook]) error {
	slug, err := db.Queries().GetOrgSlugByID(ctx, event.OrgID)
	if err != nil {
		return nil
	}

	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__persist, table)

	return Instance().Connection().Exec(
		ctx,
		stmt,
		event.Version,
		event.ID,
		event.ParentID,
		event.Hook,
		event.Scope,
		event.Action,
		event.Source,
		event.SubjectID,
		event.SubjectName,
		event.UserID,
		event.TeamID,
		event.OrgID,
		event.Timestamp,
	)
}

func PersistMessagingEvent(ctx context.Context, event events.Flat[eventsv1.MessagingHook]) error {
	slug, err := db.Queries().GetOrgSlugByID(ctx, event.OrgID)
	if err != nil {
		return nil
	}

	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__persist, table)

	return Instance().Connection().Exec(
		ctx,
		stmt,
		event.Version,
		event.ID,
		event.ParentID,
		event.Hook,
		event.Scope,
		event.Action,
		event.Source,
		event.SubjectID,
		event.SubjectName,
		event.UserID,
		event.TeamID,
		event.OrgID,
		event.Timestamp,
	)
}
