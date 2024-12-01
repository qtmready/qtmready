package defs

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/db/entities"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// signals.
const (
	SignalPush   queues.Signal = "push"   // signals a push event.
	SignalBranch queues.Signal = "branch" // signals a branch event.
	SignalTag    queues.Signal = "tag"    // signals a tag event.
	SignalPR     queues.Signal = "pr"     // signals a pull request event.
	SignalRebase queues.Signal = "rebase" // signals a rebase event.
)

const (
	QueryRepoForEventParent queues.Query = "event_parent"
)

type (
	ClonePayload struct {
		Repo   *entities.Repo    `json:"repo"`
		Hook   eventsv1.RepoHook `json:"hook"`
		Branch string            `json:"branch"`
		Path   string            `json:"path"`
		SHA    string            `json:"sha"`
	}

	DiffPayload struct {
		Path string `json:"path"`
		Base string `json:"base"`
		SHA  string `json:"sha"`
	}

	DiffFiles struct {
		Added      []string `json:"added"`
		Deleted    []string `json:"deleted"`
		Modified   []string `json:"modified"`
		Renamed    []string `json:"renamed"`
		Copied     []string `json:"copied"`
		TypeChange []string `json:"typechange"`
		Unreadable []string `json:"unreadable"`
		Ignored    []string `json:"ignored"`
		Untracked  []string `json:"untracked"`
		Conflicted []string `json:"conflicted"`
	}

	DiffLines struct {
		Added   int `json:"added"`
		Removed int `json:"removed"`
	}

	DiffResult struct {
		Files DiffFiles `json:"files"`
		Lines DiffLines `json:"lines"`
	}

	SignalBranchPayload struct {
		Signal queues.Signal  `json:"signal"`
		Repo   *entities.Repo `json:"repo"`
		Branch string         `json:"branch"`
	}

	SignalTrunkPayload struct {
		Signal queues.Signal  `json:"signal"`
		Repo   *entities.Repo `json:"repo"`
	}

	SignalQueuePayload struct{}
)

// Sum returns the sum of added and removed lines.
func (d *DiffLines) Sum() int {
	return d.Added + d.Removed
}
