package core

import (
	"encoding/json"

	"github.com/gocql/gocql"

	"go.breu.io/ctrlplane/internal/db"
)

func (stack *Stack) PreCreate() error { stack.Slug = db.CreateSlug(stack.Name); return nil }
func (stack *Stack) PreUpdate() error { return nil }

func (repo *Repo) PreCreate() error { return nil }
func (repo *Repo) PreUpdate() error { return nil }

func (config StackConfig) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(config)
}

func (config *StackConfig) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, config)
}
