package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	body, _ := ioutil.ReadAll(request.Body)
	response := &RegisterResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		utils.HandleHttpError("", err, http.StatusBadRequest, writer)
	}
}

func (a *authRoutes) login(writer http.ResponseWriter, request *http.Request)         {}
func (a *authRoutes) logout(writer http.ResponseWriter, request *http.Request)        {}
func (a *authRoutes) refreshToken(writer http.ResponseWriter, request *http.Request)  {}
func (a *authRoutes) activate(writer http.ResponseWriter, request *http.Request)      {}
func (a *authRoutes) resetPassword(writer http.ResponseWriter, request *http.Request) {}
func (a *authRoutes) recover(writer http.ResponseWriter, request *http.Request)       {}
