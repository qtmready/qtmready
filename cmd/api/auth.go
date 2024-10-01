package main

import (
	"context"

	"go.breu.io/quantm/internal/auth"
)

// user_id decodes a JWE token and returns the subject claim.
//
// It calls `auth.DecodeJWE` to decrypt and validate the token. If the decoding is successful, it returns the subject
// claim. Otherwise, it returns an error.
func user_id(_ context.Context, token string) (string, error) {
	decoded, err := auth.DecodeJWE(token)
	if err != nil {
		return "", err
	}

	return decoded.Subject, nil
}
