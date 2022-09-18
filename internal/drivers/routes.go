// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.  

package drivers

import (
	"github.com/labstack/echo/v4"
	"go.breu.io/ctrlplane/internal/drivers/github"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	github.CreateRoutes(g.Group("/github"), middlewares...)
}
