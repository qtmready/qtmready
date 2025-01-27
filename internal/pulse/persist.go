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
INSERT INTO %s (
	version,
	id,
	parents,
	hook,
	scope,
	action,
	source,
	subject_id,
	subject_name,
	user_id,
	team_id,
	org_id,
	timestamp
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
)

// Persist persists an event to clickhouse, routing it to the appropriate activity handler based on the
// event's associated hook.  It's a workflow-scoped function, mandating execution immediately post-event creation.
func Persist[H events.Hook, P events.Payload](ctx workflow.Context, event *events.Event[H, P]) error {
	ctx = dispatch.WithDefaultActivityContext(ctx)
	flat := event.Flatten()

	var future workflow.Future

	switch any(flat.Hook).(type) {
	case eventsv1.RepoHook:
		future = workflow.ExecuteActivity(ctx, PersistRepoEvent, flat)
	case eventsv1.ChatHook:
		future = workflow.ExecuteActivity(ctx, PersistChatEvent, flat)
	}

	return future.Get(ctx, nil)
}

// PersistRepoEvent persists a repo event to the database.
func PersistRepoEvent(ctx context.Context, flat events.Flat[eventsv1.RepoHook]) error {
	slug, err := db.Queries().GetOrgSlugByID(ctx, flat.OrgID)
	if err != nil {
		return nil
	}

	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__persist, table)

	return Get().
		Connection().
		Exec(
			ctx,
			stmt,
			flat.Version,
			flat.ID,
			flat.Parents,
			flat.Hook.Number(),
			flat.Scope,
			flat.Action,
			flat.Source,
			flat.SubjectID,
			flat.SubjectName,
			flat.UserID,
			flat.TeamID,
			flat.OrgID,
			flat.Timestamp,
		)
}

// PersistChatEvent persists a chat event to the database.
func PersistChatEvent(ctx context.Context, flat events.Flat[eventsv1.ChatHook]) error {
	slug, err := db.Queries().GetOrgSlugByID(ctx, flat.OrgID)
	if err != nil {
		return nil
	}

	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__persist, table)

	return Get().
		Connection().
		Exec(
			ctx,
			stmt,
			flat.Version,
			flat.ID,
			flat.Parents,
			flat.Hook.Number(),
			flat.Scope,
			flat.Action,
			flat.Source,
			flat.SubjectID,
			flat.SubjectName,
			flat.UserID,
			flat.TeamID,
			flat.OrgID,
			flat.Timestamp,
		)
}
