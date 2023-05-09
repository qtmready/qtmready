package core

import "context"

var (
	Core = &core{}
)

type (
	core struct {
		Activity     Activities
		Workflow     Workflows
		ProvidersMap providersMap
	}
)

type (
	Provider interface {
		GetLatestCommitforRepo(ctx context.Context, providerID string, branch string) (*string, error)
	}
	providersMap map[RepoProvider]Provider
)

func (c *core) Init() {
	c.ProvidersMap = make(providersMap)
}
