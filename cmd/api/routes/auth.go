package routes

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
	router.Post("/logout", a.logout)
	router.Post("/refresh-token", a.refreshToken)
	router.Post("/activate", a.activate)
	router.Post("/reset-password", a.resetPassword)
	router.Post("/recover", a.recover)
	return router
}

type authRoutes struct{}

func (a *authRoutes) register(writer http.ResponseWriter, request *http.Request) {
	data := &serializers.RegistrationRequest{}

	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		utils.HandleHttpError("", err, http.StatusBadRequest, writer)
	}

	// Validations are done in the `requests.Registration`
	if user, err := data.Save(); err != nil {
		utils.HandleHttpError("", err, http.StatusBadRequest, writer)
	} else {
		if response, err := json.Marshal(serializers.RegisterationResponse{User: user}); err != nil {
			utils.HandleHttpError("", err, http.StatusBadRequest, writer)
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(response))
		}
	}
}

func (a *authRoutes) login(writer http.ResponseWriter, request *http.Request)         {}
func (a *authRoutes) logout(writer http.ResponseWriter, request *http.Request)        {}
func (a *authRoutes) refreshToken(writer http.ResponseWriter, request *http.Request)  {}
func (a *authRoutes) activate(writer http.ResponseWriter, request *http.Request)      {}
func (a *authRoutes) resetPassword(writer http.ResponseWriter, request *http.Request) {}
func (a *authRoutes) recover(writer http.ResponseWriter, request *http.Request)       {}
