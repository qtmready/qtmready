// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package code

import (
	"github.com/gocql/gocql"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	// BranchTriggers maps branch names to their corresponding triggering event IDs.
	//
	// This data structure facilitates event lineage tracing by providing the root event for each branch.
	BranchTriggers map[string]gocql.UUID

	// StashedPushEvents stores events that are awaiting processing.
	//
	// Events are typically stashed when the associated branch does not yet exist or the event requires a parent event
	// (e.g., a push event needing a branch creation event) that has not yet been received. This scenario can arise due to
	// the distributed nature of event arrival.
	StashedPushEvents[P defs.RepoProvider] map[string][]*defs.Event[defs.Push, P]
)

// add associates a branch with its triggering event ID.
func (b BranchTriggers) add(branch string, id gocql.UUID) {
	b[branch] = id
}

// clear removes the association between a branch and its triggering event ID.
func (b BranchTriggers) clear(branch string) {
	delete(b, branch)
}

// get retrieves the event ID associated with a branch.
//
// Returns the event ID and a boolean indicating whether the branch exists.
func (b BranchTriggers) get(branch string) (gocql.UUID, bool) {
	id, ok := b[branch]

	return id, ok
}

// push adds an event to the stash for the specified branch.
func (s StashedPushEvents[P]) push(branch string, event *defs.Event[defs.Push, P]) {
	if _, ok := s[branch]; !ok {
		s[branch] = make([]*defs.Event[defs.Push, P], 0)
	}

	s[branch] = append(s[branch], event)
}

// pop retrieves and removes the oldest event from the stash for the specified branch.
//
// Returns the event and a boolean indicating whether an event was present.
func (s StashedPushEvents[P]) all(branch string) ([]*defs.Event[defs.Push, P], bool) {
	events, ok := s[branch]
	if !ok || len(events) == 0 {
		return events, false
	}

	return events, true
}

func (s StashedPushEvents[P]) clear(branch string) {
	delete(s, branch)
}
