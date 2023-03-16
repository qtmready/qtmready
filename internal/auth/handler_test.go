package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Pallinder/go-randomdata"
	pwg "github.com/sethvargo/go-password/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
	"go.breu.io/ctrlplane/internal/testutils"
)

type (
	Containers struct {
		network    testcontainers.Network
		db         *testutils.Container
		temporal   *testutils.Container
		nats       *testutils.Container
		api        *testutils.Container
		mothership *testutils.Container
	}

	RequestData struct {
		register *auth.RegisterationRequest
	}

	ResponseData struct {
		register *auth.RegisterationResponse
	}

	ServerHandlerTestSuite struct {
		suite.Suite
		context  context.Context
		ctrs     *Containers
		url      string
		client   *auth.Client
		requests *RequestData
	}
)

func (c *Containers) shutdown(ctx context.Context) {
	shared.Logger.Info("graceful shutdown test environment ...")
	_ = c.api.Shutdown()
	_ = c.mothership.Shutdown()
	_ = c.temporal.Shutdown()
	_ = c.nats.Shutdown()
	_ = c.db.DropKeyspace(db.TestKeyspace)
	_ = c.db.ShutdownCassandra()
	_ = c.network.Remove(ctx)
	db.DB.Session.Close()
	shared.Logger.Info("graceful shutdown complete.")
}

func (s *ServerHandlerTestSuite) SetupSuite() {
	s.context = context.Background()
	s.SetupContainers()
	s.SetupAPIClient()
	s.SetupRequestData()
}

func (s *ServerHandlerTestSuite) TearDownSuite() {
	s.ctrs.shutdown(context.Background())
}

func (s *ServerHandlerTestSuite) SetupContainers() {
	shared.InitForTest()
	shared.Logger.Info("setting up test environment ...")
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

	_ = db.DB.InitSessionForTests(port.Int(), "file://../db/migrations")
	shared.Logger.Warn("session gets initiated, but if we catch the error and do t.Fatal(err), the test panics!")
	if db.DB.Session.Session().S == nil {
		s.T().Fatal("session is nil")
	}

	db.DB.RunMigrations()

	temporalctr, err := testutils.StartTemporalContainer(s.context)
	if err != nil {
		s.T().Fatalf("failed to start temporal container: %v", err)
	}

	natsctr, err := testutils.StartNatsIOContainer(s.context)
	if err != nil {
		s.T().Fatalf("failed to start natsio container: %v", err)
	}

	apictr, err := testutils.StartAPIContainer(s.context)
	if err != nil {
		s.T().Fatalf("failed to start api container: %v", err)
	}

	mxctr, err := testutils.StartMothershipContainer(s.context)

	dbhost, _ := dbctr.Container.ContainerIP(s.context)
	temporalhost, _ := temporalctr.Container.ContainerIP(s.context)
	natshost, _ := natsctr.Container.ContainerIP(s.context)
	apihost, _ := apictr.Container.ContainerIP(s.context)
	mxhost, _ := mxctr.Container.ContainerIP(s.context)

	shared.Logger.Info("hosts ...", "db", dbhost, "temporal", temporalhost, "nats", natshost, "api", apihost, "mothership", mxhost)

	s.ctrs = &Containers{
		network:    network,
		db:         dbctr,
		temporal:   temporalctr,
		nats:       natsctr,
		api:        apictr,
		mothership: mxctr,
	}
}

func (s *ServerHandlerTestSuite) SetupRequestData() {
	password := pwg.MustGenerate(16, 4, 4, true, true)
	s.requests = &RequestData{}
	s.requests.register = &auth.RegisterationRequest{
		Email:           randomdata.Email(),
		Password:        password,
		ConfirmPassword: password,
		FirstName:       randomdata.FirstName(randomdata.Male),
		LastName:        randomdata.LastName(),
		TeamName:        randomdata.SillyName(),
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

func (s *ServerHandlerTestSuite) TestRegister() {
	response, err := s.client.Register(s.context, *s.requests.register)
	if err != nil {
		s.T().Fatalf("failed to register: %v", err)
	}

	parsed, err := auth.ParseRegisterResponse(response)

	assert.Equal(s.T(), http.StatusCreated, response.StatusCode)
	assert.Equal(s.T(), s.requests.register.Email, parsed.JSON201.User.Email)
	assert.Equal(s.T(), s.requests.register.FirstName, parsed.JSON201.User.FirstName)
	assert.Equal(s.T(), s.requests.register.LastName, parsed.JSON201.User.LastName)
	assert.Equal(s.T(), parsed.JSON201.User.TeamID, parsed.JSON201.Team.ID)
	assert.Equal(s.T(), s.requests.register.TeamName, parsed.JSON201.Team.Name)
}

func (s *ServerHandlerTestSuite) TestRegister_FailOnDuplicateEmail() {
	response, err := s.client.Register(s.context, *s.requests.register)
	if err != nil {
		s.T().Fatalf("failed to register: %v", err)
	}

	assert.Equal(s.T(), http.StatusBadRequest, response.StatusCode)
	parsed, err := auth.ParseRegisterResponse(response)
	assert.Contains(s.T(), parsed.JSON400.Message, "db_unique")
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(ServerHandlerTestSuite))
}
