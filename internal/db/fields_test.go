package db

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
