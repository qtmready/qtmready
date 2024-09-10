// Copyright © 2023, Breu, Inc. <info@breu.io>
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

package auth_test

import (
	"testing"

	"github.com/gosimple/slug"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/testutils"
)

func TestTeam(t *testing.T) {
	team := &auth.Team{
		Name: "Team Name",
	}
	_ = team.PreCreate()

	opsTests := testutils.TestFnMap{
		"Slug": testutils.TestFn{Args: team, Want: nil, Run: testTeamSlug},
	}

	t.Run("GetTable", testutils.TestEntityGetTable("teams", team))
	t.Run("EntityOps", testutils.TestEntityOps(team, opsTests))
}

func testTeamSlug(args, want any) func(*testing.T) {
	team := args.(*auth.Team)
	sluglen := len(slug.Make(team.Name)) + 1 + 4

	return func(t *testing.T) {
		if len(team.Slug) != sluglen {
			t.Errorf("slug length is not correct, got: %d, want: %d", len(team.Slug), sluglen)
		}
	}
}
