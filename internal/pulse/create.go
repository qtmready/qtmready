package pulse

import (
	"context"
	"fmt"

	"github.com/gobeam/stringy"
)

const (
	statement__events__create = `
CREATE TABLE IF NOT EXISTS %s (
  version String,
  id UUID,
  parents Array(UUID),
  hook Int32,
  scope String,
  action String,
  source String,
  subject_id UUID,
  subject_name String,
  user_id UUID,
  team_id UUID,
  org_id UUID,
  timestamp DateTime
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (toStartOfWeek(timestamp), toStartOfMonth(timestamp), timestamp, id);
`
)

// table_name returns the table name for the given kind and slug.
func table_name(kind, slug string) string {
	table := fmt.Sprintf("%s_%s", kind, slug)

	return stringy.New(table).SnakeCase().Get()
}

func CreateEventsTable(ctx context.Context, slug string) error {
	table := table_name("events", slug)
	stmt := fmt.Sprintf(statement__events__create, table)

	return Get().Connection().Exec(ctx, stmt)
}
