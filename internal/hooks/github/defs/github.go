package githubdefs

import (
	"time"
)

type (
	Installation struct {
		ID                  int64      `json:"id"`
		NodeID              string     `json:"node_id"`
		Account             User       `json:"account"`
		RepositorySelection string     `json:"repository_selection"`
		AccessTokensURL     string     `json:"access_tokens_url"`
		RepositoriesURL     string     `json:"repositories_url"`
		HTMLUrl             string     `json:"html_url"`
		AppID               int64      `json:"app_id"`
		TargetID            int64      `json:"target_id"`
		TargetType          string     `json:"target_type"`
		Permissions         Permission `json:"permissions"`
		Events              []string   `json:"events"`
		CreatedAt           time.Time  `json:"created_at"`
		UpdatedAt           time.Time  `json:"updated_at"`
		SingleFileName      *string    `json:"single_file_name"`
	}

	PartialRepository struct {
		ID       int64  `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
	}

	Permission struct {
		Issues             string `json:"issues"`
		Metadata           string `json:"metadata"`
		PullRequests       string `json:"pull_requests"`
		RepositoryProjects string `json:"repository_projects"`
	}

	// User represents a Github User.
	User struct {
		Login             string `json:"login"`
		ID                int64  `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLUrl           string `json:"html_url"`
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

	UserPartial struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	Commit struct {
		SHA       string      `json:"sha"`
		ID        string      `json:"id"`
		NodeID    string      `json:"node_id"`
		TreeID    string      `json:"tree_id"`
		Distinct  bool        `json:"distinct"`
		Message   string      `json:"message"`
		Timestamp time.Time   `json:"timestamp"`
		URL       string      `json:"url"`
		Author    UserPartial `json:"author"`
		Committer UserPartial `json:"committer"`
		Added     []string    `json:"added"`
		Removed   []string    `json:"removed"`
		Modified  []string    `json:"modified"`
	}

	Commits []Commit

	Push struct {
		Ref          string     `json:"ref"`
		Before       string     `json:"before"`
		After        string     `json:"after"`
		Created      bool       `json:"created"`
		Deleted      bool       `json:"deleted"`
		Forced       bool       `json:"forced"`
		BaseRef      *string    `json:"base_ref"`
		Compare      string     `json:"compare"`
		Commits      []Commit   `json:"commits"`
		HeadCommit   Commit     `json:"head_commit"`
		Repository   Repository `json:"repository"`
		Pusher       Pusher     `json:"pusher"`
		Sender       User       `json:"sender"`
		Installation int64      `json:"installation"`
	}

	Repository struct {
		ID               int64     `json:"id"`
		NodeID           string    `json:"node_id"`
		Name             string    `json:"name"`
		FullName         string    `json:"full_name"`
		Owner            User      `json:"owner"`
		Private          bool      `json:"private"`
		HTMLUrl          string    `json:"html_url"`
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
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
		PushedAt         time.Time `json:"pushed_at"`
		GitURL           string    `json:"git_url"`
		SSHUrl           string    `json:"ssh_url"`
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
)
