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

		// SendNumberOfLinesExceedMessage sends a message indicating the number of lines has been exceeded.
		SendNumberOfLinesExceedMessage(ctx context.Context, payload *defs.MessageIOLineExeededPayload) error

		// SendMergeConflictsMessage sends a message indicating merge conflicts.
		SendMergeConflictsMessage(ctx context.Context, payload *defs.MessageIOMergeConflictPayload) error
	}
)
