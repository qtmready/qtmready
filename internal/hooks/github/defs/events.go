package defs

import (
	"time"
)

type (
	Push struct {
		Ref          string       `json:"ref"`
		Before       string       `json:"before"`
		After        string       `json:"after"`
		Created      bool         `json:"created"`
		Deleted      bool         `json:"deleted"`
		Forced       bool         `json:"forced"`
		BaseRef      *string      `json:"base_ref"`
		Compare      string       `json:"compare"`
		Commits      Commits      `json:"commits"`
		HeadCommit   Commit       `json:"head_commit"`
		Repository   Repository   `json:"repository"`
		Pusher       Pusher       `json:"pusher"`
		Sender       User         `json:"sender"`
		Installation Installation `json:"installation"`
	}

	// Pull Request event.
	PR struct {
		Action       string         `json:"action"`
		Number       int64          `json:"number"`
		PullRequest  PullRequest    `json:"pull_request"`
		Repository   RepositoryPR   `json:"repository"`
		Organization *Organization  `json:"organization"`
		Installation InstallationID `json:"installation"`
		Sender       User           `json:"sender"`
		Label        *Label         `json:"label"`
	}

	PrReview struct {
		Action       string             `json:"action"`
		Number       int64              `json:"number"`
		Installation InstallationID     `json:"installation"`
		Review       *PullRequestReview `json:"review"`
		PullRequest  PullRequest        `json:"pull_request"`
		Repository   RepositoryPR       `json:"repository"`
		Sender       *User              `json:"sender"`
	}

	PrReviewComment struct {
		Action       string              `json:"action"`
		Number       int64               `json:"number"`
		Installation InstallationID      `json:"installation"`
		Comment      *PullRequestComment `json:"comment"`
		PullRequest  PullRequest         `json:"pull_request"`
		Repository   RepositoryPR        `json:"repository"`
		Sender       *User               `json:"sender"`
	}
)

// ---------------------------------- Push Event ----------------------------------.
func (p *Push) GetRef() string {
	return p.Ref
}

func (p *Push) GetBefore() string {
	return p.Before
}

func (p *Push) GetAfter() string {
	return p.After
}

func (p *Push) GetRepositoryName() string {
	return p.Repository.Name
}

func (p *Push) GetSenderID() int64 {
	return p.Sender.ID
}

func (p *Push) GetCommits() Commits {
	return p.Commits
}

func (p *Push) GetRepositoryID() int64 {
	return p.Repository.ID
}

func (p *Push) GetInstallationID() int64 {
	return p.Installation.ID
}

func (p *Push) GetPusherEmail() string {
	return p.Pusher.Email
}

// ---------------------------------- Pull Request Event ----------------------------------.
func (pr *PR) GetTitle() string {
	return pr.PullRequest.Title
}

func (pr *PR) GetAction() string {
	return pr.Action
}

func (pr *PR) GetNumber() int64 {
	return pr.Number
}

func (pr *PR) GetBody() string {
	return pr.PullRequest.Body
}

func (pr *PR) GetAuthor() string {
	return pr.Sender.Login
}

func (pr *PR) GetHeadBranch() string {
	return pr.PullRequest.Head.Ref
}

func (pr *PR) GetBaseBranch() string {
	return pr.PullRequest.Base.Ref
}

func (pr *PR) GetTimestamp() time.Time {
	return pr.PullRequest.UpdatedAt
}

func (pr *PR) GetRepositoryID() int64 {
	return pr.Repository.ID
}

func (pr *PR) GetInstallationID() int64 {
	return pr.Installation.ID
}

func (pr *PR) GetSenderEmail() *string {
	return pr.Sender.Email
}

func (pr *PR) GetLabelName() string {
	return pr.Label.Name
}

// ---------------------------------- Pull Request Review Event ----------------------------------.
func (prr *PrReview) GetAction() string {
	return prr.Action
}

func (prr *PrReview) GetPrNumber() int64 {
	return prr.Number
}

func (prr *PrReview) GetRepositoryID() int64 {
	return prr.Repository.ID
}

func (prr *PrReview) GetInstallationID() int64 {
	return prr.Installation.ID
}

func (prr *PrReview) GetSenderEmail() *string {
	return prr.Sender.Email
}

func (prr *PrReview) GetHeadBranch() string {
	return prr.PullRequest.Head.Ref
}

func (prr *PrReview) GetPrReviewID() int64 {
	return prr.Review.ID
}

func (prr *PrReview) GetSubmittedAt() time.Time {
	return prr.Review.SubmittedAt
}

func (prr *PrReview) GetState() string {
	return prr.Review.State
}

// ---------------------------------- Pull Request Review Comment Event ----------------------------------.
func (prrc *PrReviewComment) GetAction() string {
	return prrc.Action
}

func (prrc *PrReviewComment) GetCommentID() int64 {
	return prrc.Comment.ID
}

func (prrc *PrReviewComment) GetPrNumber() int64 {
	return prrc.Number
}

func (prrc *PrReviewComment) GetState() string {
	return prrc.Comment.Body // TODO - need to decide.
}

func (prrc *PrReviewComment) GetInReplyTo() *int64 {
	return prrc.Comment.InReplyTo
}

func (prrc *PrReviewComment) GetCommitSha() string {
	return prrc.Comment.CommitID
}

func (prrc *PrReviewComment) GetReviewID() int64 {
	return prrc.Comment.PullRequestReviewID
}

func (prrc *PrReviewComment) GetPath() string {
	return prrc.Comment.Path
}

func (prrc *PrReviewComment) GetPosition() int64 {
	return prrc.Comment.Position
}

func (prrc *PrReviewComment) GetRepositoryID() int64 {
	return prrc.Repository.ID
}

func (prrc *PrReviewComment) GetHeadBranch() string {
	return prrc.PullRequest.Head.Ref
}

func (prrc *PrReviewComment) GetInstallationID() int64 {
	return prrc.Installation.ID
}

func (prrc *PrReviewComment) GetSenderEmail() *string {
	return prrc.Sender.Email
}

func (prrc *PrReviewComment) GetSubmittedAt() time.Time {
	return prrc.Comment.CreatedAt
}
