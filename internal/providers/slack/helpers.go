// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
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


package slack

import (
	"encoding/base64"
	"log/slog"
	"strings"
)

// reveal decodes a base64-encoded encrypted token and decrypts it using a generated key.
func reveal(botToken, workspaceID string) (string, error) {
	// Decode the base64-encoded encrypted token.
	decoded, err := base64.StdEncoding.DecodeString(botToken)
	if err != nil {
		slog.Error("Failed to decode the token", slog.Any("e", err))
		return "", err
	}

	// Generate the same key used for encryption.
	key := generate(workspaceID)

	// Decrypt the token.
	decrypted, err := decrypt(decoded, key)
	if err != nil {
		slog.Error("Failed to decrypt the token", slog.Any("e", err))
		return "", err
	}

	return string(decrypted), nil
}

func FormatFilesList(files []string) string {
	result := ""
	for _, file := range files {
		result += "- " + file + "\n"
	}

	return result
}

func ExtractRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return parts[len(parts)-1]
}
