package shared

import (
	"github.com/gocql/gocql"
)

// workflow types

type (
	WorkflowSignal string // WorkflowSignal is the name of a workflow signal.

	PullRequestSignal struct {
		RepoID           gocql.UUID
		SenderWorkflowID string
	}
)

// workflow signals
const (
	WorkflowSignalPullRequest WorkflowSignal = "pull_request"
)

/*
 * Methods for WorkflowSignal.
 */

func (w WorkflowSignal) String() string { return string(w) }
