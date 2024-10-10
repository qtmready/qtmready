package defs

import (
	"github.com/gocql/gocql"

	"go.breu.io/quantm/internal/core/events"
	"go.breu.io/quantm/internal/db"
)

// -- Webhook Events --

type (
	AppAuthorizationEvent struct {
		Action string `json:"action"`
		Sender User   `json:"sender"`
	}

	InstallationEvent struct {
		Action       string              `json:"action"`
		Installation InstallationPayload `json:"installation"`
		Repositories []PartialRepository `json:"repositories"`
		Sender       User                `json:"sender"`
	}

	PushEvent struct {
		Ref          string         `json:"ref"`
		Before       string         `json:"before"`
		After        string         `json:"after"`
		Created      bool           `json:"created"`
		Deleted      bool           `json:"deleted"`
		Forced       bool           `json:"forced"`
		BaseRef      *string        `json:"base_ref"`
		Compare      string         `json:"compare"`
		Commits      []Commit       `json:"commits"`
		HeadCommit   Commit         `json:"head_commit"`
		Repository   Repository     `json:"repository"`
		Pusher       Pusher         `json:"pusher"`
		Sender       User           `json:"sender"`
		Installation InstallationID `json:"installation"`
	}

	CreateOrDeleteEvent struct {
		Ref          string         `json:"ref"`
		RefType      string         `json:"ref_type"`
		MasterBranch *string        `json:"master_branch"` // NOTE: This is only present in the create event.
		Description  *string        `json:"description"`   // NOTE: This is only present in the create event.
		PusherType   string         `json:"pusher_type"`
		Repository   Repository     `json:"repository"`
		Organization Organization   `json:"organization"`
		Sender       User           `json:"sender"`
		Installation InstallationID `json:"installation"`
		IsCreated    bool           `json:"is_created"`
	}

	GithubWorkflowRunEvent struct {
		Action       string             `json:"action"`
		Repository   RepositoryPR       `json:"repository"`
		Sender       User               `json:"sender"`
		Installation InstallationID     `json:"installation"`
		WR           WorkflowRunPayload `json:"workflow_run"`
		Workflow     WorkflowPayload    `json:"workflow"`
	}

	PullRequestEvent struct {
		Action       string         `json:"action"`
		Number       db.Int64       `json:"number"`
		PullRequest  PullRequest    `json:"pull_request"`
		Repository   RepositoryPR   `json:"repository"`
		Organization *Organization  `json:"organization"`
		Installation InstallationID `json:"installation"`
		Sender       User           `json:"sender"`
		Label        *Label         `json:"label"`
	}

	InstallationRepositoriesEvent struct {
		Action              string              `json:"action"`
		Installation        InstallationPayload `json:"installation"`
		RepositorySelection string              `json:"repository_selection"`
		RepositoriesAdded   []PartialRepository `json:"repositories_added"`
		RepositoriesRemoved []PartialRepository `json:"repositories_removed"`
		Requester           *User               `json:"requester"`
		Sender              User                `json:"sender"`
	}

	CompleteInstallationSignal struct {
		InstallationID db.Int64    `json:"installation_id"`
		SetupAction    SetupAction `json:"setup_action"`
		UserID         gocql.UUID  `json:"user_id"`
	}

	ArtifactReadySignal struct {
		Image    string
		Digest   string
		Registry string
	}

	GithubActionResult struct {
		Branch         string
		InstallationID db.Int64
		RepoID         string
		RepoName       string
		RepoOwner      string
	}

	PullRequestReviewEvent struct {
		Action       string             `json:"action"`
		Number       db.Int64           `json:"number"`
		Installation InstallationID     `json:"installation"`
		Review       *PullRequestReview `json:"review"`
		PullRequest  PullRequest        `json:"pull_request"`
		Repository   RepositoryPR       `json:"repository"`
		Sender       *User              `json:"sender"`
	}

	PullRequestReviewCommentEvent struct {
		Action       string              `json:"action"`
		Number       db.Int64            `json:"number"`
		Installation InstallationID      `json:"installation"`
		Comment      *PullRequestComment `json:"comment"`
		PullRequest  PullRequest         `json:"pull_request"`
		Repository   RepositoryPR        `json:"repository"`
		Sender       *User               `json:"sender"`
	}
)

func (evt *CreateOrDeleteEvent) Payload() *events.BranchOrTag {
	return &events.BranchOrTag{
		Ref:  evt.Ref,
		Kind: evt.RefType,
	}
}
