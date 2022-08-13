package main

import (
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.breu.io/ctrlplane/cmd/api/routes/auth"
	"go.breu.io/ctrlplane/internal/cmn"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/integrations"
	"go.breu.io/ctrlplane/internal/integrations/github"
)

var (
	waiter sync.WaitGroup
)

type (
	EchoValidator struct {
		validator *validator.Validate
	}
)

func (ev *EchoValidator) Validate(i interface{}) error {
	if err := ev.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func init() {
	// Reading the configuration from the environment
	cmn.Service.ReadEnv()
	cmn.Service.InitLogger()
	cmn.Service.InitValidator()
	cmn.EventStream.ReadEnv()
	db.DB.ReadEnv()
	db.DB.RegisterValidations()
	cmn.Temporal.ReadEnv()
	github.Github.ReadEnv()
	// Reading the configuration from the environment ... Done

	// Initializing reference to adapters
	waiter.Add(3)

	go func() {
		defer waiter.Done()
		db.DB.InitSessionWithMigrations()
	}()

	go func() {
		defer waiter.Done()
		cmn.EventStream.InitConnection()
	}()

	go func() {
		defer waiter.Done()
		cmn.Temporal.InitClient()
	}()

	waiter.Wait()
	// Initializing singleton objects ... Done

	cmn.Log.Info("Initializing Service ... Done")
}

func main() {
	// handling closing of the server
	defer db.DB.Session.Close()
	defer cmn.Temporal.Client.Close()

	e := echo.New()

	e.Validator = &EchoValidator{validator: cmn.Validator}

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Unauthenticated routes
	e.GET("/healthcheck", healthcheck)
	auth.CreateRoutes(e.Group("/auth"))
	integrations.CreateRoutes(e.Group("/integrations"))

	// Protected routes
	jwtconf := middleware.JWTConfig{
		Claims:     &auth.JWTClains{},
		SigningKey: []byte(cmn.Service.Secret),
	}
	protected := e.Group("/")
	protected.Use(middleware.JWTWithConfig(jwtconf))

	e.Start(":8000")
}

func healthcheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}
