package routers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.breu.io/ctrlplane/cmd/api/serializers"
	"go.breu.io/ctrlplane/internal/common/utils"
)

func AuthRouter() http.Handler {
	a := &authRoutes{}
	router := chi.NewRouter()
	router.Post("/register", a.register)
	router.Post("/login", a.login)
	router.Post("/refresh-token", a.refreshToken)
	router.Post("/activate", a.activate)
	router.Post("/reset-password", a.resetPassword)
	return router
}

type authRoutes struct{}

// Registration Route
func (a *authRoutes) register(writer http.ResponseWriter, request *http.Request) {
	serializer := &serializers.RegistrationRequest{}

	if reply, err := serializer.Reply(request.Body); err != nil {
		utils.HandleHTTPError(writer, err, http.StatusBadRequest)
	} else {
		reply, _ := json.Marshal(&reply)
		writer.WriteHeader(http.StatusCreated)
		writer.Write(reply)
	}
}

// Login Route
func (a *authRoutes) login(writer http.ResponseWriter, request *http.Request) {
	serializer := &serializers.LoginRequest{}

	if reply, err := serializer.Reply(request.Body); err != nil {
		utils.HandleHTTPError(writer, err, http.StatusBadRequest)
	} else {
		reply, _ := json.Marshal(&reply)
		writer.WriteHeader(http.StatusCreated)
		writer.Write(reply)
	}
}

func (a *authRoutes) refreshToken(writer http.ResponseWriter, request *http.Request)  {}
func (a *authRoutes) activate(writer http.ResponseWriter, request *http.Request)      {}
func (a *authRoutes) resetPassword(writer http.ResponseWriter, request *http.Request) {}
