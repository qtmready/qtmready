package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	g.POST("/register", register)
	g.POST("/login", login)
}

// register is a handler for /auth/register endpoint
func register(ctx echo.Context) error {
	request := &RegistrationRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return err
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return err
	}

	// Validating team
	team := &entities.Team{Name: request.TeamName}
	if err := ctx.Validate(team); err != nil {
		return err
	}

	user := &entities.User{FirstName: request.FirstName, LastName: request.LastName, Email: request.Email, Password: request.Password}
	if err := ctx.Validate(user); err != nil {
		return err
	}

	if err := db.Save(team); err != nil {
		return err
	}

	user.TeamID = team.ID
	if err := db.Save(user); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &RegistrationResponse{Team: team, User: user})
}

// login is a handler for /auth/login endpoint
func login(ctx echo.Context) error {
	request := &LoginRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return err
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return err
	}

	params := db.QueryParams{"email": "'" + request.Email + "'"}
	user := &entities.User{}

	if err := db.Get(user, params); err != nil {
		return err
	}

	if user.VerifyPassword(request.Password) {
		claims := &shared.JWTClaims{
			UserID:         user.ID.String(),
			TeamID:         user.TeamID.String(),
			StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), Issuer: shared.Service.Name},
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(shared.Service.Secret))
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, &TokenResponse{Token: token})
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
}
