package auth_test

import (
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/gocql/gocql"
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

func tolerate(t1, t2 time.Time, tolerance time.Duration) bool {
	return t1.Before(t2.Add(tolerance)) && t2.Before(t1.Add(tolerance))
}

// Positive Tests

func (s *CryptoTestSuite) TestEncodeJWE_Positive() {
	user := auth.User{
		CreatedAt:  time.Now(),
		Email:      "alice@example.com",
		FirstName:  "Alice",
		ID:         gocql.MustRandomUUID(),
		IsActive:   true,
		IsVerified: true,
		LastName:   "Smith",
		TeamID:     gocql.MustRandomUUID(),
		UpdatedAt:  time.Now(),
	}
	claims := auth.Claims{
		Claims: jwt.Claims{
			Issuer:   "test",
			Subject:  "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			Audience: []string{"test"},
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:       "f7654321-0987-6543-2109-876543210987",
		},
		User: user,
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
	user := auth.User{
		CreatedAt:  time.Now(),
		Email:      "alice@example.com",
		FirstName:  "Alice",
		ID:         gocql.MustRandomUUID(),
		IsActive:   true,
		IsVerified: true,
		LastName:   "Smith",
		TeamID:     gocql.MustRandomUUID(),
		UpdatedAt:  time.Now(),
	}
	claims := auth.Claims{
		Claims: jwt.Claims{
			Issuer:   "test",
			Subject:  "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			Audience: []string{"test"},
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:       "f7654321-0987-6543-2109-876543210987",
		},
		User: user,
	}

	params := auth.JWTEncodeParams{
		Claims: claims,
		MaxAge: time.Hour,
	}

	token, err := auth.EncodeJWT(params)
	require.NoError(s.T(), err)

	decodedClaims, err := auth.DecodeJWE(token)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), claims.Issuer, decodedClaims.Issuer)
	assert.Equal(s.T(), claims.Subject, decodedClaims.Subject)
	assert.Equal(s.T(), claims.Audience, decodedClaims.Audience)
	assert.Equal(s.T(), claims.ID, decodedClaims.ID)
	assert.Equal(s.T(), claims.IssuedAt, decodedClaims.IssuedAt)
	assert.Equal(s.T(), claims.Expiry, decodedClaims.Expiry)
	assert.True(s.T(), tolerate(claims.IssuedAt.Time(), decodedClaims.IssuedAt.Time(), 1*time.Second))
	assert.Equal(s.T(), claims.User.Email, decodedClaims.User.Email)
	assert.Equal(s.T(), claims.User.FirstName, decodedClaims.User.FirstName)
	assert.Equal(s.T(), claims.User.ID, decodedClaims.User.ID)
	assert.Equal(s.T(), claims.User.IsActive, decodedClaims.User.IsActive)
	assert.Equal(s.T(), claims.User.IsVerified, decodedClaims.User.IsVerified)
	assert.Equal(s.T(), claims.User.LastName, decodedClaims.User.LastName)
	assert.Equal(s.T(), claims.User.TeamID, decodedClaims.User.TeamID)
	assert.True(s.T(), tolerate(claims.User.CreatedAt, decodedClaims.User.CreatedAt, 1*time.Second))
}

func (s *CryptoTestSuite) TestEncodeJWE_EmptyClaims() {
	params := auth.JWTEncodeParams{
		Claims: auth.Claims{
			Claims: jwt.Claims{
				Issuer:   "test",
				Subject:  "a1b2c3d4-e5f6-7890-1234-567890abcdef",
				Audience: []string{"test"},
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
				ID:       "f7654321-0987-6543-2109-876543210987",
			},
			User: auth.User{},
		},
		MaxAge: time.Hour,
	}

	token, err := auth.EncodeJWT(params)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), token)

	// Optionally, you can try to decode the token to ensure it's valid
	decodedClaims, err := auth.DecodeJWE(token)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), decodedClaims.User)
}

// Negative Tests

func (s *CryptoTestSuite) TestDecodeJWE_Negative_InvalidToken() {
	_, err := auth.DecodeJWE("invalid.token")
	assert.Error(s.T(), err)
}

func (s *CryptoTestSuite) TestDecodeJWE_Negative_ExpiredToken() {
	user := auth.User{
		CreatedAt:  time.Now(),
		Email:      "alice@example.com",
		FirstName:  "Alice",
		ID:         gocql.MustRandomUUID(),
		IsActive:   true,
		IsVerified: true,
		LastName:   "Smith",
		TeamID:     gocql.MustRandomUUID(),
		UpdatedAt:  time.Now(),
	}
	claims := auth.Claims{
		Claims: jwt.Claims{
			Issuer:   "test",
			Subject:  "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			Audience: []string{"test"},
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Expiry:   jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired token
			ID:       "f7654321-0987-6543-2109-876543210987",
		},
		User: user,
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
	user := auth.User{
		CreatedAt:  time.Now(),
		Email:      "alice@example.com",
		FirstName:  "Alice",
		ID:         gocql.MustRandomUUID(),
		IsActive:   true,
		IsVerified: true,
		LastName:   "Smith",
		TeamID:     gocql.MustRandomUUID(),
		UpdatedAt:  time.Now(),
	}
	claims := auth.Claims{
		Claims: jwt.Claims{
			Issuer:   "test",
			Subject:  "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			Audience: []string{"test"},
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:       "f7654321-0987-6543-2109-876543210987",
		},
		User: user,
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
	assert.Equal(s.T(), claims.Issuer, decodedClaims.Issuer)
	assert.Equal(s.T(), claims.Subject, decodedClaims.Subject)
	assert.Equal(s.T(), claims.Audience, decodedClaims.Audience)
	assert.Equal(s.T(), claims.ID, decodedClaims.ID)
	assert.Equal(s.T(), claims.IssuedAt, decodedClaims.IssuedAt)
	assert.Equal(s.T(), claims.Expiry, decodedClaims.Expiry)
	assert.True(s.T(), tolerate(claims.User.CreatedAt, decodedClaims.User.CreatedAt, 1*time.Second))
	assert.Equal(s.T(), claims.User.Email, decodedClaims.User.Email)
	assert.Equal(s.T(), claims.User.FirstName, decodedClaims.User.FirstName)
	assert.Equal(s.T(), claims.User.ID, decodedClaims.User.ID)
	assert.Equal(s.T(), claims.User.IsActive, decodedClaims.User.IsActive)
	assert.Equal(s.T(), claims.User.IsVerified, decodedClaims.User.IsVerified)
	assert.Equal(s.T(), claims.User.LastName, decodedClaims.User.LastName)
	assert.Equal(s.T(), claims.User.TeamID, decodedClaims.User.TeamID)
	assert.True(s.T(), tolerate(claims.User.UpdatedAt, decodedClaims.User.UpdatedAt, 1*time.Second))
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

func TestEncodeDecodeJWE_UserObject(t *testing.T) {
	// Create a user object
	user := auth.User{
		CreatedAt:  time.Now(),
		Email:      "test@example.com",
		FirstName:  "Test",
		ID:         gocql.MustRandomUUID(),
		IsActive:   true,
		IsVerified: true,
		LastName:   "User",
		TeamID:     gocql.MustRandomUUID(),
		UpdatedAt:  time.Now(),
	}

	// Create claims with the user object
	claims := auth.Claims{
		Claims: jwt.Claims{
			Issuer:   "test",
			Subject:  "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			Audience: []string{"test"},
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:       "f7654321-0987-6543-2109-876543210987",
		},
		User: user,
	}

	// Encode the JWT
	token, err := auth.EncodeJWT(auth.JWTEncodeParams{
		Claims: claims,
		MaxAge: time.Hour,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Decode the JWT
	decodedClaims, err := auth.DecodeJWE(token)
	require.NoError(t, err)

	// Check if the decoded user object matches the original user object
	assert.True(t, tolerate(user.CreatedAt, decodedClaims.User.CreatedAt, 1*time.Second))
	assert.Equal(t, user.Email, decodedClaims.User.Email)
	assert.Equal(t, user.FirstName, decodedClaims.User.FirstName)
	assert.Equal(t, user.ID, decodedClaims.User.ID)
	assert.Equal(t, user.IsActive, decodedClaims.User.IsActive)
	assert.Equal(t, user.IsVerified, decodedClaims.User.IsVerified)
	assert.Equal(t, user.LastName, decodedClaims.User.LastName)
	assert.Equal(t, user.TeamID, decodedClaims.User.TeamID)
	assert.True(t, tolerate(user.UpdatedAt, decodedClaims.User.UpdatedAt, 1*time.Second))
}

// TestCrypto runs the test suite.
func TestCrypto(t *testing.T) {
	suite.Run(t, new(CryptoTestSuite))
}
