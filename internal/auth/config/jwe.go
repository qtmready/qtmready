package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-jose/go-jose/v4/jwt"
	"go.step.sm/crypto/jose"
	"golang.org/x/crypto/hkdf"

	"go.breu.io/quantm/internal/erratic"
)

const (
	alg    = "dir"                              // Algorithm used for encryption.
	enc    = "A256CBC-HS512"                    // Encryption method.
	prefix = "Auth.js Generated Encryption Key" // Key prefix.
	salt   = "__Secure-authjs.session-token"    // Salt for key derivation.
)

type (
	// Claims represents the payload of the JWT token.
	Claims struct {
		jwt.Claims        // Standard JWT claims.
		UserID     string `json:"user_id"` // User ID.
		OrgID      string `json:"org_id"`  // Organization ID.
	}

	// JWTEncodeParams contains the parameters for JWT encoding.
	JWTEncodeParams struct {
		Claims Claims        // Payload of the JWT.
		Secret []byte        // Encryption key.
		MaxAge time.Duration // Maximum age of the token.
		Salt   []byte        // Salt used for key derivation.
	}
)

// info returns a string containing the key prefix and salt for key derivation.
func info() string {
	return fmt.Sprintf("%s (%s)", prefix, salt)
}

// EncodeJWE encodes a JWT.
//
// It generates a JWE key using the `Derive` function, creates an encrypter, marshals the payload to JSON, encrypts
// the payload, serializes the JWE token, and returns the serialized token.
func EncodeJWE(secret string, params JWTEncodeParams) (string, error) {
	// Generate a JWE key.
	key := jose.JSONWebKey{
		Key:       Derive(secret),
		KeyID:     base64.RawURLEncoding.EncodeToString(Derive(secret)),
		Algorithm: alg,
		Use:       "enc",
	}

	// Create a new encrypter.
	encrypter, err := jose.NewEncrypter(jose.A256CBC_HS512, jose.Recipient{Algorithm: alg, Key: key}, nil)
	if err != nil {
		return "", err
	}

	// Marshal the payload to JSON.
	bytes, err := json.Marshal(params.Claims)
	if err != nil {
		return "", err
	}

	// Encrypt the payload.
	encrypted, err := encrypter.Encrypt(bytes)
	if err != nil {
		return "", err
	}

	// Serialize JWE token.
	serialized, err := encrypted.CompactSerialize()
	if err != nil {
		return "", err
	}

	return serialized, nil
}

// DecodeJWE decodes and validates a JWE token.
//
// It decrypts the token using the `Derive` function, unmarshals the payload, and validates the expiration time. If the
// token is valid, it returns the decoded claims.
func DecodeJWE(secret, token string) (*Claims, error) {
	enc, err := jose.Decrypt([]byte(token), jose.WithAlg(string(jose.A256CBC_HS512)), jose.WithPassword(Derive(secret)))
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	if err := json.Unmarshal(enc, claims); err != nil {
		return nil, err
	}

	// Validate expiration.
	if time.Now().Unix() > int64(*claims.Expiry) {
		return nil, erratic.NewUnauthorizedError("reason", "token expired")
	}

	return claims, nil
}

// Derive generates a derived key using HKDF.
//
// It uses the HMAC-SHA256 hash function, the secret key from the shared package, the salt, and the info string to
// derive a 64-byte key. This ensures that the key is unique and unpredictable.
func Derive(secret string) []byte {
	kdf := hkdf.New(sha256.New, []byte(secret), []byte(salt), []byte(info()))
	key := make([]byte, 64)
	_, _ = io.ReadFull(kdf, key)

	return key
}

func GeneratePassword(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
