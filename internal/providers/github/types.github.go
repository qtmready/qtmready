package github

import (
	"time"

	"go.breu.io/quantm/internal/shared"
)

// Github embedded types.
type (
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	InstallationID struct {
		ID shared.Int64 `json:"id"`
	}

	Commit struct {
		SHA       string      `json:"sha"`
		ID        string      `json:"id"`
		NodeID    string      `json:"node_id"`
		TreeID    string      `json:"tree_id"`
		Distinct  bool        `json:"distinct"`
		Message   string      `json:"message"`
		Timestamp string      `json:"timestamp"`
		URL       string      `json:"url"`
		Author    UserPartial `json:"author"`
		Committer UserPartial `json:"committer"`
		Added     []string    `json:"added"`
		Removed   []string    `json:"removed"`
		Modified  []string    `json:"modified"`
	}

	HeadCommit struct {
		ID        string      `json:"id"`
		NodeID    string      `json:"node_id"`
		TreeID    string      `json:"tree_id"`
		Distinct  bool        `json:"distinct"`
		Message   string      `json:"message"`
		Timestamp string      `json:"timestamp"`
		URL       string      `json:"url"`
		Author    UserPartial `json:"author"`
		Committer UserPartial `json:"committer"`
		Added     []string    `json:"added"`
		Removed   []string    `json:"removed"`
		Modified  []string    `json:"modified"`
	}

	UserPartial struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	Repository struct {
		ID               shared.Int64 `json:"id"`
		NodeID           string       `json:"node_id"`
		Name             string       `json:"name"`
		FullName         string       `json:"full_name"`
		Owner            User         `json:"owner"`
		Private          bool         `json:"private"`
		HTMLUrl          string       `json:"html_url"`
		Description      string       `json:"description"`
		Fork             bool         `json:"fork"`
		URL              string       `json:"url"`
		ForksURL         string       `json:"forks_url"`
		KeysURL          string       `json:"keys_url"`
		CollaboratorsURL string       `json:"collaborators_url"`
		TeamsURL         string       `json:"teams_url"`
		HooksURL         string       `json:"hooks_url"`
		IssueEventsURL   string       `json:"issue_events_url"`
		EventsURL        string       `json:"events_url"`
		AssigneesURL     string       `json:"assignees_url"`
		BranchesURL      string       `json:"branches_url"`
		TagsURL          string       `json:"tags_url"`
		BlobsURL         string       `json:"blobs_url"`
		GitTagsURL       string       `json:"git_tags_url"`
		GitRefsURL       string       `json:"git_refs_url"`
		TreesURL         string       `json:"trees_url"`
		StatusesURL      string       `json:"statuses_url"`
		LanguagesURL     string       `json:"languages_url"`
		StargazersURL    string       `json:"stargazers_url"`
		ContributorsURL  string       `json:"contributors_url"`
		SubscribersURL   string       `json:"subscribers_url"`
		SubscriptionURL  string       `json:"subscription_url"`
		CommitsURL       string       `json:"commits_url"`
		GitCommitsURL    string       `json:"git_commits_url"`
		CommentsURL      string       `json:"comments_url"`
		IssueCommentURL  string       `json:"issue_comment_url"`
		ContentsURL      string       `json:"contents_url"`
		CompareURL       string       `json:"compare_url"`
		MergesURL        string       `json:"merges_url"`
		ArchiveURL       string       `json:"archive_url"`
		DownloadsURL     string       `json:"downloads_url"`
		IssuesURL        string       `json:"issues_url"`
		PullsURL         string       `json:"pulls_url"`
		MilestonesURL    string       `json:"milestones_url"`
		NotificationsURL string       `json:"notifications_url"`
		LabelsURL        string       `json:"labels_url"`
		ReleasesURL      string       `json:"releases_url"`
		CreatedAt        Timestamp    `json:"created_at"`
		UpdatedAt        Timestamp    `json:"updated_at"`
		PushedAt         Timestamp    `json:"pushed_at"`
		GitURL           string       `json:"git_url"`
		SSHUrl           string       `json:"ssh_url"`
		CloneURL         string       `json:"clone_url"`
		SvnURL           string       `json:"svn_url"`
		Homepage         *string      `json:"homepage"`
		Size             int64        `json:"size"`
		StargazersCount  int64        `json:"stargazers_count"`
		WatchersCount    int64        `json:"watchers_count"`
		Language         *string      `json:"language"`
		HasIssues        bool         `json:"has_issues"`
		HasDownloads     bool         `json:"has_downloads"`
		HasWiki          bool         `json:"has_wiki"`
		HasPages         bool         `json:"has_pages"`
		ForksCount       shared.Int64 `json:"forks_count"`
		MirrorURL        *string      `json:"mirror_url"`
		OpenIssuesCount  shared.Int64 `json:"open_issues_count"`
		Forks            shared.Int64 `json:"forks"`
		OpenIssues       shared.Int64 `json:"open_issues"`
		Watchers         shared.Int64 `json:"watchers"`
		DefaultBranch    string       `json:"default_branch"`
		Stargazers       shared.Int64 `json:"stargazers"`
		MasterBranch     string       `json:"master_branch"`
	}

	RepositoryPR struct {
		ID               shared.Int64 `json:"id"`
		NodeID           string       `json:"node_id"`
		Name             string       `json:"name"`
		FullName         string       `json:"full_name"`
		Owner            User         `json:"owner"`
		Private          bool         `json:"private"`
		HTMLUrl          string       `json:"html_url"`
		Description      string       `json:"description"`
		Fork             bool         `json:"fork"`
		URL              string       `json:"url"`
		ForksURL         string       `json:"forks_url"`
		KeysURL          string       `json:"keys_url"`
		CollaboratorsURL string       `json:"collaborators_url"`
		TeamsURL         string       `json:"teams_url"`
		HooksURL         string       `json:"hooks_url"`
		IssueEventsURL   string       `json:"issue_events_url"`
		EventsURL        string       `json:"events_url"`
		AssigneesURL     string       `json:"assignees_url"`
		BranchesURL      string       `json:"branches_url"`
		TagsURL          string       `json:"tags_url"`
		BlobsURL         string       `json:"blobs_url"`
		GitTagsURL       string       `json:"git_tags_url"`
		GitRefsURL       string       `json:"git_refs_url"`
		TreesURL         string       `json:"trees_url"`
		StatusesURL      string       `json:"statuses_url"`
		LanguagesURL     string       `json:"languages_url"`
		StargazersURL    string       `json:"stargazers_url"`
		ContributorsURL  string       `json:"contributors_url"`
		SubscribersURL   string       `json:"subscribers_url"`
		SubscriptionURL  string       `json:"subscription_url"`
		CommitsURL       string       `json:"commits_url"`
		GitCommitsURL    string       `json:"git_commits_url"`
		CommentsURL      string       `json:"comments_url"`
		IssueCommentURL  string       `json:"issue_comment_url"`
		ContentsURL      string       `json:"contents_url"`
		CompareURL       string       `json:"compare_url"`
		MergesURL        string       `json:"merges_url"`
		ArchiveURL       string       `json:"archive_url"`
		DownloadsURL     string       `json:"downloads_url"`
		IssuesURL        string       `json:"issues_url"`
		PullsURL         string       `json:"pulls_url"`
		MilestonesURL    string       `json:"milestones_url"`
		NotificationsURL string       `json:"notifications_url"`
		LabelsURL        string       `json:"labels_url"`
		ReleasesURL      string       `json:"releases_url"`
		CreatedAt        time.Time    `json:"created_at"`
		UpdatedAt        time.Time    `json:"updated_at"`
		PushedAt         time.Time    `json:"pushed_at"`
		GitURL           string       `json:"git_url"`
		SSHUrl           string       `json:"ssh_url"`
		CloneURL         string       `json:"clone_url"`
		SvnURL           string       `json:"svn_url"`
		Homepage         *string      `json:"homepage"`
		Size             shared.Int64 `json:"size"`
		StargazersCount  shared.Int64 `json:"stargazers_count"`
		WatchersCount    shared.Int64 `json:"watchers_count"`
		Language         *string      `json:"language"`
		HasIssues        bool         `json:"has_issues"`
		HasDownloads     bool         `json:"has_downloads"`
		HasWiki          bool         `json:"has_wiki"`
		HasPages         bool         `json:"has_pages"`
		ForksCount       shared.Int64 `json:"forks_count"`
		MirrorURL        *string      `json:"mirror_url"`
		OpenIssuesCount  shared.Int64 `json:"open_issues_count"`
		Forks            shared.Int64 `json:"forks"`
		OpenIssues       shared.Int64 `json:"open_issues"`
		Watchers         shared.Int64 `json:"watchers"`
		DefaultBranch    string       `json:"default_branch"`
		Stargazers       shared.Int64 `json:"stargazers"`
		MasterBranch     string       `json:"master_branch"`
	}

	User struct {
		Login             string       `json:"login"`
		ID                shared.Int64 `json:"id"`
		NodeID            string       `json:"node_id"`
		AvatarURL         string       `json:"avatar_url"`
		GravatarID        string       `json:"gravatar_id"`
		URL               string       `json:"url"`
		HTMLUrl           string       `json:"html_url"`
		FollowersURL      string       `json:"followers_url"`
		FollowingURL      string       `json:"following_url"`
		GistsURL          string       `json:"gists_url"`
		StarredURL        string       `json:"starred_url"`
		SubscriptionsURL  string       `json:"subscriptions_url"`
		OrganizationsURL  string       `json:"organizations_url"`
		ReposURL          string       `json:"repos_url"`
		EventsURL         string       `json:"events_url"`
		ReceivedEventsURL string       `json:"received_events_url"`
		Type              string       `json:"type"`
		SiteAdmin         bool         `json:"site_admin"`
	}

	Permission struct {
		Issues             string `json:"issues"`
		Metadata           string `json:"metadata"`
		PullRequests       string `json:"pull_requests"`
		RepositoryProjects string `json:"repository_projects"`
	}

	// MileStone contains GitHub's milestone information.
	MileStone struct {
		URL          string       `json:"url"`
		HTMLUrl      string       `json:"html_url"`
		LabelsURL    string       `json:"labels_url"`
		ID           shared.Int64 `json:"id"`
		NodeID       string       `json:"node_id"`
		Number       shared.Int64 `json:"number"`
		State        string       `json:"state"`
		Title        string       `json:"title"`
		Description  string       `json:"description"`
		Creator      User         `json:"creator"`
		OpenIssues   shared.Int64 `json:"open_issues"`
		ClosedIssues shared.Int64 `json:"closed_issues"`
		CreatedAt    time.Time    `json:"created_at"`
		UpdatedAt    time.Time    `json:"updated_at"`
		ClosedAt     time.Time    `json:"closed_at"`
		DueOn        time.Time    `json:"due_on"`
	}

	InstallationPayload struct {
		ID                  shared.Int64 `json:"id"`
		NodeID              string       `json:"node_id"`
		Account             User         `json:"account"`
		RepositorySelection string       `json:"repository_selection"`
		AccessTokensURL     string       `json:"access_tokens_url"`
		RepositoriesURL     string       `json:"repositories_url"`
		HTMLUrl             string       `json:"html_url"`
		AppID               int          `json:"app_id"`
		TargetID            int          `json:"target_id"`
		TargetType          string       `json:"target_type"`
		Permissions         Permission   `json:"permissions"`
		Events              []string     `json:"events"`
		CreatedAt           time.Time    `json:"created_at"`
		UpdatedAt           time.Time    `json:"updated_at"`
		SingleFileName      *string      `json:"single_file_name"`
	}

	PartialRepository struct {
		ID       shared.Int64 `json:"id"`
		NodeID   string       `json:"node_id"`
		Name     string       `json:"name"`
		FullName string       `json:"full_name"`
		Private  bool         `json:"private"`
	}

	PullRequestRef struct {
		Label string       `json:"label"`
		Ref   string       `json:"ref"`
		SHA   string       `json:"sha"`
		User  User         `json:"user"`
		Repo  RepositoryPR `json:"repo"`
	}

	Href struct {
		Href string `json:"href"`
	}

	PullRequestLinks struct {
		Self           Href `json:"self"`
		HTML           Href `json:"html"`
		Issue          Href `json:"issue"`
		Comments       Href `json:"comments"`
		ReviewComments Href `json:"review_comments"`
		ReviewComment  Href `json:"review_comment"`
		Commits        Href `json:"commits"`
		Statuses       Href `json:"statuses"`
	}

	PullRequest struct {
		URL                string           `json:"url"`
		ID                 shared.Int64     `json:"id"`
		NodeID             string           `json:"node_id"`
		HTMLUrl            string           `json:"html_url"`
		DiffURL            string           `json:"diff_url"`
		PatchURL           string           `json:"patch_url"`
		IssueURL           string           `json:"issue_url"`
		Number             shared.Int64     `json:"number"`
		State              string           `json:"state"`
		Locked             bool             `json:"locked"`
		Title              string           `json:"title"`
		User               User             `json:"user"`
		Body               string           `json:"body"`
		CreatedAt          time.Time        `json:"created_at"`
		UpdatedAt          time.Time        `json:"updated_at"`
		ClosedAt           *time.Time       `json:"closed_at"`
		MergedAt           *time.Time       `json:"merged_at"`
		MergeCommitSha     *string          `json:"merge_commit_sha"`
		Assignee           *User            `json:"assignee"`
		Assignees          []*User          `json:"assignees"`
		Milestone          *MileStone       `json:"milestone"`
		Draft              bool             `json:"draft"`
		CommitsURL         string           `json:"commits_url"`
		ReviewCommentsURL  string           `json:"review_comments_url"`
		ReviewCommentURL   string           `json:"review_comment_url"`
		CommentsURL        string           `json:"comments_url"`
		StatusesURL        string           `json:"statuses_url"`
		RequestedReviewers []User           `json:"requested_reviewers,omitempty"`
		Labels             []Label          `json:"labels"`
		Head               PullRequestRef   `json:"head"`
		Base               PullRequestRef   `json:"base"`
		Links              PullRequestLinks `json:"_links"`
		Merged             bool             `json:"merged"`
		Mergeable          *bool            `json:"mergeable"`
		MergeableState     string           `json:"mergeable_state"`
		MergedBy           *User            `json:"merged_by"`
		Comments           shared.Int64     `json:"comments"`
		ReviewComments     shared.Int64     `json:"review_comments"`
		Commits            shared.Int64     `json:"commits"`
		Additions          shared.Int64     `json:"additions"`
		Deletions          shared.Int64     `json:"deletions"`
		ChangedFiles       shared.Int64     `json:"changed_files"`
	}

	Label struct {
		ID          shared.Int64 `json:"id"`
		NodeID      string       `json:"node_id"`
		Description string       `json:"description"`
		URL         string       `json:"url"`
		Name        string       `json:"name"`
		Color       string       `json:"color"`
		Default     bool         `json:"default"`
	}

	RequestedTeam struct {
		Name            string       `json:"name"`
		ID              shared.Int64 `json:"id"`
		NodeID          string       `json:"node_id"`
		Slug            string       `json:"slug"`
		Description     string       `json:"description"`
		Privacy         string       `json:"privacy"`
		URL             string       `json:"url"`
		HTMLUrl         string       `json:"html_url"`
		MembersURL      string       `json:"members_url"`
		RepositoriesURL string       `json:"repositories_url"`
		Permission      string       `json:"permission"`
	}

	Organization struct {
		Login            string `json:"login"`
		ID               int    `json:"id"`
		NodeID           string `json:"node_id"`
		URL              string `json:"url"`
		ReposURL         string `json:"repos_url"`
		EventsURL        string `json:"events_url"`
		HooksURL         string `json:"hooks_url"`
		IssuesURL        string `json:"issues_url"`
		MembersURL       string `json:"members_url"`
		PublicMembersURL string `json:"public_members_url"`
		AvatarURL        string `json:"avatar_url"`
		Description      string `json:"description"`
	}

	// WorkflowRun represents a repository action workflow run.
	WorkflowRunPayload struct {
		ID                 shared.Int64   `json:"id,omitempty"`
		Name               string         `json:"name,omitempty"`
		NodeID             *string        `json:"node_id,omitempty"`
		HeadBranch         string         `json:"head_branch,omitempty"`
		HeadSHA            string         `json:"head_sha,omitempty"`
		RunNumber          int            `json:"run_number,omitempty"`
		RunAttempt         int            `json:"run_attempt,omitempty"`
		Event              string         `json:"event,omitempty"`
		DisplayTitle       *string        `json:"display_title,omitempty"`
		Status             string         `json:"status,omitempty"`
		Conclusion         *string        `json:"conclusion,omitempty"`
		WorkflowID         shared.Int64   `json:"workflow_id,omitempty"`
		CheckSuiteID       *shared.Int64  `json:"check_suite_id,omitempty"`
		CheckSuiteNodeID   *string        `json:"check_suite_node_id,omitempty"`
		URL                string         `json:"url,omitempty"`
		HTMLURL            *string        `json:"html_url,omitempty"`
		PullRequests       []*PullRequest `json:"pull_requests,omitempty"`
		CreatedAt          *time.Time     `json:"created_at,omitempty"`
		UpdatedAt          *time.Time     `json:"updated_at,omitempty"`
		RunStartedAt       *time.Time     `json:"run_started_at,omitempty"`
		JobsURL            *string        `json:"jobs_url,omitempty"`
		LogsURL            *string        `json:"logs_url,omitempty"`
		CheckSuiteURL      *string        `json:"check_suite_url,omitempty"`
		ArtifactsURL       *string        `json:"artifacts_url,omitempty"`
		CancelURL          *string        `json:"cancel_url,omitempty"`
		RerunURL           *string        `json:"rerun_url,omitempty"`
		PreviousAttemptURL *string        `json:"previous_attempt_url,omitempty"`
		HeadCommit         HeadCommit     `json:"head_commit,omitempty"`
		WorkflowURL        string         `json:"workflow_url,omitempty"`
		Repository         Repository     `json:"repository,omitempty"`
		HeadRepository     Repository     `json:"head_repository,omitempty"`
		Actor              User           `json:"actor,omitempty"`
	}

	// Workflow represents a repository action workflow.
	WorkflowPayload struct {
		ID        shared.Int64 `json:"id,omitempty"`
		NodeID    *string      `json:"node_id,omitempty"`
		Name      string       `json:"name,omitempty"`
		Path      string       `json:"path,omitempty"`
		State     string       `json:"state,omitempty"`
		CreatedAt *time.Time   `json:"created_at,omitempty"`
		UpdatedAt *time.Time   `json:"updated_at,omitempty"`
		URL       *string      `json:"url,omitempty"`
		HTMLURL   *string      `json:"html_url,omitempty"`
		BadgeURL  *string      `json:"badge_url,omitempty"`
	}

	PullRequestReview struct {
		ID                *int64     `json:"id,omitempty"`
		NodeID            *string    `json:"node_id,omitempty"`
		User              *User      `json:"user,omitempty"`
		Body              *string    `json:"body,omitempty"`
		SubmittedAt       *Timestamp `json:"submitted_at,omitempty"`
		CommitID          *string    `json:"commit_id,omitempty"`
		HTMLURL           *string    `json:"html_url,omitempty"`
		PullRequestURL    *string    `json:"pull_request_url,omitempty"`
		State             *string    `json:"state,omitempty"`
		AuthorAssociation *string    `json:"author_association,omitempty"`
	}

	PullRequestComment struct {
		ID                  *int64     `json:"id,omitempty"`
		NodeID              *string    `json:"node_id,omitempty"`
		InReplyTo           *int64     `json:"in_reply_to_id,omitempty"`
		Body                *string    `json:"body,omitempty"`
		Path                *string    `json:"path,omitempty"`
		DiffHunk            *string    `json:"diff_hunk,omitempty"`
		PullRequestReviewID *int64     `json:"pull_request_review_id,omitempty"`
		Position            *int       `json:"position,omitempty"`
		OriginalPosition    *int       `json:"original_position,omitempty"`
		StartLine           *int       `json:"start_line,omitempty"`
		Line                *int       `json:"line,omitempty"`
		OriginalLine        *int       `json:"original_line,omitempty"`
		OriginalStartLine   *int       `json:"original_start_line,omitempty"`
		Side                *string    `json:"side,omitempty"`
		StartSide           *string    `json:"start_side,omitempty"`
		CommitID            *string    `json:"commit_id,omitempty"`
		OriginalCommitID    *string    `json:"original_commit_id,omitempty"`
		User                *User      `json:"user,omitempty"`
		Reactions           *Reactions `json:"reactions,omitempty"`
		CreatedAt           *Timestamp `json:"created_at,omitempty"`
		UpdatedAt           *Timestamp `json:"updated_at,omitempty"`
		AuthorAssociation   *string    `json:"author_association,omitempty"`
		URL                 *string    `json:"url,omitempty"`
		HTMLURL             *string    `json:"html_url,omitempty"`
		PullRequestURL      *string    `json:"pull_request_url,omitempty"`
		SubjectType         *string    `json:"subject_type,omitempty"`
	}

	Reactions struct {
		TotalCount *int    `json:"total_count,omitempty"`
		PlusOne    *int    `json:"+1,omitempty"`
		MinusOne   *int    `json:"-1,omitempty"`
		Laugh      *int    `json:"laugh,omitempty"`
		Confused   *int    `json:"confused,omitempty"`
		Heart      *int    `json:"heart,omitempty"`
		Hooray     *int    `json:"hooray,omitempty"`
		Rocket     *int    `json:"rocket,omitempty"`
		Eyes       *int    `json:"eyes,omitempty"`
		URL        *string `json:"url,omitempty"`
	}
)
