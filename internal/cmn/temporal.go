package cmn

import (
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type (
	QueueName string

	Queue interface {
		CreateWorkflowID(args ...string) string
		GetWorkflowOptions(args ...string) client.StartWorkflowOptions
		GetName() string
	}

	Queues map[QueueName]Queue
)

// TODO: The greater plan is to move each tennant in its own namespace.
const (
	GithubIntegrationQueue QueueName = "github"
	BuilderQueue           QueueName = "builder"
	ProvisionerQueue       QueueName = "provisioner"
	DeployerQueue          QueueName = "deployer"
)

type (
	queue struct {
		Name   QueueName
		Prefix string
	}

	temporal struct {
		ServerHost string `env:"TEMPORAL_HOST" env-default:"temporal"`
		ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
		Client     client.Client
		Queues     Queues
	}
)

func (q QueueName) ToString() string {
	return string(q)
}

func (q *queue) CreateWorkflowID(args ...string) string {
	return q.Prefix + "." + strings.Join(args, ".")
}

func (q *queue) GetName() string {
	return q.Name.ToString()
}

func (q *queue) GetWorkflowOptions(args ...string) client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:        q.CreateWorkflowID(args...),
		TaskQueue: q.GetName(),
	}
}

func (t *temporal) ReadEnv() {
	if err := cleanenv.ReadEnv(t); err != nil {
		Log.Fatal("Failed to read environment variables", zap.Error(err))
	}
}

func (t *temporal) GetConnectionString() string {
	return t.ServerHost + ":" + t.ServerPort
}

func (t *temporal) InitClient() {
	Log.Info("Initializing Temporal Client ...", zap.String("host", t.ServerHost), zap.String("port", t.ServerPort))
	options := client.Options{
		HostPort: t.GetConnectionString(),
	}

	retryTemporal := func() error {
		client, err := client.Dial(options)
		if err != nil {
			return err
		}

		t.Client = client
		Log.Info("Initializing Temporal Client ... Done")
		return nil
	}

	if err := retry.Do(
		retryTemporal,
		retry.Attempts(10),
		retry.Delay(1*time.Second),
	); err != nil {
		Log.Fatal("Failed to initialize Temporal Client", zap.Error(err))
	}
}

func (t *temporal) CreateWorkflowID(args ...string) string {
	return strings.Join(args, ".")
}

var Temporal = &temporal{
	Queues: Queues{
		GithubIntegrationQueue: &queue{GithubIntegrationQueue, "integrations.github"},
	},
}
