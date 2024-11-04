package githubdefs

import (
	"go.breu.io/durex/queues"
)

const (
	SignalRequestInstall queues.Signal = "install_from_request"
	SignalWebhookInstall queues.Signal = "install_from_webhook"
	SignalWebhookPush    queues.Signal = "push" // TODO - need to refine
)
