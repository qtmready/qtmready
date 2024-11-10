package pulse

import (
	"fmt"

	"github.com/gobeam/stringy"
)

const (
	query__events__create = `
CREATE TABLE IF NOT EXISTS %s (
  version EventVersion,
  id UUID,
  parent_id UUID,
  provider String,
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

// create_table_name returns the table name for the given org slug.
func create_table_name(kind, slug string) string {
	table := fmt.Sprintf("%s_%s", kind, slug)

	return stringy.New(table).SnakeCase().Get()
}

func QueryCreateEventsTable(slug string) string {
	table := create_table_name("events", slug)
	return fmt.Sprintf(query__events__create, table)
}
