package core

import "context"

var (
	Core = &core{}
)

type (
	core struct {
		Activity  Activities
		Workflow  Workflows
		Providers providers
	}
)

type (
	Provider interface {
		GetLatestCommitforRepo(ctx context.Context, providerID string, branch string) (string, error)
	}
	providers map[RepoProvider]Provider
)

func (c *core) Init() {
	c.Providers = make(providers)
}
