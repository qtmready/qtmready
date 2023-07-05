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
	"fmt"
	"net/http"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/golang-jwt/jwt/v4"
	pwg "github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
	"go.breu.io/ctrlplane/internal/testutils"
)

type (
	Containers struct {
		network  testcontainers.Network
		db       *testutils.Container
		temporal *testutils.Container
		// nats       *testutils.Container
		api        *testutils.Container
		mothership *testutils.Container
	}

	RequestData struct {
		register *auth.RegisterationRequest
		login    *auth.LoginRequest
	}

	ResponseData struct {
		register *auth.RegisterationResponse
	}

	ServerHandlerTestSuite struct {
		suite.Suite
		context   context.Context
		ctrs      *Containers
		url       string
		client    *auth.Client
		requests  *RequestData
		responses *ResponseData
	}
)

func (c *Containers) shutdown(ctx context.Context) {
	shared.Logger().Info("graceful shutdown test environment ...")

	_ = c.api.Shutdown()
	_ = c.mothership.Shutdown()
	_ = c.temporal.Shutdown()
	_ = c.db.DropKeyspace(db.TestKeyspace)
	_ = c.db.ShutdownCassandra()
	_ = c.network.Remove(ctx)

	db.DB().Session.Close()
	shared.Logger().Info("graceful shutdown complete.")
}

func (s *ServerHandlerTestSuite) SetupSuite() {
	shared.InitServiceForTest()

	s.context = context.Background()
	s.SetupContainers()
	s.SetupAPIClient()
	s.SetupRequestData()
	s.responses = &ResponseData{}
}

func (s *ServerHandlerTestSuite) TearDownSuite() {
	s.ctrs.shutdown(context.Background())
}

func (s *ServerHandlerTestSuite) SetupContainers() {
	shared.Logger().Info("setting up test environment ...")

	network, err := testutils.CreateTestNetwork(s.context)
	if err != nil {
		s.T().Fatalf("failed to create test network: %v", err)
	}

	dbctr, err := testutils.StartDBContainer(s.context)
	if err != nil {
		s.T().Fatalf("failed to start db container: %v", err)
	}

	if err = dbctr.CreateKeyspace(db.TestKeyspace); err != nil {
		s.T().Fatalf("failed to create keyspace: %v", err)
	}

	port, err := dbctr.Container.MappedPort(context.Background(), "9042")
	if err != nil {
		s.T().Fatalf("failed to get mapped db port: %v", err)
	}

	db.NewE2ESession(port.Int(), "file://../db/migrations")

	// _ = db.DB.InitSessionForTests(port.Int(), "file://../db/migrations")

	temporalctr, err := testutils.StartTemporalContainer(s.context)
	if err != nil {
		s.T().Fatalf("failed to start temporal container: %v", err)
	}

	apictr, err := testutils.StartAPIContainer(s.context, shared.Service().GetSecret())
	if err != nil {
		s.T().Fatalf("failed to start api container: %v", err)
	}

	mxctr, err := testutils.StartMothershipContainer(s.context, shared.Service().GetSecret())
	if err != nil {
		s.T().Fatalf("failed to start api container: %v", err)
	}

	dbhost, _ := dbctr.Container.ContainerIP(s.context)
	temporalhost, _ := temporalctr.Container.ContainerIP(s.context)
	apihost, _ := apictr.Container.ContainerIP(s.context)
	mxhost, _ := mxctr.Container.ContainerIP(s.context)

	shared.Logger().Info("hosts ...", "db", dbhost, "temporal", temporalhost, "api", apihost, "mothership", mxhost)

	s.ctrs = &Containers{
		network:    network,
		db:         dbctr,
		temporal:   temporalctr,
		api:        apictr,
		mothership: mxctr,
	}
}

func (s *ServerHandlerTestSuite) GenRegistrationRequest() *auth.RegisterationRequest {
	password := pwg.MustGenerate(16, 4, 4, true, true)

	return &auth.RegisterationRequest{
		Email:           randomdata.Email(),
		Password:        password,
		ConfirmPassword: password,
		FirstName:       randomdata.FirstName(randomdata.Male),
		LastName:        randomdata.LastName(),
		TeamName:        randomdata.SillyName(),
	}
}

func (s *ServerHandlerTestSuite) SetupRequestData() {
	s.requests = &RequestData{}
	s.requests.register = s.GenRegistrationRequest()
	s.requests.login = &auth.LoginRequest{
		Email:    s.requests.register.Email,
		Password: s.requests.register.Password,
	}
}

func (s *ServerHandlerTestSuite) SetupLoginData() {
	s.requests.login = &auth.LoginRequest{
		Email:    s.requests.register.Email,
		Password: s.requests.register.Password,
	}
}

func (s *ServerHandlerTestSuite) SetupAPIClient() {
	port, _ := s.ctrs.api.Container.MappedPort(context.Background(), "8000")
	s.url = fmt.Sprintf("http://localhost:%d", port.Int())

	client, err := auth.NewClient(s.url)
	if err != nil {
		s.T().Fatalf("failed to create api client: %v", err)
	}

	s.client = client
}

func (s *ServerHandlerTestSuite) Test_0001_Register() {
	response, err := s.client.Register(s.context, *s.requests.register)
	if err != nil {
		s.T().Fatalf("failed to register: %v", err)
	}

	defer response.Body.Close()

	parsed, err := auth.ParseRegisterResponse(response)
	if err != nil {
		s.T().Fatalf("failed to parse register response: %v", err)
	}

	s.Assert().Equal(http.StatusCreated, response.StatusCode)
	s.Assert().NotNil(parsed.JSON201)
	s.Assert().Equal(s.requests.register.Email, parsed.JSON201.User.Email)
	s.Assert().Equal(s.requests.register.FirstName, parsed.JSON201.User.FirstName)
	s.Assert().Equal(s.requests.register.LastName, parsed.JSON201.User.LastName)
	s.Assert().Equal(parsed.JSON201.User.TeamID, parsed.JSON201.Team.ID)
	s.Assert().Equal(s.requests.register.TeamName, parsed.JSON201.Team.Name)

	s.responses.register = parsed.JSON201
}

func (s *ServerHandlerTestSuite) Test_0002_Register_DuplicateEmail() {
	response, err := s.client.Register(s.context, *s.requests.register)
	if err != nil {
		s.T().Fatalf("failed to register: %v", err)
	}

	defer response.Body.Close()
	s.Assert().Equal(http.StatusBadRequest, response.StatusCode)

	parsed, _ := auth.ParseRegisterResponse(response)
	s.Assert().NotNil(parsed.JSON400)
	s.Assert().Equal(parsed.JSON400.Message, "validation error")

	emailerr, ok := parsed.JSON400.Errors.Get("email")
	s.Assert().True(ok)
	s.Assert().Equal(emailerr, "already exists")
}

func (s *ServerHandlerTestSuite) Test_0003_Register_InvalidEmail() {
	request := s.GenRegistrationRequest()
	request.Email = "invalid"

	response, err := s.client.Register(s.context, *request)
	if err != nil {
		s.T().Fatalf("failed to register: %v", err)
	}

	defer response.Body.Close()

	s.Assert().Equal(http.StatusBadRequest, response.StatusCode)

	parsed, _ := auth.ParseRegisterResponse(response)
	s.Assert().NotNil(parsed.JSON400)
	s.Assert().Equal(parsed.JSON400.Message, "validation error")

	emailerr, ok := parsed.JSON400.Errors.Get("email")
	s.Assert().True(ok)
	s.Assert().Equal(emailerr, "invalid format")
}

func (s *ServerHandlerTestSuite) Test_0101_Login() {
	response, err := s.client.Login(s.context, *s.requests.login)
	if err != nil {
		s.T().Fatalf("failed to login: %v", err)
	}

	defer response.Body.Close()

	s.Assert().Equal(http.StatusOK, response.StatusCode)

	parsed, err := auth.ParseLoginResponse(response)
	if err != nil {
		s.T().Fatalf("failed to parse login response: %v", err)
	}

	s.Assert().NotNil(parsed.JSON200)
	s.Assert().NotNil(parsed.JSON200.AccessToken)

	access := parsed.JSON200.AccessToken

	paccess, err := jwt.NewParser().ParseWithClaims(access, &auth.JWTClaims{}, auth.SecretFn)
	if err != nil {
		s.T().Fatalf("failed to parse access token: %v", err)
	}

	if claims, ok := paccess.Claims.(*auth.JWTClaims); ok {
		s.Assert().Equal(claims.UserID, s.responses.register.User.ID.String())
		s.Assert().Equal(claims.TeamID, s.responses.register.Team.ID.String())
	} else {
		s.T().Fatalf("failed to parse claims")
	}
}

// TestHandler runs the test suite for the server handler. Please note that suite tests run in alphabetical order.
// In the light of this, we prefix the test functions with numbers to ensure the correct order of execution.
//
//   - 00xx - Register
//   - 01xx - Login
//   - 02xx - API Key
//
// This is a bookkeeping nightmare, but it's the easiest way to ensure that the tests are run in the correct order. For
// more information, please see [github] issue.
//
// [github]: https://github.com/stretchr/testify/issues/194
func TestHandler(t *testing.T) {
	suite.Run(t, new(ServerHandlerTestSuite))
}
