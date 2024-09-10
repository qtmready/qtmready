// Copyright © 2023, 2024, Breu, Inc. <info@breu.io>
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

package auth

import (
	"errors"
)

var (
	ErrInvalidAPIKey         = errors.New("invalid API key")
	ErrInvalidAuthHeader     = errors.New("invalid authorization header")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrMalformedAPIKey       = errors.New("malformed API key")
	ErrMissingAuthHeader     = errors.New("no authorization header provided")
	ErrCrypto                = errors.New("crypto error")
	ErrTokenExpired          = errors.New("token has expired")
)
