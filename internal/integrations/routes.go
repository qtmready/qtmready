package integrations

import (
	"github.com/labstack/echo/v4"
	"go.breu.io/ctrlplane/internal/integrations/github"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	github.CreateRoutes(g.Group("/github"))
}
