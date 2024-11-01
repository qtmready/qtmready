package githubwfs

import (
	"go.temporal.io/sdk/workflow"
)

type (
	StatusInstall struct {
		Webhook bool
		Request bool
	}
)

// Install installs the Github Integration.
func Install(ctx workflow.Context) error {
	log := workflow.GetLogger(ctx)

	log.Info(pfx_install("installing..."))

	return nil
}
