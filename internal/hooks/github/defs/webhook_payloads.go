package githubdefs

type (
	WebhookInstall struct {
		Action       string              `json:"action"`
		Installation Installation        `json:"installation"`
		Repositories []PartialRepository `json:"repositories"`
		Sender       User                `json:"sender"`
	}
)
