package github

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/shared"
)

type (
	Timestamp time.Time // Timestamp is hack around github's funky use of time
)

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case float64:
		*t = Timestamp(time.Unix(int64(v), 0))
	case string:
		if strings.HasSuffix(v, "Z") {
			t_, err := time.Parse("2006-01-02T15:04:05Z", v)
			if err != nil {
				return err
			}

			*t = Timestamp(t_)
		} else {
			t_, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}

			*t = Timestamp(t_)
		}
	}

	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	t_ := time.Time(t)
	return json.Marshal(t_.Format(time.RFC3339))
}

func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

// Webhook events & Workflow singals.
type (
	WebhookEvent         string                               // WebhookEvent defines the event type.
	WebhookEventHandler  func(ctx echo.Context) error         // EventHandler is the signature of the event handler function.
	WebhookEventHandlers map[WebhookEvent]WebhookEventHandler // EventHandlers maps event types to their respective event handlers.

	RepoEvent interface {
		RepoID() shared.Int64
		RepoName() string
		InstallationID() shared.Int64
		SenderID() string
	}
)

type (
	CreateMembershipsPayload struct {
		UserID        gocql.UUID   `json:"user_id"`
		TeamID        gocql.UUID   `json:"team_id"`
		IsAdmin       bool         `json:"is_admin"`
		GithubOrgName string       `json:"github_org_name"`
		GithubOrgID   shared.Int64 `json:"github_org_id"`
		GithubUserID  shared.Int64 `json:"github_user_id"`
	}

	PostInstallPayload struct {
		InstallationID    shared.Int64 `json:"installation_id"`
		InstallationLogin string       `json:"installation_login"`
	}

	SyncReposFromGithubPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		Owner          string       `json:"owner"`
		TeamID         gocql.UUID   `json:"team_id"`
	}

	SyncOrgUsersFromGithubPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		GithubOrgName  string       `json:"github_org_name"`
		GithubOrgID    shared.Int64 `json:"github_org_id"`
	}
)

// Payloads for internal events & signals.
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
		HeadCommit   HeadCommit     `json:"head_commit"`
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
		Number       shared.Int64   `json:"number"`
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
		InstallationID shared.Int64 `json:"installation_id"`
		SetupAction    SetupAction  `json:"setup_action"`
		UserID         gocql.UUID   `json:"user_id"`
	}

	ArtifactReadySignal struct {
		Image    string
		Digest   string
		Registry string
	}

	GithubActionResult struct {
		Branch         string
		InstallationID shared.Int64
		RepoID         string
		RepoName       string
		RepoOwner      string
	}

	PullRequestReviewEvent struct {
		Action       string             `json:"action"`
		Number       shared.Int64       `json:"number"`
		Installation InstallationID     `json:"installation"`
		Review       *PullRequestReview `json:"review"`
		PullRequest  PullRequest        `json:"pull_request"`
		Repository   RepositoryPR       `json:"repository"`
		Sender       *User              `json:"sender"`
	}

	PullRequestReviewCommentEvent struct {
		Action       string              `json:"action"`
		Number       shared.Int64        `json:"number"`
		Installation InstallationID      `json:"installation"`
		Comment      *PullRequestComment `json:"comment"`
		PullRequest  PullRequest         `json:"pull_request"`
		Repository   RepositoryPR        `json:"repository"`
		Sender       *User               `json:"sender"`
	}
)

// Webhook event types. We get this from the header `X-Github-Event`.
// For payload information, see https://developer.github.com/webhooks/event-payloads/.
const (
	WebhookEventAppAuthorization                    WebhookEvent = "github_app_authorization" // nolint:gosec
	WebhookEventCheckRun                            WebhookEvent = "check_run"
	WebhookEventCheckSuite                          WebhookEvent = "check_suite"
	WebhookEventCommitComment                       WebhookEvent = "commit_comment"
	WebhookEventCreate                              WebhookEvent = "create"
	WebhookEventDelete                              WebhookEvent = "delete"
	WebhookEventDeployKey                           WebhookEvent = "deploy_key"
	WebhookEventDeployment                          WebhookEvent = "deployment"
	WebhookEventDeploymentStatus                    WebhookEvent = "deployment_status"
	WebhookEventFork                                WebhookEvent = "fork"
	WebhookEventGollum                              WebhookEvent = "gollum"
	WebhookEventInstallation                        WebhookEvent = "installation"
	WebhookEventInstallationRepositories            WebhookEvent = "installation_repositories"
	WebhookEventIntegrationInstallation             WebhookEvent = "integration_installation"
	WebhookEventIntegrationInstallationRepositories WebhookEvent = "integration_installation_repositories"
	WebhookEventIssueComment                        WebhookEvent = "issue_comment"
	WebhookEventIssues                              WebhookEvent = "issues"
	WebhookEventLabel                               WebhookEvent = "label"
	WebhookEventMember                              WebhookEvent = "member"
	WebhookEventMembership                          WebhookEvent = "membership"
	WebhookEventMilestone                           WebhookEvent = "milestone"
	WebhookEventMeta                                WebhookEvent = "meta"
	WebhookEventOrganization                        WebhookEvent = "organization"
	WebhookEventOrgBlock                            WebhookEvent = "org_block"
	WebhookEventPageBuild                           WebhookEvent = "page_build"
	WebhookEventPing                                WebhookEvent = "ping"
	WebhookEventProjectCard                         WebhookEvent = "project_card"
	WebhookEventProjectColumn                       WebhookEvent = "project_column"
	WebhookEventProject                             WebhookEvent = "project"
	WebhookEventPublic                              WebhookEvent = "public"
	WebhookEventPullRequest                         WebhookEvent = "pull_request"
	WebhookEventPullRequestReview                   WebhookEvent = "pull_request_review"
	WebhookEventPullRequestReviewComment            WebhookEvent = "pull_request_review_comment"
	WebhookEventPush                                WebhookEvent = "push"
	WebhookEventRelease                             WebhookEvent = "release"
	WebhookEventRepository                          WebhookEvent = "repository"
	WebhookEventRepositoryVulnerabilityAlert        WebhookEvent = "repository_vulnerability_alert"
	WebhookEventSecurityAdvisory                    WebhookEvent = "security_advisory"
	WebhookEventStatus                              WebhookEvent = "status"
	WebhookEventTeam                                WebhookEvent = "team"
	WebhookEventTeamAdd                             WebhookEvent = "team_add"
	WebhookEventWatch                               WebhookEvent = "watch"
	WebhookEventWorkflowDispatch                    WebhookEvent = "workflow_dispatch"
	WebhookEventWorkflowJob                         WebhookEvent = "workflow_job"
	WebhookEventWorkflowRun                         WebhookEvent = "workflow_run"
)

const (
	NoCommit = "0000000000000000000000000000000000000000"
)

func (e WebhookEvent) String() string { return string(e) }

// Workflow signal types.
const (
	WorkflowSignalInstallationEvent    shared.WorkflowSignal = "installation_event"
	WorkflowSignalCompleteInstallation shared.WorkflowSignal = "complete_installation"
	WorkflowSignalPullRequestProcessed shared.WorkflowSignal = "pull_request_processed"
	WorkflowSignalArtifactReady        shared.WorkflowSignal = "artifact_ready"
	WorkflowSignalActionResult         shared.WorkflowSignal = "action_result"
	WorkflowSignalPullRequestLabeled   shared.WorkflowSignal = "pull_request_labeled"
	WorkflowSignalPushEvent            shared.WorkflowSignal = "push_event"
)

type (
	RepoEventState struct {
		CoreRepo *defs.Repo     `json:"core_repo"`
		Repo     *Repo          `json:"repo"`
		User     *auth.TeamUser `json:"user"`
	}
)

func (p *PushEvent) RepoID() shared.Int64 {
	return p.Repository.ID
}

func (p *PushEvent) InstallationID() shared.Int64 {
	return p.Installation.ID
}

func (p *PushEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PushEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *PullRequestEvent) RepoID() shared.Int64 {
	return p.Repository.ID
}

func (p *PullRequestEvent) InstallationID() shared.Int64 {
	return p.Installation.ID
}

func (p *PullRequestEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PullRequestEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *CreateOrDeleteEvent) RepoID() shared.Int64 {
	return p.Repository.ID
}

func (p *CreateOrDeleteEvent) InstallationID() shared.Int64 {
	return p.Installation.ID
}

func (p *CreateOrDeleteEvent) RepoName() string {
	return p.Repository.Name
}

func (p *CreateOrDeleteEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *PullRequestReviewEvent) RepoID() shared.Int64 {
	return p.Repository.ID
}

func (p *PullRequestReviewEvent) InstallationID() shared.Int64 {
	return p.Installation.ID
}

func (p *PullRequestReviewEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PullRequestReviewEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *PullRequestReviewCommentEvent) RepoID() shared.Int64 {
	return p.Repository.ID
}

func (p *PullRequestReviewCommentEvent) InstallationID() shared.Int64 {
	return p.Installation.ID
}

func (p *PullRequestReviewCommentEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PullRequestReviewCommentEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *GithubWorkflowRunEvent) RepoID() shared.Int64 {
	return p.Repository.ID
}

func (p *GithubWorkflowRunEvent) InstallationID() shared.Int64 {
	return p.Installation.ID
}

func (p *GithubWorkflowRunEvent) RepoName() string {
	return p.Repository.Name
}

func (p *GithubWorkflowRunEvent) SenderID() string {
	return p.Sender.ID.String()
}
