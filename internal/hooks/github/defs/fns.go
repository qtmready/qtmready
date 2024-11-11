package githubdefs

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// TODO - move the events cast.
func NormalizeCommit(c Commit) *eventsv1.Commit {
	return &eventsv1.Commit{
		Sha:       c.ID,
		Message:   c.Message,
		Url:       c.URL,
		Timestamp: timestamppb.New(c.Timestamp),
		Added:     c.Added,
		Removed:   c.Removed,
		Modified:  c.Modified,
	}
}
