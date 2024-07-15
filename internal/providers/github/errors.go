// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
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
