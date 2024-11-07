package events

import (
	"time"

	"go.breu.io/quantm/internal/db"
)

// -- helpers --

type (
	// Commit represents a git commit.
	Commit struct {
		SHA       string    `json:"sha"`       // SHA is the SHA of the commit.
		Message   string    `json:"message"`   // Message is the commit message.
		URL       string    `json:"url"`       // URL is the URL of the commit.
		Added     []string  `json:"added"`     // Added is a list of files added in the commit.
		Removed   []string  `json:"removed"`   // Removed is a list of files removed in the commit.
		Modified  []string  `json:"modified"`  // Modified is a list of files modified in the commit.
		Author    string    `json:"author"`    // Author is the author of the commit.
		Committer string    `json:"committer"` // Committer is the committer of the commit.
		Timestamp time.Time `json:"timestamp"` // Timestamp is the timestamp of the commit.
	}

	Commits []Commit
)

// -- payloads --

type (
	// BranchOrTag represents a git branch or tag.
	BranchOrTag struct {
		Ref  string `json:"ref"`  // Ref is the name of the branch or tag.
		Kind string `json:"kind"` // Kind is the kind of the ref (branch or tag).
	}

	// Push represents a git push.
	Push struct {
		Ref        string    `json:"ref"`        // Ref is the ref that was pushed to.
		Before     string    `json:"before"`     // Before is the SHA of the commit before the push.
		After      string    `json:"after"`      // After is the SHA of the commit after the push.
		Repository string    `json:"repository"` // Repository is the repository that was pushed to.
		SenderID   db.Int64  `json:"sender_id"`  // SenderID is the id of the user who pushed the changes.
		Commits    Commits   `json:"commits"`    // Commits is a list of commits that were pushed.
		Timestamp  time.Time `json:"timestamp"`  // Timestamp is the timestamp of the push.
	}

	// EventPayload represents all available event payloads.
	EventPayload interface {
		BranchOrTag | Push // actions
	}
)
