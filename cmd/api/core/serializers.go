package core

import "github.com/gocql/gocql"

type (
	AppCreateRequest struct {
		Name string `json:"name"`
	}

	AppRepoCreateRequest struct {
		RepoID        gocql.UUID `json:"repo_id"`
		DefaultBranch string     `json:"default_branch"`
		IsMonorepo    bool       `json:"is_monorepo"`
		Provider      string     `json:"provider"`
	}
)
