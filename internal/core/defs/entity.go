// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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
