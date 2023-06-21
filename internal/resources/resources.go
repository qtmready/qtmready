package resources

import "go.breu.io/ctrlplane/internal/core"

type (
	Deployable struct {
		Resource resource
		Workload workload
	}

	workload struct {
		image string
	}
	resource interface {
		Deploy(d *Deployable)
		Parse(assets *core.Assets, yaml *string)
	}
)
