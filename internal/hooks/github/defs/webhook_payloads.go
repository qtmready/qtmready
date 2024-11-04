package githubdefs

type (
	WebhookInstall struct {
		Action       string              `json:"action"`
		Installation Installation        `json:"installation"`
		Repositories []PartialRepository `json:"repositories"`
		Sender       User                `json:"sender"`
	}

	WebhookInstallRepos struct {
		Action              string              `json:"action"`
		Installation        Installation        `json:"installation"`
		RepositoriesAdeed   []PartialRepository `json:"repositories_added"`
		RepositoriesRemoved []PartialRepository `json:"repositories_removed"`
		RepositorySelection string              `json:"repository_selection"`
	}
)
