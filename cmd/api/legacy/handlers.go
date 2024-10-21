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

package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	HealthzResponse struct {
		Status string `json:"status"`
	}
)

func healthz(ctx echo.Context) error {
	if _, err := shared.Temporal().Client().CheckHealth(ctx.Request().Context(), &client.CheckHealthRequest{}); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	if db.Cassandra().Session.Session().S.Closed() {
		return shared.NewAPIError(http.StatusInternalServerError, errors.New("database connection is closed"))
	}

	return ctx.JSON(http.StatusOK, &HealthzResponse{"ok"})
}
