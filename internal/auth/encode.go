// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

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
