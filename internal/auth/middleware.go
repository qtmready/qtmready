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

	"github.com/labstack/echo/v4"

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

// Middleware provides JWE & API Key authentication.
// It checks for bearer token (JWE) and API key based on the context requirements.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		keyScopes, requiresKey := ctx.Get(APIKeyAuthScopes).([]string)
		bearerScopes, requiresBearer := ctx.Get(BearerAuthScopes).([]string)

		shared.Logger().Debug("auth requirements", "bearer", requiresBearer, "key", requiresKey)

		if !requiresKey && !requiresBearer {
			return next(ctx)
		}

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
func GenerateAccessToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 15)
	if shared.Service().GetDebug() {
		expires = time.Now().Add(time.Hour * 24)
	}

	return generateJWE(userID, teamID, expires)
}

// GenerateRefreshToken creates a long-lived JWE token for the given user.
func GenerateRefreshToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 60)
	if shared.Service().GetDebug() {
		expires = time.Now().Add(time.Hour * 24 * 30)
	}

	return generateJWE(userID, teamID, expires)
}

// BearerFn handles the JWE token authentication.
// It decrypts the JWE token, validates its contents, and sets user and team IDs in the context.
func BearerFn(next echo.HandlerFunc, ctx echo.Context, token string) error {
	claims, err := DecodeJWE(token)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if info, ok := claims["user"].(map[string]any); ok {
		ctx.Set("user_id", info["id"])
		ctx.Set("team_id", info["team_id"])
	}

	if next != nil {
		return next(ctx)
	}

	return nil
}

// generateJWE creates a JWE token with the given user ID, team ID, and expiration time.
func generateJWE(userID, teamID string, expires time.Time) (string, error) {
	claims := map[string]any{
		"user": map[string]any{
			"id":      userID,
			"team_id": teamID,
			"exp":     expires.Unix(),
			"iss":     shared.Service().GetName(),
		},
	}

	return EncodeJWT(JWTEncodeParams{
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
// It verifies the API key, determines the lookup type (team or user),
// and sets the appropriate IDs in the context.
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
