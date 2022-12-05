// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package github

import (
	"encoding/json"
	"errors"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
)

// Webhook events & Workflow singals.
type (
	WebhookEvent         string                               // WebhookEvent defines the event type.
	WebhookEventHandler  func(ctx echo.Context) error         // EventHandler is the signature of the event handler function.
	WebhookEventHandlers map[WebhookEvent]WebhookEventHandler // EventHandlers maps event types to their respective event handlers.
	WorkflowSignal       string                               // WorkflowSignal is the name of a workflow signal.
	WorkflowSignalMap    map[string]WorkflowSignal            // WorkflowSignalMap maps strings to their respective signal.
)

// maps for openapi generated enums.
type (
	SetupActionMap    map[string]SetupAction    // SetupActionMap maps strings to their respective SetupAction.
	WorkflowStatusMap map[string]WorkflowStatus // WorkflowStatusMap maps strings to their respective WorkflowStatus.
)

// Payloads for internal events & signals.
type (
	AppAuthorizationEventPayload struct {
		Action string `json:"action"`
		Sender User   `json:"sender"`
	}

	InstallationEventPayload struct {
		Action       string              `json:"action"`
		Installation Installation        `json:"installation"`
		Repositories []PartialRepository `json:"repositories"`
		Sender       User                `json:"sender"`
	}

	PushEventPayload struct {
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

	PullRequestEventPayload struct {
		Action       string                `json:"action"`
		Number       int64                 `json:"number"`
		PullRequest  PullRequest           `json:"pull_request"`
		Repository   PullRequestRepository `json:"repository"`
		Organization *Organization         `json:"organization"`
		Installation InstallationID        `json:"installation"`
		Sender       User                  `json:"sender"`
	}

	CompleteInstallationSignalPayload struct {
		InstallationID int64       `json:"installation_id"`
		SetupAction    SetupAction `json:"setup_action"`
		TeamID         gocql.UUID  `json:"team_id"`
	}
)

var (
	SetupActions = SetupActionMap{
		SetupActionCreated.String(): SetupActionCreated,
		SetupActionDeleted.String(): SetupActionDeleted,
		SetupActionUpdated.String(): SetupActionUpdated,
	}

	WorkflowSignals = WorkflowSignalMap{
		WorkflowSignalInstallationEvent.String():    WorkflowSignalInstallationEvent,
		WorkflowSignalCompleteInstallation.String(): WorkflowSignalCompleteInstallation,
		WorkflowSignalPullRequest.String():          WorkflowSignalPullRequest,
	}

	WorkflowStatuses = WorkflowStatusMap{
		WorkflowStatusFailure.String():  WorkflowStatusFailure,
		WorkflowStatusQueued.String():   WorkflowStatusQueued,
		WorkflowStatusSuccess.String():  WorkflowStatusSuccess,
		WorkflowStatusSignaled.String(): WorkflowStatusSignaled,
		WorkflowStatusSkipped.String():  WorkflowStatusSkipped,
	}
)

// Webhook event types. We get this from the header `X-Github-Event`.
// For payload information, see https://developer.github.com/webhooks/event-payloads/.
const (
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

// Workflow signal types.
const (
	WorkflowSignalInstallationEvent    WorkflowSignal = "installation_event"
	WorkflowSignalCompleteInstallation WorkflowSignal = "complete_installation"
	WorkflowSignalPullRequest          WorkflowSignal = "pull_request"
)

const (
	NoCommit = "0000000000000000000000000000000000000000"
)

func (e WebhookEvent) String() string { return string(e) }

// Methods for SetupAction.

func (a SetupAction) String() string { return string(a) }

func (a SetupAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *SetupAction) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := SetupActions[s]
	if !ok {
		return errors.New("invalid setup action")
	}

	*a = val

	return nil
}

/*
 * Methods for WorkflowSignal.
 */

func (w WorkflowSignal) String() string { return string(w) }

func (w WorkflowSignal) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

func (w *WorkflowSignal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := WorkflowSignals[s]
	if !ok {
		return errors.New("invalid workflow signal")
	}

	*w = val

	return nil
}

/*
 * Methods for WorkflowStatus.
 */

func (w WorkflowStatus) String() string { return string(w) }

func (w WorkflowStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

func (w *WorkflowStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val, ok := WorkflowStatuses[s]
	if !ok {
		return errors.New("invalid workflow status")
	}

	*w = val

	return nil
}
