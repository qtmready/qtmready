package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"go.step.sm/crypto/jose"
	"golang.org/x/crypto/hkdf"

	"go.breu.io/quantm/internal/shared"
)

const (
	alg    = "dir"
	enc    = "A256CBC-HS512"
	prefix = "Auth.js Generated Encryption Key"
	salt   = "__Secure-authjs.session-token"
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
		Key:       Derive(),
		KeyID:     base64.RawURLEncoding.EncodeToString(Derive()),
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

// DecodeJWE decodes and validates a JWE token.
// It returns the decoded claims if the token is valid.
func DecodeJWE(token string) (map[string]any, error) {
	enc, err := jose.Decrypt([]byte(token), jose.WithAlg(string(jose.A256CBC_HS512)), jose.WithPassword(Derive()))
	if err != nil {
		return nil, err
	}

	var claims map[string]any
	if err = json.Unmarshal(enc, &claims); err != nil {
		return nil, err
	}

	// Validate expiration
	if user, ok := claims["user"].(map[string]any); ok {
		if exp, ok := user["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, ErrTokenExpired
			}
		}
	}

	return claims, nil
}

func info() string {
	return fmt.Sprintf("%s (%s)", prefix, salt)
}

// Derive generates a derived key using HKDF.
func Derive() []byte {
	kdf := hkdf.New(sha256.New, []byte(shared.Service().GetSecret()), []byte(salt), []byte(info()))
	key := make([]byte, 64)
	_, _ = io.ReadFull(kdf, key)

	return key
}
