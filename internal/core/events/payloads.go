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

	// PullRequest represents a pull request.
	PullRequest struct {
		Number         db.Int64  `json:"number"`                     // Number is the pull request number.
		Title          string    `json:"title"`                      // Title is the pull request title.
		Body           string    `json:"body"`                       // Body is the pull request body.
		State          string    `json:"state"`                      // State is the pull request state.
		MergeCommitSHA *string   `json:"merge_commit_sha,omitempty"` // MergeCommitSHA is the SHA of the merge commit.
		AuthorID       db.Int64  `json:"author_id"`                  // AuthorID is the author_id of the pull request.
		HeadBranch     string    `json:"head_branch"`                // HeadBranch is the head branch of the pull request.
		BaseBranch     string    `json:"base_branch"`                // BaseBranch is the base branch of the pull request.
		Timestamp      time.Time `json:"timestamp"`                  // Timestamp is the timestamp when the pull request was created.
	}

	// PullRequestReview represents a pull request review.
	PullRequestReview struct {
		ID                db.Int64  `json:"id"`                  // ID is the pull request review ID.
		PullRequestNumber db.Int64  `json:"pull_request_number"` // PullRequestNumber is the pull request number.
		Branch            string    `json:"branch"`              // Branch is the branch the review belongs to.
		State             string    `json:"state"`               // State is the pull request review state.
		AuthorID          db.Int64  `json:"author_id"`           // AuthorID is the author of the review.
		Timestamp         time.Time `json:"submitted_at"`        // SubmittedAt is the timestamp when the review was submitted.
	}

	// PullRequestLabel represents a pull request label.
	PullRequestLabel struct {
		Name              string    `json:"name"`                // Name is the text of the label e.g. "ready", "fix" etc.
		PullRequestNumber db.Int64  `json:"pull_request_number"` // PullRequestNumber is the pull request number.
		Branch            string    `json:"branch"`              // Branch is the branch the label belongs to.
		Timestamp         time.Time `json:"timestamp"`           // Timestamp is the timestamp of the label.
	}

	// PullRequestComment represents a pull request comment.
	PullRequestComment struct {
		ID                db.Int64  `json:"id"`                    // ID is the pull request review comment ID.
		PullRequestNumber db.Int64  `json:"pull_request_number"`   // PullRequestNumber is the pull request number.
		Branch            string    `json:"branch"`                // Branch is the branch the comment belongs to.
		ReviewID          db.Int64  `json:"review_id"`             // ReviewID is the ID of the pull request review the comment belongs.
		InReplyTo         *db.Int64 `json:"in_reply_to,omitempty"` // InReplyTo is the ID of the parent comment.
		CommitSHA         string    `json:"commit_sha"`            // CommitSHA is the SHA of the commit associated with the comment.
		Path              string    `json:"path"`                  // Path is the path to the file where the comment was made.
		Position          db.Int64  `json:"position"`              // Position is the line number where the comment was made.
		AuthorID          db.Int64  `json:"author_id"`             // AuthorID is the author_id of the comment.
		Timestamp         time.Time `json:"timestamp"`             // Timestamp is the timestamp of the comment.
	}

	// PullRequestThread represents a pull request thread.
	PullRequestThread struct {
		ID                db.Int64   `json:"id"`                  // ID is the pull request thread ID.
		PullRequestNumber db.Int64   `json:"pull_request_number"` // PullRequestNumber is the pull request number.
		CommentIDs        []db.Int64 `json:"comment_ids"`         // CommentIDs is the list of comment IDs associated with the thread.
		Timestamp         time.Time  `json:"timestamp"`           // Timestamp is the timestamp of the thread.
	}

	// TODO - set the event payload.
	CommitDiff struct {
		Added     db.Int64 `json:"added"`     // Number of lines added in the commit.
		Removed   db.Int64 `json:"removed"`   // Number of lines removed in the commit.
		Threshold db.Int64 `json:"threshold"` // Set threshold for PR.
		Delta     db.Int64 `json:"delta"`     // Net change in lines (added - removed).
	}

	// MergeConflict represents a git merge conflict.
	MergeConflict struct {
		HeadBranch string    `json:"head_branch"` // HeadBranch is the name of the head branch.
		HeadCommit Commit    `json:"head_commit"` // HeadCommit is the last commit on the head branch before rebasing.
		BaseBranch string    `json:"base_branch"` // BaseBranch is the name of the base branch.
		BaseCommit Commit    `json:"base_commit"` // BaseCommit is the last commit on the base branch before rebasing.
		Files      []string  `json:"files"`       // Files is the list of files with conflicts.
		Timestamp  time.Time `json:"timestamp"`   // Timestamp is the timestamp of the merge conflict.
	}

	LinesExceed struct {
		Branch    string     `json:"branch"`    // Branch is the name of the head or feature branch.
		Commit    Commit     `json:"commit"`    // Commit is the last commit on the head branch.
		Diff      CommitDiff `json:"diff"`      // LineStats contains details about lines added, removed, and the delta.
		Timestamp time.Time  `json:"timestamp"` // Timestamp is the timestamp of the merge conflict.
	}

	RebaseRequest struct {
		Ref        string    `json:"ref"`         // Ref is the ref that was pushed to.
		Before     string    `json:"before"`      // Before is the SHA of the commit before the push.
		After      string    `json:"after"`       // After is the SHA of the commit after the push.
		HeadBranch string    `json:"head_branch"` // HeadBranch is the name of the head branch.
		HeadCommit Commit    `json:"head_commit"` // HeadCommit is the last commit on the head branch before rebasing.
		BaseBranch string    `json:"base_branch"` // BaseBranch is the name of the base branch.
		BaseCommit Commit    `json:"base_commit"` // BaseCommit is the last commit on the base branch before rebasing.
		Timestamp  time.Time `json:"timestamp"`   // Timestamp is the timestamp of the merge conflict.
	}

	// EventPayload represents all available event payloads.
	EventPayload interface {
		BranchOrTag | Push | // actions
			PullRequest | PullRequestReview | PullRequestLabel | PullRequestComment | PullRequestThread | // activities
			RebaseRequest | // requests and responses
			MergeConflict | LinesExceed // errors or results
	}
)

// ToEvent converts a BranchOrTag struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "branches" or "tags". The action parameter specifies the type of the action, such as "created", "deleted", or
// "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "branch" or "tag". The Timestamp
// field of the EventContext is set to the current time.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	branch := &BranchOrTag{
//	  Ref: "main",
//	}
//	event := branch.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	  EventScopeBranch,
//	)
func (bt *BranchOrTag) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction, scope EventScope,
) *Event[BranchOrTag, RepoProvider] {
	event := &Event[BranchOrTag, RepoProvider]{
		Version: EventVersionDefault,
		ID:      db.MustUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     scope,
			Action:    action,
			Timestamp: time.Now(),
		},
		Subject: subject,
		Payload: *bt,
	}

	return event
}

// ToEvent converts a Push struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pushes". The action parameter specifies the type of the action, such as "created" or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "push". The Timestamp field of
// the EventContext is set to the Timestamp field of the Push struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	push := &Push{
//	  Ref: "main",
//	  Before: "old_sha",
//	  After: "new_sha",
//	  Repository: "example/repo",
//	  Pusher: "user",
//	  Commits: []Commit{},
//	  Timestamp: time.Now(),
//	}
//	event := push.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	)
func (p *Push) ToEvent() *Event[Push, RepoProvider] {
	event := &Event[Push, RepoProvider]{
		Version: EventVersionDefault,
		Subject: EventSubject{
			Name: "repos",
		},
		Payload: *p,
	}

	return event
}

// ToEvent converts a PullRequest struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_requests". The action parameter specifies the type of the action, such as "updated", "closed", or "merged".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request". The Timestamp
// field of the EventContext is set to the UpdatedAt field of the PullRequest struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	pr := &PullRequest{
//	  Number: 1,
//	  Title:  "Test Pull Request",
//	  Body:   "This is a test pull request",
//	  State:  "open",
//	  UpdatedAt: time.Now(),
//	}
//	event := pr.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionUpdated,
//	)
func (pr *PullRequest) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequest, RepoProvider] {
	event := &Event[PullRequest, RepoProvider]{
		Version: EventVersionDefault,
		ID:      db.MustUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequest,
			Action:    action,
			Timestamp: pr.Timestamp,
		},
		Subject: subject,
		Payload: *pr,
	}

	return event
}

// ToEvent converts a PullRequestReview struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_reviews". The action parameter specifies the type of the action, such as "submitted" or "edited".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_review". The
// Timestamp field of the EventContext is set to the SubmittedAt field of the PullRequestReview struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	review := &PullRequestReview{
//	  ID: 1,
//	  Body: "This is a review",
//	  State: "approved",
//	  SubmittedAt: time.Now(),
//	  Author: "user",
//	  PullRequestNumber: 1,
//	}
//	event := review.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionSubmitted,
//	)
func (prr *PullRequestReview) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestReview, RepoProvider] {
	event := &Event[PullRequestReview, RepoProvider]{
		Version: EventVersionDefault,
		ID:      db.MustUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestReview,
			Action:    action,
			Timestamp: prr.Timestamp,
		},
		Subject: subject,
		Payload: *prr,
	}

	return event
}

// ToEvent converts a PullRequestLabel struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_labels". The action parameter specifies the type of the action, such as "added" or "removed".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_label". The
// Timestamp field of the EventContext is set to the UpdatedAt field of the PullRequestLabel struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	label := &PullRequestLabel{
//	  Name: "bug",
//	  Color: "red",
//	  Description: "This is a bug label",
//	  CreatedAt: time.Now(),
//	  UpdatedAt: time.Now(),
//	  PullRequestNumber: 1,
//	}
//	event := label.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionAdded,
//	)
func (prl *PullRequestLabel) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestLabel, RepoProvider] {
	event := &Event[PullRequestLabel, RepoProvider]{
		Version: EventVersionDefault,
		ID:      db.MustUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestLabel,
			Action:    action,
			Timestamp: prl.Timestamp,
		},
		Subject: subject,
		Payload: *prl,
	}

	return event
}

// ToEvent converts a PullRequestComment struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_comments". The action parameter specifies the type of the action, such as "created" or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_comment". The
// Timestamp field of the EventContext is set to the Timestamp field of the PullRequestComment struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	comment := &PullRequestComment{
//	  ID: 1,
//	  Path: "path/to/file.go",
//	  Position: 15,
//	  Author: "user",
//	  PullRequestNumber: 1,
//	  ReviewID: 5,
//	  CommitSHA: "abcdef1234567890",
//	  InReplyTo: nil,
//	  Timestamp: time.Now(),
//	}
//	event := comment.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	)
func (prc *PullRequestComment) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestComment, RepoProvider] {
	event := &Event[PullRequestComment, RepoProvider]{
		Version: EventVersionDefault,
		ID:      db.MustUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestComment,
			Action:    action,
			Timestamp: prc.Timestamp,
		},
		Subject: subject,
		Payload: *prc,
	}

	return event
}

// ToEvent converts a PullRequestThread struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_threads". The action parameter specifies the type of the action, such as "created" or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_thread". The
// Timestamp field of the EventContext is set to the UpdatedAt field of the PullRequestThread struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	thread := &PullRequestThread{
//	  ID: 1,
//	  Title: "Question about implementation",
//	  Comments: []Comment{},
//	  CreatedAt: time.Now(),
//	  UpdatedAt: time.Now(),
//	  Path: "path/to/file.go",
//	  Position: 15,
//	}
//	event := thread.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	)
func (prt *PullRequestThread) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestThread, RepoProvider] {
	event := &Event[PullRequestThread, RepoProvider]{
		Version: EventVersionDefault,
		ID:      db.MustUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestThread,
			Action:    action,
			Timestamp: prt.Timestamp,
		},
		Subject: subject,
		Payload: *prt,
	}

	return event
}
