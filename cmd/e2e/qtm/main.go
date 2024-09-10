// Copyright © 2023, Breu, Inc. <info@breu.io>
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

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"
	"regexp"

	ghi "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/uuid"

	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/shared"
)

// block, _ := pem.Decode([]byte(github.Instance().PrivateKey))
// key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

// shared.Logger().Info("decoded ....", slog.Any("block", block), slog.Any("key", key))

func main() {
	client, _ := ghi.New(http.DefaultTransport, github.Instance().AppID, 50886707, []byte(github.Instance().PrivateKey))
	token, _ := client.Token(context.Background())

	url := fmt.Sprintf("https://git:%s@github.com/breuHQ/governance.git", token)

	id, _ := uuid.NewV7()

	shared.Logger().Info("token ....", slog.Any("token", token), slog.String("url", url), slog.String("id", id.String()))

	{
		cmd := exec.Command("git", "clone", url, "--single-branch", "--branch", "updated_user_list", "--depth", "1", fmt.Sprintf("/tmp/%s", id.String()))
		out, _ := cmd.CombinedOutput()

		shared.Logger().Info("output ....", slog.Any("output", out))
	}

	{
		cmd := exec.Command("git", "-C", fmt.Sprintf("/tmp/%s", id.String()), "fetch", "origin", "main")
		out, _ := cmd.CombinedOutput()

		shared.Logger().Info("output ....", slog.Any("output", out))
	}

	{
		cmd := exec.CommandContext(context.Background(), "git", "-C", fmt.Sprintf("/tmp/%s", id.String()), "rebase", "d2c649da85e1ba213643542501987a5b6696f6ea")

		out, err := cmd.CombinedOutput()
		if err != nil {
			var exerr *exec.ExitError

			if errors.As(err, &exerr) {
				str := err.Error()
				pattern := `(?m)^Could not apply ([0-9a-fA-F]{7})\.\.\. (.*)$`

				// Compile the regex
				re := regexp.MustCompile(pattern)

				// Find all matches
				matches := re.FindAllStringSubmatch(str, -1)
				for _, match := range matches {
					shared.Logger().Info(match[0])
					shared.Logger().Info(match[1])
					shared.Logger().Info(match[2])
				}
			}
		}

		pattern := `(?m)^Could not apply ([0-9a-fA-F]{7})\.\.\. (.*)$`

		// Compile the regex
		re := regexp.MustCompile(pattern)

		// Find all matches
		matches := re.FindAllStringSubmatch(string(out), -1)
		// for _, match := range matches {
		// 	shared.Logger().Info(match[0])
		// 	shared.Logger().Info(match[1])
		// 	shared.Logger().Info(match[2])
		// }

		if len(matches) > 0 {
			sha, msg := matches[0][1], matches[0][2]

			shared.Logger().Info("matches ....", slog.String("sha", sha), slog.String("msg", msg))
		}

		shared.Logger().Info("output ....", slog.Any("output", out))
	}
}

// git clone https://git:<token>@github.com/owner/repo.git
