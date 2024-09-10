package shared

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared/queue"
)

// workflow related shared types and contants.
type (
	Int64 int64
	// WorkflowSignal is a type alias to define the name of the workflow signal.
	WorkflowSignal string

	Duration struct {
		time.Duration
	}

	// PullRequestSignal is the sent to PR workflows to trigger a deployment.
	PullRequestSignal struct {
		RepoID           gocql.UUID
		SenderWorkflowID string
		TriggerID        Int64
		Image            string
		Digest           string
		ImageRegistry    string //TODO: move registry enum generation to shared
	}

	FutureHandler    func(workflow.Future)               // FutureHandler is the signature of the future handler for temporal.
	ChannelHandler   func(workflow.ReceiveChannel, bool) // ChannelHandler is the signature of the channel handler for temporal.
	CoroutineHandler func(workflow.Context)              // CoroutineHandler is the signature of the coroutine handler for temporal.

	WorkflowOption = queue.WorkflowOptions

	CreateChangesetSignal struct {
		RepoTableID gocql.UUID
		RepoID      string
		CommitID    string
	}

	PushEventSignal struct {
		RefBranch      string
		RepoProvider   string
		RepoID         Int64
		RepoName       string
		RepoOwner      string
		DefaultBranch  string
		InstallationID Int64
	}

	MergeQueueSignal struct {
		PullRequestID  Int64
		InstallationID Int64
		RepoOwner      string
		RepoName       string
		Branch         string
		RepoProvider   string
	}
)

var (
	WithWorkflowParent    = queue.WithWorkflowParent    // Sets the parent workflow ID.
	WithWorkflowBlock     = queue.WithWorkflowBlock     // Sets the block name for the workflow ID.
	WithWorkflowBlockID   = queue.WithWorkflowBlockID   // Sets the block value for the workflow ID.
	WithWorkflowElement   = queue.WithWorkflowElement   // Sets the element name for the workflow ID.
	WithWorkflowElementID = queue.WithWorkflowElementID // Sets the element value for the workflow ID.
	WithWorkflowMod       = queue.WithWorkflowMod       // Sets the modifier for the workflow ID.
	WithWorkflowModID     = queue.WithWorkflowModID     // Sets the modifier value for the workflow ID.
	WithWorkflowProp      = queue.WithWorkflowProp      // Sets the property for the workflow ID.
	NewWorkflowOptions    = queue.NewWorkflowOptions    // Creates a new workflow ID. (see queue.NewWorkflowOptions)
)

// queue definitions.
const (
	CoreQueue      queue.Name = "core"      // core queue
	ProvidersQueue queue.Name = "providers" // messaging related to providers
	MutexQueue     queue.Name = "mutex"     // mutex workflow queue
	WebSocketQueue queue.Name = "websocket" // websocket workflow queue
)

// workflow signal definitions.
const (
	WorkflowSignalDeploymentStarted WorkflowSignal = "deployment_trigger"
	WorkflowSignalCreateChangeset   WorkflowSignal = "create_changeset"
	WorkflowPushEvent               WorkflowSignal = "push_event_triggered"
	MergeQueueStarted               WorkflowSignal = "merge_queue_trigger"
	MergeTriggered                  WorkflowSignal = "merge_trigger"
)

/*
 * Methods for WorkflowSignal.
 */
func (w WorkflowSignal) String() string { return string(w) }

// MarshalJSON implements the json.Marshaler interface for WorkflowSignal.
func (w WorkflowSignal) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(w))
}

// UnmarshalJSON implements the json.Unmarshaler interface for WorkflowSignal.
func (w *WorkflowSignal) UnmarshalJSON(data []byte) error {
	var signal string
	if err := json.Unmarshal(data, &signal); err != nil {
		return err
	}

	*w = WorkflowSignal(signal)

	return nil
}

type (
	// EchoValidator is a wrapper for the instantiated validator.
	EchoValidator struct {
		Validator *validator.Validate
	}
)

func (ev *EchoValidator) Validate(i any) error {
	return ev.Validator.Struct(i)
}

// String returns the string representation of the Int64 value.
func (i Int64) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Int64 returns the int64 value of the Int64 .
func (i Int64) Int64() int64 {
	return int64(i)
}

// MarshalJSON implements the json.Marshaler interface.
func (i Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

// UnmarshalJSON implements the json.Unmarshaler .
func (i *Int64) UnmarshalJSON(data []byte) error {
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}

func (i Int64) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(info, i.Int64())
}

func (i *Int64) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	var v int64
	if err := gocql.Unmarshal(info, data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}

// UnmarshalCQL unmarshals a Cassandra duration (stored as text or varchar) into Duration.
// NOTE - saving duration as a string in database.
func (d *Duration) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	if info.Type() != gocql.TypeVarchar && info.Type() != gocql.TypeText {
		return fmt.Errorf("expected varchar or text type, got %v", info.Type())
	}

	if len(data) == 0 {
		return nil // Handle null values if needed
	}

	s := string(data)

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("error parsing duration %s: %w", s, err)
	}

	d.Duration = duration

	return nil
}

// MarshalCQL marshals Duration into a Cassandra duration.
func (d Duration) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(d.String()), nil
}

// MarshalJSON marshals Duration to JSON string.
func (d Duration) MarshalJSON() ([]byte, error) {
	if d.Duration == 0 {
		return []byte("null"), nil // Handle null case for JSON
	}

	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

// UnmarshalJSON unmarshals JSON string into Duration.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "null" {
		return nil // Handle null case for JSON
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("error parsing duration %s: %w", s, err)
	}

	d.Duration = duration

	return nil
}

// String returns the string representation of Duration.
func (d Duration) String() string {
	return d.Duration.String()
}

// NewDuration creates a new Duration from a time.Duration.
func NewDuration(d time.Duration) Duration {
	return Duration{d}
}
