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

package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/shared"
)

type (
	CryptoTestSuite struct {
		suite.Suite
	}
)

func (s *CryptoTestSuite) SetupSuite() {
	shared.InitServiceForTest()
}

// Positive Tests

func (s *CryptoTestSuite) TestEncodeJWE_Positive() {
	claims := map[string]any{
		"user": map[string]any{
			"id":   "123",
			"name": "Test User",
		},
	}

	params := auth.JWTEncodeParams{
		Claims: claims,
		MaxAge: time.Hour,
	}

	token, err := auth.EncodeJWT(params)
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), token)
}

func (s *CryptoTestSuite) TestDecodeJWE_Positive() {
	claims := map[string]any{
		"user": map[string]any{
			"id":   "123",
			"name": "Test User",
			"exp":  float64(time.Now().Add(time.Hour).Unix()),
		},
	}

	params := auth.JWTEncodeParams{
		Claims: claims,
		MaxAge: time.Hour,
	}

	token, err := auth.EncodeJWT(params)
	require.NoError(s.T(), err)

	decodedClaims, err := auth.DecodeJWE(token)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), claims["user"].(map[string]any)["id"], decodedClaims["user"].(map[string]any)["id"])
	assert.Equal(s.T(), claims["user"].(map[string]any)["name"], decodedClaims["user"].(map[string]any)["name"])
}

func (s *CryptoTestSuite) TestEncodeJWE_EmptyClaims() {
	params := auth.JWTEncodeParams{
		Claims: map[string]any{},
		MaxAge: time.Hour,
	}

	token, err := auth.EncodeJWT(params)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), token)

	// Optionally, you can try to decode the token to ensure it's valid
	decodedClaims, err := auth.DecodeJWE(token)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), decodedClaims)
}

// Negative Tests

func (s *CryptoTestSuite) TestDecodeJWE_Negative_InvalidToken() {
	_, err := auth.DecodeJWE("invalid.token")
	assert.Error(s.T(), err)
}

func (s *CryptoTestSuite) TestDecodeJWE_Negative_ExpiredToken() {
	claims := map[string]any{
		"user": map[string]any{
			"id":   "123",
			"name": "Test User",
			"exp":  float64(time.Now().Add(-time.Hour).Unix()), // Expired token
		},
	}

	params := auth.JWTEncodeParams{
		Claims: claims,
		MaxAge: time.Hour,
	}

	token, err := auth.EncodeJWT(params)
	require.NoError(s.T(), err)

	_, err = auth.DecodeJWE(token)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), auth.ErrTokenExpired, err)
}

// Smoke Tests

func (s *CryptoTestSuite) TestEncodeDecodeJWE_Smoke() {
	// Test a variety of claim types
	claims := map[string]any{
		"string": "test",
		"int":    123,
		"float":  3.14,
		"bool":   true,
		"array":  []string{"a", "b", "c"},
		"map": map[string]any{
			"nested": "value",
		},
		"user": map[string]any{
			"id":   "456",
			"name": "Jane Doe",
			"exp":  float64(time.Now().Add(time.Hour).Unix()),
		},
	}

	params := auth.JWTEncodeParams{
		Claims: claims,
		MaxAge: time.Hour,
	}

	token, err := auth.EncodeJWT(params)
	require.NoError(s.T(), err)

	decodedClaims, err := auth.DecodeJWE(token)
	require.NoError(s.T(), err)

	// Compare each claim individually, allowing for type differences
	for key, expectedValue := range claims {
		decodedValue, exists := decodedClaims[key]
		assert.True(s.T(), exists, "Key %s not found in decoded claims", key)

		switch expected := expectedValue.(type) {
		case []string:
			// Convert expected []string to []interface{} for comparison
			expectedInterface := make([]any, len(expected))
			for i, v := range expected {
				expectedInterface[i] = v
			}

			assert.Equal(s.T(), expectedInterface, decodedValue, "Mismatch for key: %s", key)
		case int:
			// Allow for int to be decoded as float64
			assert.Equal(s.T(), float64(expected), decodedValue, "Mismatch for key: %s", key)
		default:
			assert.Equal(s.T(), expected, decodedValue, "Mismatch for key: %s", key)
		}
	}
}

func (s *CryptoTestSuite) TestDerive_Smoke() {
	// Test multiple derivations
	key1 := auth.Derive()
	key2 := auth.Derive()

	assert.NotEmpty(s.T(), key1)
	assert.NotEmpty(s.T(), key2)
	assert.Len(s.T(), key1, 64)
	assert.Len(s.T(), key2, 64)
	assert.Equal(s.T(), key1, key2) // Keys should be the same for the same secret
}

// TestEncode runs the test suite.
func TestEncode(t *testing.T) {
	suite.Run(t, new(CryptoTestSuite))
}
