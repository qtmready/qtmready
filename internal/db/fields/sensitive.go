// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

// Copyright © 2024, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package fields

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
	// Sensitive represents a string encrypted using AES-GCM.
	//
	// It provides encryption and decryption of sensitive data within the application, ensuring that the data is stored
	// and transmitted securely without being exposed in plain text. This is particularly useful for protecting sensitive
	// values, both at rest and in motion.
	//
	// The encryption key is derived from a secret by calling shared.Service(). shared.Service is a singleton initialized
	// at application startup using environment variables.
	//
	// Usage:
	//
	//     type Entity struct {
	//         ID        gocql.UUID        `json:"id" cql:"id"`
	//         Sensitive db.Sensitive `json:"sensitive" cql:"sensitive"`
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
	//
	// Open Questions:
	//
	// - What is the ciphertext size?
	// - How are errors handled during encryption/decryption?
	// - What are the performance implications for large data?
	// - Are there any known security vulnerabilities with the AES-GCM implementation?
	// - What is the maximum plaintext length?
	// - How are non-ASCII characters handled?
	// - Is key configuration possible? How is key exposure prevented?
	// - What is the performance impact on the application?
	// - Are there any specific use cases where Sensitive is well-suited?
	// - How can we handle secret rotation? Maybe a version prefix with Vault integration?
	Sensitive string
)

func (sen Sensitive) String() string {
	return string(sen)
}

// secret truncates or pads the secret to 32 bytes for AES-GCM.
func (sen Sensitive) secret() []byte {
	s := []byte(shared.Service().GetSecret())

	if len(s) > 32 {
		return s[:32]
	}

	if len(s) < 32 {
		return append(s, make([]byte, 32-len(s))...)
	}

	return s
}

// encrypt encrypts the string value and returns the encrypted data.
func (sen Sensitive) encrypt() ([]byte, error) {
	plain := []byte(string(sen))

	// Use the secret directly as the key
	block, err := aes.NewCipher(sen.secret())
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
func (sen *Sensitive) from(encrypted []byte) error {
	// Use the secret directly as the key
	block, err := aes.NewCipher(sen.secret())
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

	*sen = Sensitive(string(plain))

	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (sen Sensitive) MarshalJSON() ([]byte, error) {
	encrypted, err := sen.encrypt()
	if err != nil {
		return nil, err
	}

	return json.Marshal(base64.StdEncoding.EncodeToString(encrypted))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (sen *Sensitive) UnmarshalJSON(data []byte) error {
	var encrypted string
	if err := json.Unmarshal(data, &encrypted); err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return err
	}

	return sen.from(decoded)
}

// MarshalCQL returns the encrypted string for storing in Cassandra.
func (sen Sensitive) MarshalCQL() ([]byte, error) {
	encrypted, err := sen.encrypt()
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// UnmarshalCQL decodes the encrypted string from Cassandra.
func (sen *Sensitive) UnmarshalCQL(b []byte) error {
	return sen.from(b)
}
