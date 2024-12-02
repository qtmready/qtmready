package defs

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
		RepositorySelection string              `json:"repository_selection"`
		RepositoriesAdded   []PartialRepository `json:"repositories_added"`
		RepositoriesRemoved []PartialRepository `json:"repositories_removed"`
	}

	WebhookRef struct {
		Ref          string         `json:"ref"`
		RefType      string         `json:"ref_type"`
		MasterBranch *string        `json:"master_branch"` // NOTE: This is only present in the create event.
		Description  *string        `json:"description"`   // NOTE: This is only present in the create event.
		PusherType   string         `json:"pusher_type"`
		Repository   Repository     `json:"repository"`
		Organization Organization   `json:"organization"`
		Sender       User           `json:"sender"`
		Installation InstallationID `json:"installation"`
		IsCreated    bool           `json:"is_created"`
	}
)

func (wr *WebhookRef) GetRef() string {
	return wr.Ref
}

func (wr *WebhookRef) GetRefType() string {
	return wr.RefType
}
