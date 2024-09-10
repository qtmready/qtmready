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

	if db.DB().Session.Session().S.Closed() {
		return shared.NewAPIError(http.StatusInternalServerError, errors.New("database connection is closed"))
	}

	return ctx.JSON(http.StatusOK, &HealthzResponse{"ok"})
}
