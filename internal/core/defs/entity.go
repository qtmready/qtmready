package defs

import (
	"encoding/json"

	"github.com/gocql/gocql"
)

func (repo *Repo) PreCreate() error { return nil }
func (repo *Repo) PreUpdate() error { return nil }

func (mp MessageProviderData) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(mp)
}

func (mp *MessageProviderData) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	if len(data) == 0 {
		*mp = MessageProviderData{}
		return nil
	}

	return json.Unmarshal(data, mp)
}

func (rp RepoProvider) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(info, rp.String())
}

func (rp *RepoProvider) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	*rp = RepoProviderMap[string(data)]

	return nil
}
