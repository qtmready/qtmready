package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"

	"go.breu.io/quantm/internal/shared"
)

type (
	// EncryptedField represents a string encrypted using AES-GCM.
	//
	// It provides encryption and decryption of sensitive data within the application,ensuring that the data is stored and transmitted
	// securely without being exposed in plain text. This is particularly useful for protecting sensitive values in databases or during
	// network communication.
	//
	// The encryption key is derived from the secret stored in `shared.Service()`.
	// Example usage:
	//
	//     type Entity struct {
	//         ID        gocql.UUID        `json:"id" cql:"id"`
	//         Sensitive db.EncryptedField `json:"sensitive" cql:"sensitive"`
	//     }
	//
	//     // Create a new instance of Entity with an encrypted field.
	//     entity := Entity{
	//         ID:        gocql.TimeUUID(),
	//         Sensitive: db.NewEncryptedField("my_secret_value"),
	//     }
	//
	//     // Save the entity to the database, ensuring the sensitive field is encrypted.
	//     err := db.Save(entity)
	//     if err != nil {
	//         // Handle error
	//     }
	//
	//     // Retrieve the entity from the database by ID, automatically decrypting the sensitive field.
	//     var retrievedEntity Entity
	//     err = db.Get(&retrievedEntity, db.QueryParams{"id": entity.ID.String()})
	//     if err != nil {
	//         // Handle error
	//     }
	//
	//     // Access the decrypted sensitive value directly.
	//     fmt.Println(retrievedEntity.Sensitive.String()) // Outputs: "my_secret_value"
	EncryptedField string
)

func (ef EncryptedField) String() string {
	return string(ef)
}

// encrypt encrypts the string value and returns the encrypted data.
func (ef EncryptedField) encrypt() ([]byte, error) {
	plain := []byte(string(ef))

	// Use the secret directly as the key
	block, err := aes.NewCipher([]byte(shared.Service().GetSecret()))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	sealed := gcm.Seal(nonce, nonce, plain, nil)

	return sealed, nil
}

// from rebuilds the encrypted string from the encrypted data.
func (ef *EncryptedField) from(encrypted []byte) error {
	// Use the secret directly as the key
	block, err := aes.NewCipher([]byte(shared.Service().GetSecret()))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := encrypted[:gcm.NonceSize()]
	ciphertext := encrypted[gcm.NonceSize():]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	*ef = EncryptedField(string(plain))

	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (ef EncryptedField) MarshalJSON() ([]byte, error) {
	encrypted, err := ef.encrypt()
	if err != nil {
		return nil, err
	}

	return json.Marshal(base64.StdEncoding.EncodeToString(encrypted))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (ef *EncryptedField) UnmarshalJSON(data []byte) error {
	var encrypted string
	if err := json.Unmarshal(data, &encrypted); err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return err
	}

	return ef.from(decoded)
}

// MarshalCQL returns the encrypted string for storing in Cassandra.
func (ef EncryptedField) MarshalCQL() ([]byte, error) {
	encrypted, err := ef.encrypt()
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// UnmarshalCQL decodes the encrypted string from Cassandra.
func (ef *EncryptedField) UnmarshalCQL(b []byte) error {
	return ef.from(b)
}
