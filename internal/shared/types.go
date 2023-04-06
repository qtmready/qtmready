package shared

import (
	"encoding/json"

	"github.com/gocql/gocql"
)

// workflow types

type (
	WorkflowSignal        string                    // WorkflowSignal is the name of a workflow signal.
	WorkflowSignalMapType map[string]WorkflowSignal // WorkflowSignalMap maps strings to their respective signal.

	PullRequestSignal struct {
		RepoID           gocql.UUID
		SenderWorkflowID string
	}
)

// Workflow signal types.
const (

	// github signals
	GithubWorkflowSignalInstallationEvent    WorkflowSignal = "installation_event"
	GithubWorkflowSignalCompleteInstallation WorkflowSignal = "complete_installation"
	GithubWorkflowSignalPullRequestProcessed WorkflowSignal = "pull_request_processed"

	// core signals
	CoreWorkflowSignalPullRequest  WorkflowSignal = "pull_request"
	CoreWorkflowSignalLockAcquired WorkflowSignal = "lock_acquired"
	CoreWorkflowSignalRequestLock  WorkflowSignal = "request_lock"
	CoreWorkflowSignalReleaseLock  WorkflowSignal = "release_lock"
)

var (
	WorkflowSignalMap = WorkflowSignalMapType{
		GithubWorkflowSignalInstallationEvent.String():    GithubWorkflowSignalInstallationEvent,
		GithubWorkflowSignalCompleteInstallation.String(): GithubWorkflowSignalCompleteInstallation,
		GithubWorkflowSignalPullRequestProcessed.String(): GithubWorkflowSignalPullRequestProcessed,

		CoreWorkflowSignalPullRequest.String():  CoreWorkflowSignalPullRequest,
		CoreWorkflowSignalLockAcquired.String(): CoreWorkflowSignalLockAcquired,
		CoreWorkflowSignalRequestLock.String():  CoreWorkflowSignalRequestLock,
		CoreWorkflowSignalReleaseLock.String():  CoreWorkflowSignalReleaseLock,
	}
)

/*
 * Methods for WorkflowSignal.
 */

func (w WorkflowSignal) String() string               { return string(w) }
func (w WorkflowSignal) MarshalJSON() ([]byte, error) { return json.Marshal(w.String()) }
func (w *WorkflowSignal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := WorkflowSignalMap[s]
	if !ok {
		return ErrInvalidRolloutState
	}

	*w = val

	return nil
}
