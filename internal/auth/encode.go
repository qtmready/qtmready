package auth

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"go.step.sm/crypto/jose"
)

const (
	alg = "dir"
	enc = "A256CBC-HS512"
)

// JWTEncodeParams represents the parameters for JWT encoding.
type (
	JWTEncodeParams struct {
		Claims map[string]any // Payload of the JWT
		Secret []byte         // Encryption key
		MaxAge time.Duration
		Salt   []byte // Salt used for key derivation
	}
)

// EncodeJWT encodes a JWT.
func EncodeJWT(params JWTEncodeParams) (string, error) {
	// Generate a JWE key
	key := jose.JSONWebKey{
		Key:       params.Secret,
		KeyID:     base64.RawURLEncoding.EncodeToString(params.Secret),
		Algorithm: alg,
		Use:       "enc",
	}

	// Create a new encrypter
	encrypter, err := jose.NewEncrypter(jose.A256CBC_HS512, jose.Recipient{Algorithm: alg, Key: key}, nil)
	if err != nil {
		return "", err
	}

	// Marshal the payload to JSON
	bytes, err := json.Marshal(params.Claims)
	if err != nil {
		return "", err
	}

	// Encrypt the payload
	encrypted, err := encrypter.Encrypt(bytes)
	if err != nil {
		return "", err
	}

	// Serialize JWE token
	serialized, err := encrypted.CompactSerialize()
	if err != nil {
		return "", err
	}

	return serialized, nil
}
