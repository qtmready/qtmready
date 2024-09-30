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

package kernel

import (
	"context"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	// MessageIO defines the interface for sending various types of messages.
	MessageIO interface {
		// SendStaleBranchMessage sends a message indicating a stale branch.
		SendStaleBranchMessage(ctx context.Context, payload *defs.MessageIOStaleBranchPayload) error

		// NotifyLinesExceed sends a message indicating a line exceed message.
		NotifyLinesExceed(ctx context.Context, event *defs.Event[defs.LinesExceed, defs.RepoProvider]) error

		// NotifyMergeConflict sends a message indicating a merge conflict.
		NotifyMergeConflict(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) error
	}
)
