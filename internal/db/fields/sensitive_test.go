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

// Copyright © 2024, Breu, Inc. <info@breu.io>
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

package fields

import (
	"encoding/json"
	"testing"

	"github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/suite"

	"go.breu.io/quantm/internal/shared"
)

type (
	EncryptedFieldTestSuite struct {
		suite.Suite

		sensitive Sensitive
	}
)

func (s *EncryptedFieldTestSuite) SetupSuite() {
	shared.InitServiceForTest() // instantiate the service with a random secret

	sensitive := password.MustGenerate(32, 8, 8, false, false)
	s.sensitive = Sensitive(sensitive)
}

func TestEncryptedField(t *testing.T) {
	suite.Run(t, new(EncryptedFieldTestSuite))
}

func (s *EncryptedFieldTestSuite) TestEncryptDecrypt() {
	// Encrypt the string
	encrypted, err := s.sensitive.encrypt()
	s.NoError(err)

	// Assert that the encrypted value is not equal to the original string
	s.NotEqual(s.sensitive.String(), string(encrypted))

	// Decrypt the string
	var decrypted Sensitive
	err = decrypted.from(encrypted)
	s.NoError(err)

	s.Equal(s.sensitive.String(), decrypted.String())
}

func (s *EncryptedFieldTestSuite) TestMarshalJSON() {
	// Marshal the string to JSON
	data, err := json.Marshal(s.sensitive)
	s.NoError(err)

	// Unmarshal the JSON data
	var decrypted Sensitive
	err = json.Unmarshal(data, &decrypted)
	s.NoError(err)

	s.Equal(s.sensitive.String(), decrypted.String())
}

func (s *EncryptedFieldTestSuite) TestMarshalCQL() {
	// Marshal the string to CQL
	cql, err := s.sensitive.MarshalCQL()
	s.NoError(err)

	// Unmarshal the CQL data
	var decrypted Sensitive
	err = decrypted.UnmarshalCQL(cql)
	s.NoError(err)

	s.Equal(s.sensitive.String(), decrypted.String())
}
