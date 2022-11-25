// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

const (
	JwtPrefix       = "Token"
	APIKeyPrefix    = "API-KEY"
	GuardLookupTeam = "team"
	GuardLookupUser = "user"
)

type (
	JWTClaims struct {
		UserID string `json:"user_id"`
		TeamID string `json:"team_id"`
		jwt.StandardClaims
	}
)

var (
	ErrMalformedToken        = errors.New("malformed token")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
	ErrMissingAuthHeader     = errors.New("no authorization header provided")
	ErrInvalidAuthHeader     = errors.New("invalid authorization header")
)

// GenerateAccessToken generates a short lived JWT token for the given user.
func GenerateAccessToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 15).Unix()
	if shared.Service.Debug {
		expires = time.Now().Add(time.Hour * 24).Unix()
	}

	claims := &JWTClaims{
		UserID:         userID,
		TeamID:         teamID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expires, Issuer: shared.Service.Name},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(shared.Service.Secret))
}

// GenerateRefreshToken generates a long lived JWT token for the given user.
//
// TODO: Implement the logic for refreshing tokens using the refresh token.
func GenerateRefreshToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 60).Unix()
	if shared.Service.Debug {
		expires = time.Now().Add(time.Hour * 24 * 30).Unix()
	}

	claims := &JWTClaims{
		UserID:         userID,
		TeamID:         teamID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expires, Issuer: shared.Service.Name},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(shared.Service.Secret))
}

// Middleware to provide JWT & API Key authentication.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// get the authorization header
		header := ctx.Request().Header.Get("Authorization")
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, ErrMissingAuthHeader)
		}

		// split the header at the space to get the scheme and token
		fields := strings.Split(header, " ")
		if len(fields) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAuthHeader)
		}

		// apply the correct validation function based on the scheme
		schema, secret := fields[0], fields[1]

		switch schema {
		case JwtPrefix:
			if err := validateToken(ctx, secret); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}
		case APIKeyPrefix:
			if err := validateKey(ctx, secret); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}
		default:
			return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAuthHeader)
		}

		return next(ctx)
	}
}

// validateToken validates the JWT token.
func validateToken(ctx echo.Context, token string) error {
	parsed, err := jwt.ParseWithClaims(token, &JWTClaims{}, getSecret)
	if err != nil {
		return ErrMalformedToken
	}

	if claims, ok := parsed.Claims.(*JWTClaims); ok && parsed.Valid {
		ctx.Set("user_id", claims.UserID)
		ctx.Set("team_id", claims.TeamID)
	} else {
		return ErrInvalidOrExpiredToken
	}

	return nil
}

// validateKey validates the API key.
func validateKey(ctx echo.Context, key string) error {
	guard := &entities.Guard{}
	err := guard.VerifyAPIKey(key) // This will always return true if err is nil

	if err != nil {
		return err
	}

	switch guard.LookupType {
	case GuardLookupTeam:
		ctx.Set("team_id", guard.LookupID.String())
	case GuardLookupUser: // NOTE: this uses two db queries. we should optimize this. use k/v ?
		user := &entities.User{}
		if err := db.Get(user, db.QueryParams{"id": guard.LookupID.String()}); err != nil {
			return err
		}

		ctx.Set("user_id", user.ID.String()) // NOTE: IMHO, we shouldn't be converting to string here
		ctx.Set("team_id", user.TeamID.String())
	}

	return nil
}

// getSecret provides the secret for the JWT token.
func getSecret(t *jwt.Token) (interface{}, error) {
	return []byte(shared.Service.Secret), nil
}
