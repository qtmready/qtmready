// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package entities_test

import (
	"testing"

	"go.breu.io/ctrlplane/internal/entities"
)

func TestGithubInstallation(t *testing.T) {
	gi := &entities.GithubInstallation{}
	t.Run("GetTable", testEntityGetTable("github_installations", gi))
}

func TestGithubRepo(t *testing.T) {
	gi := &entities.GithubRepo{}
	t.Run("GetTable", testEntityGetTable("github_repos", gi))
}
