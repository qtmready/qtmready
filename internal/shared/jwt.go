package shared

import (
	"github.com/golang-jwt/jwt"
)

type (
	JWTClaims struct {
		UserID string `json:"user_id"`
		TeamID string `json:"team_id"`
		jwt.StandardClaims
	}
)
