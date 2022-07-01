package common

import "time"

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
