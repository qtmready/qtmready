// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.step.sm/crypto/jose"
	"golang.org/x/crypto/hkdf"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

const (
	BearerHeaderName    = "Authorization"
	APIKeyHeaderName    = "X-API-KEY"
	GuardLookupTypeTeam = "team"
	GuardLookupTypeUser = "user"
)

var (
	BearerPrefixes = []string{"Token", "Bearer"}
)

// TODO: handle salt.
const (
	prefix = "Auth.js Generated Encryption Key"
	salt   = "authjs.session-token"
)

type (
	JWTClaims struct {
		UserID string `json:"user_id"`
		TeamID string `json:"team_id"`
		jwt.RegisteredClaims
	}
)

// GenerateAccessToken generates a short lived JWT token for the given user.
func GenerateAccessToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 15)
	if shared.Service().GetDebug() {
		expires = time.Now().Add(time.Hour * 24)
	}

	return generateJWE(userID, teamID, expires)
}

// GenerateRefreshToken generates a long lived JWT token for the given user.
//
// TODO: Implement the logic for refreshing tokens using the refresh token.
func GenerateRefreshToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 60)
	if shared.Service().GetDebug() {
		expires = time.Now().Add(time.Hour * 24 * 30)
	}

	return generateJWE(userID, teamID, expires)
}

// Middleware to provide JWT & API Key authentication.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		keyScopes, requiresKey := ctx.Get(APIKeyAuthScopes).([]string)
		bearerScopes, requiresBearer := ctx.Get(BearerAuthScopes).([]string)

		shared.Logger().Debug("requires bearer", "bearer", requiresBearer, "scopes", bearerScopes)
		shared.Logger().Debug("requires key", "key", requiresKey, "scopes", keyScopes)
		// if requiredKey and requiresBearer are both false, then we don't need to do any auth
		if !requiresKey && !requiresBearer {
			shared.Logger().Debug("no auth required")
			return next(ctx)
		}

		// do bearer authentication
		if requiresBearer && len(bearerScopes) > -1 {
			shared.Logger().Debug("Authenticate with Bearer Token")

			header := ctx.Request().Header.Get(BearerHeaderName)
			if header == "" {
				if !requiresKey {
					return shared.NewAPIError(http.StatusBadRequest, ErrMissingAuthHeader)
				}
				// at this point, although the bearer is invalid, we know that endpoint can also be accessed with an API key
				// so we continue with the API key auth
				goto APIKeyAuth
			}

			parts := strings.Split(header, " ")

			if len(parts) != 2 || !isValidPrefix(parts[0]) {
				return shared.NewAPIError(http.StatusBadRequest, ErrInvalidAuthHeader)
			}

			return bearerFn(next, ctx, parts[1])
		}

	APIKeyAuth:

		// do api key authentication
		if requiresKey && len(keyScopes) > -1 {
			shared.Logger().Debug("Authenticate with API Key")

			key := ctx.Request().Header.Get(APIKeyHeaderName)
			if key == "" {
				return shared.NewAPIError(http.StatusBadRequest, ErrMissingAuthHeader)
			}

			return KeyFn(next, ctx, key)
		}

		return shared.NewAPIError(http.StatusBadRequest, ErrInvalidAuthHeader)
	}
}

// bearerFn is the function that handles the JWT token authentication.
func bearerFn(next echo.HandlerFunc, ctx echo.Context, token string) error {
	enc, err := jose.Decrypt(
		[]byte(token),
		jose.WithAlg(string(jose.A256CBC_HS512)),
		jose.WithPassword(derive()),
	)

	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	var result map[string]any

	if err = json.Unmarshal(enc, &result); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	// Set user.id as user_id and team_id from user to the Echo context
	if info, ok := result["user"].(map[string]any); ok {
		ctx.Set("user_id", info["id"])
		ctx.Set("team_id", info["team_id"])
	}

	return next(ctx)
}

// KeyFn validates the API key.
func KeyFn(next echo.HandlerFunc, ctx echo.Context, key string) error {
	guard := &Guard{}
	err := guard.VerifyAPIKey(key) // This will always return true if err is nil

	if err != nil {
		return shared.NewAPIError(http.StatusUnauthorized, err)
	}

	switch guard.LookupType {
	case GuardLookupTypeTeam:
		ctx.Set("team_id", guard.LookupID.String())

	case GuardLookupTypeUser: // NOTE: this uses two db queries. we should optimize this. use k/v ?
		user := &User{}
		if err := db.Get(user, db.QueryParams{"id": guard.LookupID.String()}); err != nil {
			return shared.NewAPIError(http.StatusUnauthorized, ErrInvalidAuthHeader)
		}

		ctx.Set("user_id", user.ID.String()) // NOTE: IMHO, we shouldn't be converting to string here
		ctx.Set("team_id", user.TeamID.String())

	default:
		return shared.NewAPIError(http.StatusUnauthorized, ErrMalformedAPIKey)
	}

	return next(ctx)
}

// TODO: change to other info logger.
func info() string {
	return fmt.Sprintf("%s (%s)", prefix, salt)
}

func derive() []byte {
	kdf := hkdf.New(sha256.New, []byte(shared.Service().GetSecret()), []byte(salt), []byte(info()))
	key := make([]byte, 64)
	_, _ = io.ReadFull(kdf, key)

	return key
}

// generate generates a JWE token for the given user with specified expiration.
func generateJWE(userID, teamID string, expires time.Time) (string, error) {
	claims := map[string]any{
		"id":      userID,
		"team_id": teamID,
		"exp":     expires.Unix(),
		"iss":     shared.Service().GetName(),
	}

	result := map[string]any{
		"user": claims,
	}

	// Define JWT encode parameters
	params := JWTEncodeParams{
		Claims: result,
		Secret: derive(),
		MaxAge: time.Hour * 24,
		Salt:   nil,
	}

	// Encode JWT
	return EncodeJWT(params)
}

// Function to check if the prefix is valid.
func isValidPrefix(prefix string) bool {
	for _, valid := range BearerPrefixes {
		if prefix == valid {
			return true
		}
	}

	return false
}

// SecretFn provides the secret for the JWT token.
func SecretFn(*jwt.Token) (any, error) {
	return []byte(shared.Service().GetSecret()), nil
}
