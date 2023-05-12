package shared

import (
	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"
)

// workflow types.
type (
	WorkflowSignal string // WorkflowSignal is the name of a workflow signal.

	PullRequestSignal struct {
		RepoID           gocql.UUID
		SenderWorkflowID string
		TriggerID        int64
	}

	FutureHandler  func(workflow.Future)               // FutureHandler is the signature of the future handler function.
	ChannelHandler func(workflow.ReceiveChannel, bool) // ChannelHandler is the signature of the channel handler function.
)

// workflow signals.
const (
	WorkflowSignalPullRequest WorkflowSignal = "pull_request"

	WorkflowMaxAttempts = 10
)

/*
 * Methods for WorkflowSignal.
 */

func (w WorkflowSignal) String() string { return string(w) }
