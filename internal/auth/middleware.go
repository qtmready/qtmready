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
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

const (
	BearerHeaderName    = "Authorization"
	BearerPrefix        = "Token"
	APIKeyHeaderName    = "X-API-KEY"
	GuardLookupTypeTeam = "team"
	GuardLookupTypeUser = "user"
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

	claims := &JWTClaims{
		UserID:           userID,
		TeamID:           teamID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expires), Issuer: shared.Service().GetName()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(shared.Service().GetSecret()))
}

// GenerateRefreshToken generates a long lived JWT token for the given user.
//
// TODO: Implement the logic for refreshing tokens using the refresh token.
func GenerateRefreshToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 60)
	if shared.Service().GetDebug() {
		expires = time.Now().Add(time.Hour * 24 * 30)
	}

	claims := &JWTClaims{
		UserID:           userID,
		TeamID:           teamID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expires), Issuer: shared.Service().GetName()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(shared.Service().GetSecret()))
}

// Middleware to provide JWT & API Key authentication.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// inspect.Context(ctx, true)
		// inspect.EchoHeaders(ctx)

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
					return echo.NewHTTPError(http.StatusBadRequest, ErrMissingAuthHeader)
				}
				// at this point, although the bearer is invalid, we know that endpoint can also be accessed with an API key
				// so we continue with the API key auth
				goto APIKeyAuth
			}

			parts := strings.Split(header, " ")

			if len(parts) != 2 || parts[0] != BearerPrefix {
				return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAuthHeader)
			}

			return bearerFn(next, ctx, parts[1])
		}

	APIKeyAuth:

		// do api key authentication
		if requiresKey && len(keyScopes) > -1 {
			shared.Logger().Debug("Authenticate with API Key")

			key := ctx.Request().Header.Get(APIKeyHeaderName)
			if key == "" {
				return echo.NewHTTPError(http.StatusBadRequest, ErrMissingAuthHeader)
			}

			return KeyFn(next, ctx, key)
		}

		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidAuthHeader)
	}
}

// bearerFn is the function that handles the JWT token authentication.
func bearerFn(next echo.HandlerFunc, ctx echo.Context, token string) error {
	parsed, err := jwt.ParseWithClaims(token, &JWTClaims{}, SecretFn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if claims, ok := parsed.Claims.(*JWTClaims); ok && parsed.Valid {
		ctx.Set("user_id", claims.UserID)
		ctx.Set("team_id", claims.TeamID)
	} else {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidOrExpiredToken)
	}

	return next(ctx)
}

// KeyFn validates the API key.
func KeyFn(next echo.HandlerFunc, ctx echo.Context, key string) error {
	guard := &Guard{}
	err := guard.VerifyAPIKey(key) // This will always return true if err is nil

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	switch guard.LookupType {
	case GuardLookupTypeTeam:
		ctx.Set("team_id", guard.LookupID.String())

	case GuardLookupTypeUser: // NOTE: this uses two db queries. we should optimize this. use k/v ?
		user := &User{}
		if err := db.Get(user, db.QueryParams{"id": guard.LookupID.String()}); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAuthHeader)
		}

		ctx.Set("user_id", user.ID.String()) // NOTE: IMHO, we shouldn't be converting to string here
		ctx.Set("team_id", user.TeamID.String())

	default:
		return echo.NewHTTPError(http.StatusUnauthorized, ErrMalformedAPIKey)
	}

	return next(ctx)
}

// SecretFn provides the secret for the JWT token.
func SecretFn(*jwt.Token) (any, error) {
	return []byte(shared.Service().GetSecret()), nil
}
