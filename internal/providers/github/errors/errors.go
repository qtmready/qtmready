package errors

import (
	"errors"
	"fmt"

	"go.breu.io/quantm/internal/db"
)

type (
	MissingConfigurationError struct{}

	RepoEventError struct {
		InstallationID db.Int64 `json:"installation_id"`
		GithubRepoID   db.Int64 `json:"github_repo_id"`
		RepoName       string   `json:"repo_name"`
		Details        string   `json:"details"`
	}
)

var (
	ErrNoEventToParse               = errors.New("no event specified to parse")
	ErrInvalidHTTPMethod            = errors.New("invalid HTTP Method")
	ErrMissingHeaderGithubEvent     = errors.New("missing X-GitHub-Event Header")
	ErrMissingHeaderGithubSignature = errors.New("missing X-Hub-Signature Header")
	ErrInvalidEvent                 = errors.New("event not defined to be parsed")
	ErrPayloadParser                = errors.New("error parsing payload")
	ErrVerifySignature              = errors.New("HMAC verification failed")
)

func (e *RepoEventError) Error() string {
	return fmt.Sprintf(
		"repo_event_error: installation_id: %d, github_repo_id: %d, repo_name: %s, details: %s",
		e.InstallationID, e.GithubRepoID, e.RepoName, e.Details,
	)
}

func NewRepoNotFoundRepoEventError(installationID, githubRepoID db.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "repo_not_found",
	}
}

func NewMultipleReposFoundRepoEventError(installationID, githubRepoID db.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "multiple_repos_associated",
	}
}

func NewInactiveRepoRepoEventError(installationID, githubRepoID db.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "repo_not_active",
	}
}

func NewHasNoEarlyWarningRepoEventError(installationID, githubRepoID db.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "repo_has_no_early_warning",
	}
}
