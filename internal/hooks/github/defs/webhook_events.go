package defs

import (
	"encoding/json"
)

type (
	// WebhookEvent defines a Github Webhook event name.
	WebhookEvent string
)

const (
	WebhookEventUnspecified                         WebhookEvent = ""
	WebhookEventAppAuthorization                    WebhookEvent = "github_app_authorization"
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

var (
	WebhookEventMap = map[string]WebhookEvent{
		"github_app_authorization":              WebhookEventAppAuthorization,
		"check_run":                             WebhookEventCheckRun,
		"check_suite":                           WebhookEventCheckSuite,
		"commit_comment":                        WebhookEventCommitComment,
		"create":                                WebhookEventCreate,
		"delete":                                WebhookEventDelete,
		"deploy_key":                            WebhookEventDeployKey,
		"deployment":                            WebhookEventDeployment,
		"deployment_status":                     WebhookEventDeploymentStatus,
		"fork":                                  WebhookEventFork,
		"gollum":                                WebhookEventGollum,
		"installation":                          WebhookEventInstallation,
		"installation_repositories":             WebhookEventInstallationRepositories,
		"integration_installation":              WebhookEventIntegrationInstallation,
		"integration_installation_repositories": WebhookEventIntegrationInstallationRepositories,
		"issue_comment":                         WebhookEventIssueComment,
		"issues":                                WebhookEventIssues,
		"label":                                 WebhookEventLabel,
		"member":                                WebhookEventMember,
		"membership":                            WebhookEventMembership,
		"milestone":                             WebhookEventMilestone,
		"meta":                                  WebhookEventMeta,
		"organization":                          WebhookEventOrganization,
		"org_block":                             WebhookEventOrgBlock,
		"page_build":                            WebhookEventPageBuild,
		"ping":                                  WebhookEventPing,
		"project_card":                          WebhookEventProjectCard,
		"project_column":                        WebhookEventProjectColumn,
		"project":                               WebhookEventProject,
		"public":                                WebhookEventPublic,
		"pull_request":                          WebhookEventPullRequest,
		"pull_request_review":                   WebhookEventPullRequestReview,
		"pull_request_review_comment":           WebhookEventPullRequestReviewComment,
		"push":                                  WebhookEventPush,
		"release":                               WebhookEventRelease,
		"repository":                            WebhookEventRepository,
		"repository_vulnerability_alert":        WebhookEventRepositoryVulnerabilityAlert,
		"security_advisory":                     WebhookEventSecurityAdvisory,
		"status":                                WebhookEventStatus,
		"team":                                  WebhookEventTeam,
		"team_add":                              WebhookEventTeamAdd,
		"watch":                                 WebhookEventWatch,
		"workflow_dispatch":                     WebhookEventWorkflowDispatch,
		"workflow_job":                          WebhookEventWorkflowJob,
		"workflow_run":                          WebhookEventWorkflowRun,
	}
)

func (e WebhookEvent) String() string { return string(e) }

func (e WebhookEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

func (e *WebhookEvent) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if event, ok := WebhookEventMap[s]; ok {
		*e = event
	} else {
		*e = WebhookEventUnspecified
	}

	return nil
}
