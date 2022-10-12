package entities_test

import (
	"testing"

	"go.breu.io/ctrlplane/internal/entities"
)

func TestApp(t *testing.T) {
	app := &entities.App{}
	t.Run("GetTable", testEntityGetTable("apps", app))
}

func TestRepo(t *testing.T) {
	repo := &entities.Repo{}
	t.Run("GetTable", testEntityGetTable("repos", repo))
}

func TestWorkload(t *testing.T) {
	workload := &entities.Workload{}
	t.Run("GetTable", testEntityGetTable("workloads", workload))
}

func TestResource(t *testing.T) {
	resource := &entities.Resource{}
	t.Run("GetTable", testEntityGetTable("resources", resource))
}

func TestBlueprint(t *testing.T) {
	blueprint := &entities.Blueprint{}
	t.Run("GetTable", testEntityGetTable("blueprints", blueprint))
}

func TestRollout(t *testing.T) {
	rollout := &entities.Rollout{}
	t.Run("GetTable", testEntityGetTable("rollouts", rollout))
}
