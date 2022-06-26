package webhooks

// GithubEvent defines a GitHub hook event type
type GithubEvent string

// GitHub hook types
const (
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
