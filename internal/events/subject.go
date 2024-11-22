package events

import (
	"github.com/google/uuid"
)

type (
	// Subject represents the entity within the system that is the subject of an event.
	//
	// It encapsulates:
	//   - ID: The unique identifier of the entity i.e. the primary key within its respective database table.
	//   - Name: The name of the entity's corresponding database table. This provides context for the event's subject.
	//     For instance, an event related to a branch would have "repos" as the subject name, as branches are associated
	//     with repositories.
	//   - TeamID: The unique identifier of the team to which this entity belongs. This allows for team-based filtering
	//     and organization
	//     of events.
	Subject struct {
		Name   string    `json:"name"`    // Name of the database table.
		ID     uuid.UUID `json:"id"`      // ID is the ID of the subject.
		OrgID  uuid.UUID `json:"org_id"`  // OrgID is the ID of the organization that the subject belongs to.
		TeamID uuid.UUID `json:"team_id"` // Team ID of the subject's team in the organization. It can be null uuid.
		UserID uuid.UUID `json:"user_id"` // UserID is the ID of the user that the subject belongs to. It can be null uuid.
	}
)

const (
	SubjectNameRepos = "repos"
)
