package types

import "time"

// GithubEvent defines a GitHub hook event type
type GithubEvent string

// GitHub hook types
const (
	GithubAppAuthorizationEvent                    GithubEvent = "github_app_authorization"
	GithubCheckRunEvent                            GithubEvent = "check_run"
	GithubCheckSuiteEvent                          GithubEvent = "check_suite"
	GithubCommitCommentEvent                       GithubEvent = "commit_comment"
	GithubCreateEvent                              GithubEvent = "create"
	GithubDeleteEvent                              GithubEvent = "delete"
	GithubDeployKeyEvent                           GithubEvent = "deploy_key"
	GithubDeploymentEvent                          GithubEvent = "deployment"
	GithubDeploymentStatusEvent                    GithubEvent = "deployment_status"
	GithubForkEvent                                GithubEvent = "fork"
	GithubGollumEvent                              GithubEvent = "gollum"
	GithubInstallationEvent                        GithubEvent = "installation"
	GithubInstallationRepositoriesEvent            GithubEvent = "installation_repositories"
	GithubIntegrationInstallationEvent             GithubEvent = "integration_installation"
	GithubIntegrationInstallationRepositoriesEvent GithubEvent = "integration_installation_repositories"
	GithubIssueCommentEvent                        GithubEvent = "issue_comment"
	GithubIssuesEvent                              GithubEvent = "issues"
	GithubLabelEvent                               GithubEvent = "label"
	GithubMemberEvent                              GithubEvent = "member"
	GithubMembershipEvent                          GithubEvent = "membership"
	GithubMilestoneEvent                           GithubEvent = "milestone"
	GithubMetaEvent                                GithubEvent = "meta"
	GithubOrganizationEvent                        GithubEvent = "organization"
	GithubOrgBlockEvent                            GithubEvent = "org_block"
	GithubPageBuildEvent                           GithubEvent = "page_build"
	GithubPingEvent                                GithubEvent = "ping"
	GithubProjectCardEvent                         GithubEvent = "project_card"
	GithubProjectColumnEvent                       GithubEvent = "project_column"
	GithubProjectEvent                             GithubEvent = "project"
	GithubPublicEvent                              GithubEvent = "public"
	GithubPullRequestEvent                         GithubEvent = "pull_request"
	GithubPullRequestReviewEvent                   GithubEvent = "pull_request_review"
	GithubPullRequestReviewCommentEvent            GithubEvent = "pull_request_review_comment"
	GithubPushEvent                                GithubEvent = "push"
	GithubReleaseEvent                             GithubEvent = "release"
	GithubRepositoryEvent                          GithubEvent = "repository"
	GithubRepositoryVulnerabilityAlertEvent        GithubEvent = "repository_vulnerability_alert"
	GithubSecurityAdvisoryEvent                    GithubEvent = "security_advisory"
	GithubStatusEvent                              GithubEvent = "status"
	GithubTeamEvent                                GithubEvent = "team"
	GithubTeamAddEvent                             GithubEvent = "team_add"
	GithubWatchEvent                               GithubEvent = "watch"
	GithubWorkflowDispatchEvent                    GithubEvent = "workflow_dispatch"
	GithubWorkflowJobEvent                         GithubEvent = "workflow_job"
	GithubWorkflowRunEvent                         GithubEvent = "workflow_run"
)

// GithubEventSubtype defines a GitHub Hook Event subtype
type GithubEventSubtype string

// GitHub hook event subtypes
const (
	NoSubtype     GithubEventSubtype = ""
	BranchSubtype GithubEventSubtype = "branch"
	TagSubtype    GithubEventSubtype = "tag"
	PullSubtype   GithubEventSubtype = "pull"
	IssueSubtype  GithubEventSubtype = "issues"
)

type GithubAppAuthorizationEventPayload struct {
	Action string             `json:"action"`
	Sender installationSender `json:"sender"`
}

// GithubInstallationEventPayload contains the information for GitHub's installation and integration_installation hook events
type GithubInstallationEventPayload struct {
	Action       string             `json:"action"`
	Installation installation       `json:"installation"`
	Repositories []installationRepo `json:"repositories"`
	Sender       installationSender `json:"sender"`
}

// types local to package

type installationAccount struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type installationPermissions struct {
	Issues             string `json:"issues"`
	Metadata           string `json:"metadata"`
	PullRequests       string `json:"pull_requests"`
	RepositoryProjects string `json:"repository_projects"`
}

type installation struct {
	ID                  int64                   `json:"id"`
	NodeID              string                  `json:"node_id"`
	Account             installationAccount     `json:"account"`
	RepositorySelection string                  `json:"repository_selection"`
	AccessTokensURL     string                  `json:"access_tokens_url"`
	RepositoriesURL     string                  `json:"repositories_url"`
	HTMLURL             string                  `json:"html_url"`
	AppID               int                     `json:"app_id"`
	TargetID            int                     `json:"target_id"`
	TargetType          string                  `json:"target_type"`
	Permissions         installationPermissions `json:"permissions"`
	Events              []string                `json:"events"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
	SingleFileName      *string                 `json:"single_file_name"`
}

type installationRepo struct {
	ID       int64  `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

type installationSender struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
