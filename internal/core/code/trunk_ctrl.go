// Copyright © 2024, Breu, Inc. <info@breu.io>
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
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

// TrunkCtrl is the event loop to process events during the lifecycle of the main branch.
//
// It processes the following events:
//
//   - push
//   - create_delete
func TrunkCtrl(ctx workflow.Context, repo *defs.Repo) error {
	state := NewTrunkState(ctx, repo)
	selector := workflow.NewSelector(ctx)

	// channels
	// push event
	push := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// create_delete
	create_delete := workflow.GetSignalChannel(ctx, defs.RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	// main event loop
	for state.is_active() {
		selector.Select(ctx)

		if state.needs_reset() {
			return state.as_new(ctx, "event history exceeded threshold", TrunkCtrl, repo)
		}
	}

	// graceful shutdown
	state.terminate(ctx)

	return nil
}
