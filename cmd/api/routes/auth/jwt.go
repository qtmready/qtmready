package auth

import (
	"github.com/golang-jwt/jwt"
)

type (
	JWTClains struct {
		UserID string `json:"user_id"`
		TeamID string `json:"team_id"`
		jwt.StandardClaims
	}
)
