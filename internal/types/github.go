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
	Action string `json:"action"`
	Sender sender `json:"sender"`
}

// GithubInstallationEventPayload contains the information for GitHub's installation and integration_installation hook events
type GithubInstallationEventPayload struct {
	Action       string             `json:"action"`
	Installation installation       `json:"installation"`
	Repositories []installationRepo `json:"repositories"`
	Sender       sender             `json:"sender"`
}

type GithubPushEventPayload struct {
	Ref          string          `json:"ref"`
	Before       string          `json:"before"`
	After        string          `json:"after"`
	Created      bool            `json:"created"`
	Deleted      bool            `json:"deleted"`
	Forced       bool            `json:"forced"`
	BaseRef      *string         `json:"base_ref"`
	Compare      string          `json:"compare"`
	Commits      []commit        `json:"commits"`
	HeadCommit   headCommit      `json:"head_commit"`
	Repository   pushRepo        `json:"repository"`
	Pusher       pusher          `json:"pusher"`
	Sender       sender          `json:"sender"`
	Installation pushInstllation `json:"installation"`
}

// types local to package

type pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type pushInstllation struct {
	ID int `json:"id"`
}

type commit struct {
	Sha       string     `json:"sha"`
	ID        string     `json:"id"`
	NodeID    string     `json:"node_id"`
	TreeID    string     `json:"tree_id"`
	Distinct  bool       `json:"distinct"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	URL       string     `json:"url"`
	Author    githubUser `json:"author"`
	Committer githubUser `json:"committer"`
	Added     []string   `json:"added"`
	Removed   []string   `json:"removed"`
	Modified  []string   `json:"modified"`
}

type headCommit struct {
	ID        string     `json:"id"`
	NodeID    string     `json:"node_id"`
	TreeID    string     `json:"tree_id"`
	Distinct  bool       `json:"distinct"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	URL       string     `json:"url"`
	Author    githubUser `json:"author"`
	Committer githubUser `json:"committer"`
	Added     []string   `json:"added"`
	Removed   []string   `json:"removed"`
	Modified  []string   `json:"modified"`
}

type githubUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type pushRepo struct {
	ID               int64     `json:"id"`
	NodeID           string    `json:"node_id"`
	Name             string    `json:"name"`
	FullName         string    `json:"full_name"`
	Owner            repoOwner `json:"owner"`
	Private          bool      `json:"private"`
	HTMLURL          string    `json:"html_url"`
	Description      string    `json:"description"`
	Fork             bool      `json:"fork"`
	URL              string    `json:"url"`
	ForksURL         string    `json:"forks_url"`
	KeysURL          string    `json:"keys_url"`
	CollaboratorsURL string    `json:"collaborators_url"`
	TeamsURL         string    `json:"teams_url"`
	HooksURL         string    `json:"hooks_url"`
	IssueEventsURL   string    `json:"issue_events_url"`
	EventsURL        string    `json:"events_url"`
	AssigneesURL     string    `json:"assignees_url"`
	BranchesURL      string    `json:"branches_url"`
	TagsURL          string    `json:"tags_url"`
	BlobsURL         string    `json:"blobs_url"`
	GitTagsURL       string    `json:"git_tags_url"`
	GitRefsURL       string    `json:"git_refs_url"`
	TreesURL         string    `json:"trees_url"`
	StatusesURL      string    `json:"statuses_url"`
	LanguagesURL     string    `json:"languages_url"`
	StargazersURL    string    `json:"stargazers_url"`
	ContributorsURL  string    `json:"contributors_url"`
	SubscribersURL   string    `json:"subscribers_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommitsURL       string    `json:"commits_url"`
	GitCommitsURL    string    `json:"git_commits_url"`
	CommentsURL      string    `json:"comments_url"`
	IssueCommentURL  string    `json:"issue_comment_url"`
	ContentsURL      string    `json:"contents_url"`
	CompareURL       string    `json:"compare_url"`
	MergesURL        string    `json:"merges_url"`
	ArchiveURL       string    `json:"archive_url"`
	DownloadsURL     string    `json:"downloads_url"`
	IssuesURL        string    `json:"issues_url"`
	PullsURL         string    `json:"pulls_url"`
	MilestonesURL    string    `json:"milestones_url"`
	NotificationsURL string    `json:"notifications_url"`
	LabelsURL        string    `json:"labels_url"`
	ReleasesURL      string    `json:"releases_url"`
	CreatedAt        int64     `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PushedAt         int64     `json:"pushed_at"`
	GitURL           string    `json:"git_url"`
	SSHURL           string    `json:"ssh_url"`
	CloneURL         string    `json:"clone_url"`
	SvnURL           string    `json:"svn_url"`
	Homepage         *string   `json:"homepage"`
	Size             int64     `json:"size"`
	StargazersCount  int64     `json:"stargazers_count"`
	WatchersCount    int64     `json:"watchers_count"`
	Language         *string   `json:"language"`
	HasIssues        bool      `json:"has_issues"`
	HasDownloads     bool      `json:"has_downloads"`
	HasWiki          bool      `json:"has_wiki"`
	HasPages         bool      `json:"has_pages"`
	ForksCount       int64     `json:"forks_count"`
	MirrorURL        *string   `json:"mirror_url"`
	OpenIssuesCount  int64     `json:"open_issues_count"`
	Forks            int64     `json:"forks"`
	OpenIssues       int64     `json:"open_issues"`
	Watchers         int64     `json:"watchers"`
	DefaultBranch    string    `json:"default_branch"`
	Stargazers       int64     `json:"stargazers"`
	MasterBranch     string    `json:"master_branch"`
}

type repoOwner struct {
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

type sender struct {
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
