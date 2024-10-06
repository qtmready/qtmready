// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package github

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
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
		RepoID() db.Int64
		RepoName() string
		InstallationID() db.Int64
		SenderID() string
	}
)

type (
	CreateMembershipsPayload struct {
		UserID        gocql.UUID `json:"user_id"`
		TeamID        gocql.UUID `json:"team_id"`
		IsAdmin       bool       `json:"is_admin"`
		GithubOrgName string     `json:"github_org_name"`
		GithubOrgID   db.Int64   `json:"github_org_id"`
		GithubUserID  db.Int64   `json:"github_user_id"`
	}

	PostInstallPayload struct {
		InstallationID    db.Int64 `json:"installation_id"`
		InstallationLogin string   `json:"installation_login"`
	}

	SyncReposFromGithubPayload struct {
		InstallationID db.Int64   `json:"installation_id"`
		Owner          string     `json:"owner"`
		TeamID         gocql.UUID `json:"team_id"`
	}

	SyncOrgUsersFromGithubPayload struct {
		InstallationID db.Int64 `json:"installation_id"`
		GithubOrgName  string   `json:"github_org_name"`
		GithubOrgID    db.Int64 `json:"github_org_id"`
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
	WorkflowSignalInstallationEvent    queues.Signal = "installation_event"
	WorkflowSignalCompleteInstallation queues.Signal = "complete_installation"
	WorkflowSignalPullRequestProcessed queues.Signal = "pull_request_processed"
	WorkflowSignalArtifactReady        queues.Signal = "artifact_ready"
	WorkflowSignalActionResult         queues.Signal = "action_result"
	WorkflowSignalPullRequestLabeled   queues.Signal = "pull_request_labeled"
	WorkflowSignalPushEvent            queues.Signal = "push_event"
)

type (
	RepoEventMetadataQuery struct {
		RepoID         db.Int64 `json:"repo_id"`
		RepoName       string   `json:"repo_name"`
		InstallationID db.Int64 `json:"installation_id"`
		SenderID       string   `json:"sender_id"`
	}
	RepoEventMetadata struct {
		CoreRepo *defs.Repo     `json:"core_repo"`
		Repo     *Repo          `json:"repo"`
		User     *auth.TeamUser `json:"user"`
	}
)

func (p *PushEvent) RepoID() db.Int64 {
	return p.Repository.ID
}

func (p *PushEvent) InstallationID() db.Int64 {
	return p.Installation.ID
}

func (p *PushEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PushEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *PullRequestEvent) RepoID() db.Int64 {
	return p.Repository.ID
}

func (p *PullRequestEvent) InstallationID() db.Int64 {
	return p.Installation.ID
}

func (p *PullRequestEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PullRequestEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *CreateOrDeleteEvent) RepoID() db.Int64 {
	return p.Repository.ID
}

func (p *CreateOrDeleteEvent) InstallationID() db.Int64 {
	return p.Installation.ID
}

func (p *CreateOrDeleteEvent) RepoName() string {
	return p.Repository.Name
}

func (p *CreateOrDeleteEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *PullRequestReviewEvent) RepoID() db.Int64 {
	return p.Repository.ID
}

func (p *PullRequestReviewEvent) InstallationID() db.Int64 {
	return p.Installation.ID
}

func (p *PullRequestReviewEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PullRequestReviewEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *PullRequestReviewCommentEvent) RepoID() db.Int64 {
	return p.Repository.ID
}

func (p *PullRequestReviewCommentEvent) InstallationID() db.Int64 {
	return p.Installation.ID
}

func (p *PullRequestReviewCommentEvent) RepoName() string {
	return p.Repository.Name
}

func (p *PullRequestReviewCommentEvent) SenderID() string {
	return p.Sender.ID.String()
}

func (p *GithubWorkflowRunEvent) RepoID() db.Int64 {
	return p.Repository.ID
}

func (p *GithubWorkflowRunEvent) InstallationID() db.Int64 {
	return p.Installation.ID
}

func (p *GithubWorkflowRunEvent) RepoName() string {
	return p.Repository.Name
}

func (p *GithubWorkflowRunEvent) SenderID() string {
	return p.Sender.ID.String()
}

// prelude is a helper function to create a base event structure with common fields.
//
// It takes a `defs.Repo` pointer and returns a `gocql.UUID`, `defs.EventVersion`, `defs.EventContext`, and
// `defs.EventSubject` that can be used to construct a `defs.Event`.
func prelude(
	repo *defs.Repo,
) (gocql.UUID, defs.EventVersion, defs.EventContext[defs.RepoProvider], defs.EventSubject) {
	id, _ := db.NewUUID()
	version := defs.EventVersionDefault

	ctx := defs.EventContext[defs.RepoProvider]{
		Provider:  repo.Provider,
		Timestamp: time.Now(),
	}

	sub := defs.EventSubject{
		ID:     repo.ID,
		Name:   "repos",
		TeamID: repo.TeamID,
	}

	return id, version, ctx, sub
}

// payload converts the CreateOrDeleteEvent struct to the relevant EventPayload.
//
// It returns a `BranchOrTag` struct containing the `ref` and `default_branch` fields.
func (coe *CreateOrDeleteEvent) payload() defs.BranchOrTag {
	result := defs.BranchOrTag{
		Ref: coe.Ref,
	}

	if coe.MasterBranch != nil {
		result.DefaultBranch = *coe.MasterBranch
	}

	return result
}

// normalize converts the CreateOrDeleteEvent struct to an Event struct.
//
// It uses the provided Repo struct to extract relevant information for the EventContext and EventSubject.
// The action is set to either "created" or "deleted" based on the `IsCreated` flag.
func (coe *CreateOrDeleteEvent) normalize(repo *defs.Repo) *defs.Event[defs.BranchOrTag, defs.RepoProvider] {
	id, version, ctx, sub := prelude(repo)
	payload := coe.payload()

	if payload.DefaultBranch == "" {
		payload.DefaultBranch = repo.DefaultBranch
	}

	event := &defs.Event[defs.BranchOrTag, defs.RepoProvider]{
		ID:      id,
		Version: version,
		Context: ctx,
		Subject: sub,
		Payload: payload,
	}

	event.SetSource(coe.Repository.URL)
	event.SetActionCreated()
	event.SetScopeBranch()

	if !coe.IsCreated {
		event.SetActionDeleted()
	}

	if coe.RefType != "branch" {
		event.SetScopeTag()
	}

	return event
}

// payload converts the PushEvent struct to the relevant EventPayload.
//
// It returns a `defs.Push` struct containing the relevant information for a push event.
func (pe PushEvent) payload() defs.Push {
	commits := make(defs.Commits, len(pe.Commits))
	for i, c := range pe.Commits {
		commits[i] = c.normalize()
	}

	return defs.Push{
		Ref:        pe.Ref,
		Before:     pe.Before,
		After:      pe.After,
		Repository: pe.Repository.Name,
		SenderID:   pe.Sender.ID,
		Commits:    commits,
		Timestamp:  pe.HeadCommit.Timestamp.Time(),
	}
}

// normalize converts the PushEvent struct to an Event struct.
//
// It uses the provided Repo struct to extract relevant information for the EventContext and EventSubject.
func (pe PushEvent) normalize(repo *defs.Repo) *defs.Event[defs.Push, defs.RepoProvider] {
	id, version, ctx, sub := prelude(repo)
	event := &defs.Event[defs.Push, defs.RepoProvider]{
		ID:      id,
		Version: version,
		Context: ctx,
		Subject: sub,
		Payload: pe.payload(),
	}

	event.SetSource(pe.Repository.URL)
	event.SetScopePush()
	event.SetActionCreated()

	return event
}

// payload converts the PullRequestEvent struct to the relevant EventPayload.
//
// It returns a `defs.PullRequest` struct containing the relevant information for a pull request event.
func (pre PullRequestEvent) payload() defs.PullRequest {
	return defs.PullRequest{
		Number:         pre.Number,
		Title:          pre.PullRequest.Title,
		Body:           pre.PullRequest.Body,
		State:          pre.PullRequest.State,
		MergeCommitSHA: pre.PullRequest.MergeCommitSha,
		AuthorID:       pre.PullRequest.User.ID,
		HeadBranch:     pre.PullRequest.Head.Ref,
		BaseBranch:     pre.PullRequest.Base.Ref,
		Timestamp:      pre.PullRequest.UpdatedAt,
	}
}

// normalize converts the PullRequestEvent struct to an Event struct.
//
// It uses the provided Repo struct to extract relevant information for the EventContext and EventSubject.
// The action is set based on the `Action` field of the PullRequestEvent struct and whether MergedCommitSHA is set.
func (pre PullRequestEvent) normalize(repo *defs.Repo) *defs.Event[defs.PullRequest, defs.RepoProvider] {
	id, version, ctx, sub := prelude(repo)
	event := &defs.Event[defs.PullRequest, defs.RepoProvider]{
		ID:      id,
		Version: version,
		Context: ctx,
		Subject: sub,
		Payload: pre.payload(),
	}

	event.SetSource(pre.Repository.URL)
	event.SetScopePullRequest()

	switch pre.Action {
	case "opened":
		event.SetActionCreated()
	case "reopened":
		event.SetActionCreated()
	case "closed":
		event.SetActionClosed()
	case "edited": //nolint
		event.SetActionUpdated()
	case "assigned":
		event.SetActionUpdated()
	case "unassigned":
		event.SetActionUpdated()
	case "review_requested":
		event.SetActionUpdated()
	case "review_request_removed":
		event.SetActionUpdated()
	case "synchronized":
		event.SetActionUpdated()
	case "labeled":
		event.SetActionAdded()
		event.SetScopePullRequestLabel()
	case "unlabeled":
		event.SetActionDeleted()
		event.SetScopePullRequestLabel()
	default:
		return nil
	}

	// Determine "merged" action based on MergedCommitSHA
	if pre.PullRequest.MergeCommitSha != nil {
		event.SetActionMerged()
	}

	return event
}

func (pre PullRequestEvent) as_label(
	event *defs.Event[defs.PullRequest, defs.RepoProvider],
) *defs.Event[defs.PullRequestLabel, defs.RepoProvider] {
	if pre.Label == nil {
		return nil
	}

	label := defs.PullRequestLabel{
		Name:              pre.Label.Name,
		PullRequestNumber: pre.Number,
		Branch:            pre.PullRequest.Head.Ref,
		Timestamp:         pre.PullRequest.UpdatedAt,
	}

	return &defs.Event[defs.PullRequestLabel, defs.RepoProvider]{
		ID:      event.ID,
		Version: event.Version,
		Context: event.Context,
		Subject: event.Subject,
		Payload: label,
	}
}

// payload converts the PullRequestReviewEvent struct to the relevant EventPayload.
//
// It returns a `defs.PullRequestReview` struct containing the relevant information for a pull request review event.
func (pre PullRequestReviewEvent) payload() defs.PullRequestReview {
	return defs.PullRequestReview{
		ID:                pre.Review.ID,
		State:             pre.Review.State,
		AuthorID:          pre.Review.User.ID,
		PullRequestNumber: pre.Number,
		Timestamp:         pre.Review.SubmittedAt.Time(),
	}
}

// normalize converts the PullRequestReviewEvent struct to an Event struct.
//
// It uses the provided Repo struct to extract relevant information for the EventContext and EventSubject.
// The action is set based on the `Action` field of the PullRequestReviewEvent struct.
func (pre PullRequestReviewEvent) normalize(repo *defs.Repo) *defs.Event[defs.PullRequestReview, defs.RepoProvider] {
	id, version, ctx, sub := prelude(repo)
	event := &defs.Event[defs.PullRequestReview, defs.RepoProvider]{
		ID:      id,
		Version: version,
		Context: ctx,
		Subject: sub,
		Payload: pre.payload(),
	}

	event.SetSource(pre.Repository.URL)
	event.SetScopePullRequestReview()

	switch pre.Action {
	case "submitted":
		event.SetActionCreated()
	case "edited":
		event.SetActionUpdated()
	case "dismissed":
		event.SetActionDismissed()
	default:
		log.Warnf("unknown pull request review event action: %s", pre.Action)
		return nil
	}

	return event
}

// payload converts the PullRequestReviewCommentEvent struct to the relevant EventPayload.
//
// It returns a `defs.PullRequestComment` struct containing the relevant information for a pull request review comment
// event.
func (pre PullRequestReviewCommentEvent) payload() defs.PullRequestComment {
	return defs.PullRequestComment{
		ID:                pre.Comment.ID,
		PullRequestNumber: pre.Number,
		ReviewID:          pre.Comment.PullRequestReviewID,
		InReplyTo:         pre.Comment.InReplyTo,
		CommitSHA:         pre.Comment.CommitID,
		Path:              pre.Comment.Path,
		Position:          pre.Comment.Position,
		AuthorID:          pre.Comment.User.ID,
		Timestamp:         pre.Comment.UpdatedAt.Time(),
	}
}

// normalize converts the PullRequestReviewCommentEvent struct to an Event struct.
//
// It uses the provided Repo struct to extract relevant information for the EventContext and EventSubject. The action is
// set based on the `Action` field of the PullRequestReviewCommentEvent struct.
func (pre PullRequestReviewCommentEvent) normalize(
	repo *defs.Repo,
) *defs.Event[defs.PullRequestComment, defs.RepoProvider] {
	id, version, ctx, sub := prelude(repo)
	event := &defs.Event[defs.PullRequestComment, defs.RepoProvider]{
		ID:      id,
		Version: version,
		Context: ctx,
		Subject: sub,
		Payload: pre.payload(),
	}

	event.SetSource(pre.Repository.URL)
	event.SetScopePullRequestComment()

	switch pre.Action {
	case "created":
		event.SetActionCreated()
	case "edited":
		event.SetActionUpdated()
	case "deleted": //nolint
		event.SetActionDeleted()
	default:
		return nil
	}

	return event
}

func PrepareRepoEventPayload(event RepoEvent) *RepoEventMetadataQuery {
	return &RepoEventMetadataQuery{
		RepoID:         event.RepoID(),
		RepoName:       event.RepoName(),
		InstallationID: event.InstallationID(),
		SenderID:       event.SenderID(),
	}
}
