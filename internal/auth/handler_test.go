// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Guilospanck/gocqlxmock"
	"github.com/labstack/echo/v4"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/stretchr/testify/assert"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entity"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	// TestHandler is the handler for testing the auth endpoints.
	TestHandler struct {
		session *gocqlxmock.SessionxMock
	}

	RequestData struct {
		data *[]byte
	}
)

func (r *RequestData) Read(p []byte) (n int, err error) {
	return copy(p, *r.data), nil
}

func (r *RequestData) Reset() {
	r.data = nil
}

func (r *RequestData) String() string {
	return string(*r.data)
}

func (r *RequestData) FromRegistrationRequest(request auth.RegisterationRequest) {
	data, _ := json.Marshal(&request)
	r.data = &data
}

// setup creates global mocks for testing environment.
func (handler *TestHandler) setup() {
	shared.InitValidator()

	handler.session = &gocqlxmock.SessionxMock{}
	db.DB.InitMockSession(handler.session)
	db.DB.RegisterValidations()
}

// teardown cleans up the db.
func (handler *TestHandler) teardown() {
	handler.session.On("Close").Return()
	db.DB.Session.Close()

	handler.session = nil
}

// register returns the test function for the register endpoint.
func (handler *TestHandler) register() func(*testing.T) {
	return func(t *testing.T) {
		handler.setup()

		e := echo.New()
		e.Validator = &shared.EchoValidator{Validator: shared.Validator}

		data := &RequestData{}
		reg := auth.RegisterationRequest{
			FirstName:       "John",
			LastName:        "Doe",
			Email:           "johndoe@example.com",
			Password:        "password",
			ConfirmPassword: "password",
			TeamName:        "team",
		}
		data.FromRegistrationRequest(reg)

		eu := entity.User{}
		uclause := qb.EqLit("email", "'johndoe@example.com'")
		ustmt, unames := db.SelectBuilder(eu.GetTable().Name()).AllowFiltering().Columns("id", "email").Where(uclause).ToCql()
		uquerymock := &gocqlxmock.QueryxMock{Ctx: context.Background(), Stmt: ustmt, Names: unames}
		uitermock := &gocqlxmock.IterxMock{}

		uquerymock.On("Iter").Return(uitermock)
		uitermock.On("Unsafe").Return(uitermock)
		uitermock.On("Get", db.NewGetPlaceholder()).Return(nil)
		handler.session.On("Query", ustmt, unames).Return(uquerymock)

		request := httptest.NewRequest("POST", "/auth/register", strings.NewReader(data.String()))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		record := httptest.NewRecorder()
		ctx := e.NewContext(request, record)
		server := &auth.ServerHandler{}

		if assert.Error(t, server.Register(ctx)) {
			assert.Equal(t, http.StatusOK, record.Code) // FIXME: should be 400
		}

		t.Cleanup(handler.teardown)
	}
}

// login returns the test function for the login endpoint.
func (handler *TestHandler) login() func(*testing.T) {
	return func(t *testing.T) {
		handler.setup()
		t.Cleanup(handler.teardown)
	}
}

func TestServerHandler(t *testing.T) {
	s := &TestHandler{}
	t.Run("Register", s.register())
	t.Run("Login", s.login())
}
