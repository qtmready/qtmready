// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package shared

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	JWTClaims struct {
		UserID string `json:"user_id"`
		TeamID string `json:"team_id"`
		jwt.StandardClaims
	}
)

func GenerateAccessToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 15).Unix()
	if Service.Debug {
		expires = time.Now().Add(time.Hour * 24).Unix()
	}

	claims := &JWTClaims{
		UserID:         userID,
		TeamID:         teamID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expires, Issuer: Service.Name},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Service.Secret))
}

func GenerateRefreshToken(userID, teamID string) (string, error) {
	expires := time.Now().Add(time.Minute * 60).Unix()
	if Service.Debug {
		expires = time.Now().Add(time.Hour * 24 * 30).Unix()
	}

	claims := &JWTClaims{
		UserID:         userID,
		TeamID:         teamID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expires, Issuer: Service.Name},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Service.Secret))
}

func GetTeamIDFromContext(ctx echo.Context) string {
	claims := ctx.Get("user").(*jwt.Token).Claims.(*JWTClaims)
	return claims.TeamID
}

func GetUserIDFromContext(ctx echo.Context) string {
	claims := ctx.Get("user").(*jwt.Token).Claims.(*JWTClaims)
	return claims.UserID
}
