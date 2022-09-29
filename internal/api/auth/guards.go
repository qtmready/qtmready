// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

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
	ErrInvalidAPIKey         = errors.New("invalid API key")
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
// TODO: this is not implemented yet.
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
		switch fields[0] {
		case JwtPrefix:
			if err := validateToken(ctx, fields[1]); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}
		case APIKeyPrefix:
			if err := validateKey(ctx, fields[1]); err != nil {
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
	// This is where the magic happens
	ok, err := guard.VerifyAPIKey(key)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidAPIKey
	}

	switch guard.LookupType {
	case GuardLookupTeam:
		ctx.Set("team_id", guard.LookupID.String())
	case GuardLookupUser: // FIXME: we have two db calls here, we should be able to do this in one
		user := &entities.User{}
		if err := db.Get(user, db.QueryParams{"id": guard.LookupID.String()}); err != nil {
			return ErrInvalidAPIKey
		}

		ctx.Set("user_id", user.ID.String())
		ctx.Set("team_id", user.TeamID.String())
	}

	return nil
}

// getSecret provides the secret for the JWT token.
func getSecret(t *jwt.Token) (interface{}, error) {
	return []byte(shared.Service.Secret), nil
}
