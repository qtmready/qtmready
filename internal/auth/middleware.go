// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
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

package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

const (
	BearerHeaderName    = "Authorization" // Header name for bearer token authentication.
	APIKeyHeaderName    = "X-API-KEY"     // Header name for API key authentication.
	GuardLookupTypeTeam = "team"          // Lookup type for API key authentication: team.
	GuardLookupTypeUser = "user"          // Lookup type for API key authentication: user.
)

var (
	BearerPrefixes = []string{"Token", "Bearer"} // Valid prefixes for bearer tokens.
)

// Middleware provides JWE & API Key authentication.
//
// It checks for bearer token (JWE) and API key based on the context requirements. If both are required, it prioritizes
// bearer token authentication. If neither is required, it proceeds to the next handler.
//
// The context requirements are determined by the `APIKeyAuthScopes` and `BearerAuthScopes` values, set by middleware
// handlers upstream.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		keyScopes, requiresKey := ctx.Get(APIKeyAuthScopes).([]string)
		bearerScopes, requiresBearer := ctx.Get(BearerAuthScopes).([]string)

		shared.Logger().Debug("auth requirements", "bearer", requiresBearer, "key", requiresKey)

		if !requiresKey && !requiresBearer {
			return next(ctx)
		}

		// Prioritize bearer token authentication if both are required.
		if requiresBearer && len(bearerScopes) > -1 {
			header := ctx.Request().Header.Get(BearerHeaderName)
			if header != "" {
				parts := strings.Split(header, " ")
				if len(parts) == 2 && isValidPrefix(parts[0]) {
					if err := BearerFn(nil, ctx, parts[1]); err == nil {
						return next(ctx)
					}
				}
			}
		}

		if requiresKey && len(keyScopes) > -1 {
			key := ctx.Request().Header.Get(APIKeyHeaderName)
			if key != "" {
				if err := KeyFn(nil, ctx, key); err == nil {
					return next(ctx)
				}
			}
		}

		return shared.NewAPIError(http.StatusUnauthorized, ErrInvalidAuthHeader)
	}
}

// GenerateAccessToken creates a short-lived JWE token for the given user.
//
// The token's expiration time is set to 15 minutes by default, but it can be overridden to 24 hours in debug mode.
func GenerateAccessToken(user *User) (string, error) {
	expires := time.Now().Add(time.Minute * 15)
	if shared.Service().GetDebug() {
		expires = time.Now().Add(time.Hour * 24)
	}

	return generateJWE(user, expires)
}

// GenerateRefreshToken creates a long-lived JWE token for the given user.
//
// The token's expiration time is set to 60 minutes by default, but it can be overridden to 30 days in debug mode.
func GenerateRefreshToken(user *User) (string, error) {
	expires := time.Now().Add(time.Minute * 60)
	if shared.Service().GetDebug() {
		expires = time.Now().Add(time.Hour * 24 * 30)
	}

	return generateJWE(user, expires)
}

// BearerFn handles the JWE token authentication.
//
// It decrypts the JWE token, validates its contents, and sets user and team IDs in the context.
func BearerFn(next echo.HandlerFunc, ctx echo.Context, token string) error {
	claims, err := DecodeJWE(token)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	ctx.Set("user_id", claims.User.ID.String())
	ctx.Set("team_id", claims.User.TeamID.String())

	if next != nil {
		return next(ctx)
	}

	return nil
}

// generateJWE creates a JWE token with the given user ID, team ID, and expiration time.
func generateJWE(user *User, expires time.Time) (string, error) {
	claims := Claims{
		Claims: jwt.Claims{
			Issuer:   shared.Service().GetName(),
			Subject:  user.ID.String(),
			Audience: jwt.Audience{shared.Service().GetName()},
			Expiry:   jwt.NewNumericDate(expires),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		User: *user,
	}

	return EncodeJWE(JWTEncodeParams{
		Claims: claims,
		Secret: Derive(),
		MaxAge: time.Hour * 24,
	})
}

// isValidPrefix checks if the given prefix is a valid bearer token prefix.
func isValidPrefix(prefix string) bool {
	for _, valid := range BearerPrefixes {
		if prefix == valid {
			return true
		}
	}

	return false
}

// KeyFn handles the API key authentication.
//
// It verifies the API key, determines the lookup type (team or user), and sets the appropriate IDs in the context.
func KeyFn(next echo.HandlerFunc, ctx echo.Context, key string) error {
	guard := &Guard{}
	if err := guard.VerifyAPIKey(key); err != nil {
		return shared.NewAPIError(http.StatusUnauthorized, err)
	}

	switch guard.LookupType {
	case GuardLookupTypeTeam:
		ctx.Set("team_id", guard.LookupID.String())
	case GuardLookupTypeUser:
		user := &User{}
		if err := db.Get(user, db.QueryParams{"id": guard.LookupID.String()}); err != nil {
			return shared.NewAPIError(http.StatusUnauthorized, ErrInvalidAuthHeader)
		}

		ctx.Set("user_id", user.ID.String())
		ctx.Set("team_id", user.TeamID.String())
	default:
		return shared.NewAPIError(http.StatusUnauthorized, ErrMalformedAPIKey)
	}

	if next != nil {
		return next(ctx)
	}

	return nil
}
