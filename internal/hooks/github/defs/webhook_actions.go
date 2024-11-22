package defs

const (
	NoCommit = "0000000000000000000000000000000000000000"
)

// Installation actions.
const (
	InstallationCreated                = "created"
	InstallationDeleted                = "deleted"
	InstallationUnsuspended            = "unsuspended"
	InstallationSuspended              = "suspended"
	InstallationNewPermissionsAccepted = "new_permissions_accepted"
)

// Installation repository actions.
const (
	InstallationRepositoriesAdded   = "added"
	InstallationRepositoriesRemoved = "removed"
)

// Installation target actions.
const (
	InstallationTargetRenamed = "renamed"
)

// GitHub App authorization actions.
const (
	GitHubAppAuthorizationRevoked = "revoked"
)

// Membership actions.
const (
	MembershipAdded   = "added"
	MembershipRemoved = "removed"
)

// Member actions.
const (
	MemberAdded                 = "added"
	MemberEdited                = "edited"
	MemberRemoved               = "removed"
	MemberInvited               = "invited"
	MemberRemovedFromRepository = "removed_from_repository"
)

// Team actions.
const (
	TeamAddedToRepository     = "added_to_repository"
	TeamRemovedFromRepository = "removed_from_repository"
)

// Organization actions.
const (
	OrganizationMemberAdded   = "member_added"
	OrganizationMemberInvited = "member_invited"
	OrganizationMemberRemoved = "member_removed"
	OrganizationDeleted       = "deleted"
	OrganizationRenamed       = "renamed"
	OrganizationSuspended     = "suspended"
	OrganizationUnsuspended   = "unsuspended"
)

// Project actions.
const (
	ProjectClosed      = "closed"
	ProjectConverted   = "converted"
	ProjectCreated     = "created"
	ProjectEdited      = "edited"
	ProjectReopened    = "reopened"
	ProjectTransferred = "transferred"
)

// Project column actions.
const (
	ProjectColumnCreated = "created"
	ProjectColumnEdited  = "edited"
	ProjectColumnMoved   = "moved"
	ProjectColumnDeleted = "deleted"
)

// Project card actions.
const (
	ProjectCardConverted = "converted"
	ProjectCardCreated   = "created"
	ProjectCardEdited    = "edited"
	ProjectCardMoved     = "moved"
	ProjectCardDeleted   = "deleted"
)

// Milestone actions.
const (
	MilestoneClosed      = "closed"
	MilestoneCreated     = "created"
	MilestoneDeleted     = "deleted"
	MilestoneEdited      = "edited"
	MilestoneOpened      = "opened"
	MilestoneRenamed     = "renamed"
	MilestoneTransferred = "transferred"
)

// Label actions.
const (
	LabelCreated = "created"
	LabelEdited  = "edited"
	LabelDeleted = "deleted"
)

// Issue actions.
const (
	IssueAssigned     = "assigned"
	IssueClosed       = "closed"
	IssueCreated      = "created"
	IssueDemilestoned = "demilestoned"
	IssueEdited       = "edited"
	IssueLabeled      = "labeled"
	IssueLocked       = "locked"
	IssueMilestoned   = "milestoned"
	IssueOpened       = "opened"
	IssueReopened     = "reopened"
	IssueTransferred  = "transferred"
	IssueUnlabeled    = "unlabeled"
	IssueUnlocked     = "unlocked"
)

// Pull request actions.
const (
	PullRequestAssigned           = "assigned"
	PullRequestAutoMergeEnabled   = "auto_merge_enabled"
	PullRequestAutoMergeDisabled  = "auto_merge_disabled"
	PullRequestAutoMergeRejected  = "auto_merge_rejected"
	PullRequestAutoSquashEnabled  = "auto_squash_enabled"
	PullRequestAutoSquashDisabled = "auto_squash_disabled"
	PullRequestClosed             = "closed"
	PullRequestConverted          = "converted"
	PullRequestCreated            = "created"
	PullRequestDemilestoned       = "demilestoned"
	PullRequestEdited             = "edited"
	PullRequestLabeled            = "labeled"
	PullRequestLocked             = "locked"
	PullRequestMilestoned         = "milestoned"
	PullRequestOpened             = "opened"
	PullRequestReadyForReview     = "ready_for_review"
	PullRequestReopened           = "reopened"
	PullRequestTransferred        = "transferred"
	PullRequestUnlabeled          = "unlabeled"
	PullRequestUnlocked           = "unlocked"
)

// Pull request review actions.
const (
	PullRequestReviewSubmitted = "submitted"
	PullRequestReviewDismissed = "dismissed"
	PullRequestReviewEdited    = "edited"
)

// Pull request review comment actions.
const (
	PullRequestReviewCommentCreated = "created"
	PullRequestReviewCommentEdited  = "edited"
	PullRequestReviewCommentDeleted = "deleted"
)

// Pull request review thread actions.
const (
	PullRequestReviewThreadResolved    = "resolved"
	PullRequestReviewThreadUnresolved  = "unresolved"
	PullRequestReviewThreadRerequested = "rerequested"
)

// Commit comment actions.
const (
	CommitCommentCreated = "created"
	CommitCommentEdited  = "edited"
	CommitCommentDeleted = "deleted"
)

// Gollum actions.
const (
	GollumCreated = "created"
	GollumEdited  = "edited"
	GollumDeleted = "deleted"
)

// Release actions.
const (
	ReleasePublished = "published"
	ReleaseCreated   = "created"
	ReleaseEdited    = "edited"
	ReleaseDeleted   = "deleted"
	ReleaseReleased  = "released"
)

// Push actions.
const (
	PushCreated = "created"
	PushDeleted = "deleted"
)

// Deployment actions.
const (
	DeploymentCreated = "created"
)

// Deployment status actions.
const (
	DeploymentStatusCreated = "created"
)

// Deployment review actions.
const (
	DeploymentReviewApproved  = "approved"
	DeploymentReviewRejected  = "rejected"
	DeploymentReviewCommented = "commented"
	DeploymentReviewDismissed = "dismissed"
)

// Meta actions.
const (
	MetaDeleted = "deleted"
)

// Dependabot alert actions.
const (
	DependabotAlertAutoDismissed = "auto_dismissed"
)

// Security advisory actions.
const (
	SecurityAdvisoryPublished = "published"
)

// Repository actions.
const (
	RepositoryArchived    = "archived"
	RepositoryCreated     = "created"
	RepositoryDeleted     = "deleted"
	RepositoryEdited      = "edited"
	RepositoryPublished   = "published"
	RepositoryPrivatized  = "privatized"
	RepositoryRenamed     = "renamed"
	RepositoryTransferred = "transferred"
	RepositoryUnarchived  = "unarchived"
	RepositoryUnpublished = "unpublished"
)

// Org block actions.
const (
	OrgBlockBlocked   = "blocked"
	OrgBlockUnblocked = "unblocked"
)

// Watch actions.
const (
	WatchStarted = "started"
)

// Star actions.
const (
	StarCreated = "created"
	StarDeleted = "deleted"
)

// Check run actions.
const (
	CheckRunCompleted = "completed"
)

// Check suite actions.
const (
	CheckSuiteCompleted = "completed"
)

// Package actions.
const (
	PackagePublished = "published"
)

// Page build actions.
const (
	PageBuildCreated = "created"
)

// Repository dispatch actions.
const (
	RepositoryDispatch = "repository_dispatch"
)

// Workflow run actions.
const (
	WorkflowRunCompleted = "completed"
)

// Workflow job actions.
const (
	WorkflowJobCompleted = "completed"
)

// Workflow dispatch actions.
const (
	WorkflowDispatch = "workflow_dispatch"
)

// Secret scanning alert actions.
const (
	SecretScanningAlertCreated = "created"
)

// Secret scanning alert location actions.
const (
	SecretScanningAlertLocationCreated = "created"
)

// Code scanning alert actions.
const (
	CodeScanningAlertAppearedInBranch = "appeared_in_branch"
	CodeScanningAlertReopenedByUser   = "reopened_by_user"
	CodeScanningAlertClosedByUser     = "closed_by_user"
)

// Repository advisory actions.
const (
	RepositoryAdvisoryPublished = "published"
)

// Repository ruleset actions.
const (
	RepositoryRulesetCreated = "created"
	RepositoryRulesetEdited  = "edited"
	RepositoryRulesetDeleted = "deleted"
)

// Merge group actions.
const (
	MergeGroupChecksRequested = "checks_requested"
)

// Discussion actions.
const (
	DiscussionCreated  = "created"
	DiscussionEdited   = "edited"
	DiscussionDeleted  = "deleted"
	DiscussionPinned   = "pinned"
	DiscussionUnpinned = "unpinned"
	DiscussionAnswered = "answered"
)

// Discussion comment actions.
const (
	DiscussionCommentCreated = "created"
	DiscussionCommentEdited  = "edited"
	DiscussionCommentDeleted = "deleted"
)

// Personal Access Token Request actions.
const (
	PersonalAccessTokenRequestApproved = "approved"
	PersonalAccessTokenRequestDenied   = "denied"
	PersonalAccessTokenRequestRevoked  = "revoked"
)

// Sponsorship actions.
const (
	SponsorshipCreated         = "created"
	SponsorshipCancelled       = "cancelled"
	SponsorshipEdited          = "edited"
	SponsorshipPendingApproval = "pending_approval"
	SponsorshipTierChanged     = "tier_changed"
	SponsorshipPaymentSuccess  = "payment_success"
	SponsorshipPaymentFailed   = "payment_failed"
)

// Projects V2 actions.
const (
	ProjectsV2Closed   = "closed"
	ProjectsV2Created  = "created"
	ProjectsV2Edited   = "edited"
	ProjectsV2Reopened = "reopened"
)

// Projects V2 item actions.
const (
	ProjectsV2ItemArchived   = "archived"
	ProjectsV2ItemCreated    = "created"
	ProjectsV2ItemEdited     = "edited"
	ProjectsV2ItemMoved      = "moved"
	ProjectsV2ItemUnarchived = "unarchived"
)

// Projects V2 status update actions.
const (
	ProjectsV2StatusUpdateCreated = "created"
	ProjectsV2StatusUpdateEdited  = "edited"
	ProjectsV2StatusUpdateDeleted = "deleted"
)

// Security & Analysis actions.
const (
	SecurityAndAnalysisEnabled  = "enabled"
	SecurityAndAnalysisDisabled = "disabled"
)
