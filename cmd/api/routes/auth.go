package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.breu.io/ctrlplane/internal/common"
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

func (a *authRoutes) register(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	common.Logger.Info(string(body))
	data := &RegRequest{}
	if err := json.Unmarshal(body, data); err != nil {
		utils.HandleHttpError("", err, http.StatusBadRequest, response)
	}

	// Validations are done in the `regRequest` struct. see `cmd/api/routes/requests.go`
	if err := data.save(); err != nil {
		utils.HandleHttpError("", err, http.StatusBadRequest, response)
	}
}

func (a *authRoutes) login(response http.ResponseWriter, request *http.Request)         {}
func (a *authRoutes) logout(response http.ResponseWriter, request *http.Request)        {}
func (a *authRoutes) refreshToken(response http.ResponseWriter, request *http.Request)  {}
func (a *authRoutes) activate(response http.ResponseWriter, request *http.Request)      {}
func (a *authRoutes) resetPassword(response http.ResponseWriter, request *http.Request) {}
func (a *authRoutes) recover(response http.ResponseWriter, request *http.Request)       {}
