package shared

import (
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

func GetTeamIDFromContext(ctx echo.Context) string {
	claims := ctx.Get("user").(*jwt.Token).Claims.(*JWTClaims)
	return claims.TeamID
}

func GetUserIDFromContext(ctx echo.Context) string {
	claims := ctx.Get("user").(*jwt.Token).Claims.(*JWTClaims)
	return claims.UserID
}
