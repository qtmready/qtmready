package slack

import (
	"encoding/base64"
	"log/slog"
)

// DecodeAndDecryptToken decodes a base64-encoded encrypted token and decrypts it using a generated key.
func decodeAndDecryptToken(botToken, workspaceID string) (string, error) {
	// Decode the base64-encoded encrypted token.
	decoded, err := base64.StdEncoding.DecodeString(botToken)
	if err != nil {
		slog.Error("Failed to decode the token", slog.Any("e", err))
		return "", err
	}

	// Generate the same key used for encryption.
	key := generateKey(workspaceID)

	// Decrypt the token.
	decryptedToken, err := decrypt(decoded, key)
	if err != nil {
		slog.Error("Failed to decrypt the token", slog.Any("e", err))
		return "", err
	}

	return string(decryptedToken), nil
}
