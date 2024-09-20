// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.


package github

import (
	"errors"
	"fmt"

	"go.breu.io/quantm/internal/shared"
)

type (
	MissingConfigurationError struct{}

	RepoEventError struct {
		InstallationID shared.Int64 `json:"installation_id"`
		GithubRepoID   shared.Int64 `json:"github_repo_id"`
		RepoName       string       `json:"repo_name"`
		Details        string       `json:"details"`
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

func NewRepoNotFoundRepoEventError(installationID, githubRepoID shared.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "repo_not_found",
	}
}

func NewMultipleReposFoundRepoEventError(installationID, githubRepoID shared.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "multiple_repos_associated",
	}
}

func NewInactiveRepoRepoEventError(installationID, githubRepoID shared.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "repo_not_active",
	}
}

func NewHasNoEarlyWarningRepoEventError(installationID, githubRepoID shared.Int64, repoName string) error {
	return &RepoEventError{
		InstallationID: installationID,
		GithubRepoID:   githubRepoID,
		RepoName:       repoName,
		Details:        "repo_has_no_early_warning",
	}
}
