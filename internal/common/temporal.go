package common

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type temporalconf struct {
	ServerHost string `env:"TEMPORAL_HOST" env-default:"temporal"`
	ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
	Client     client.Client
	Queues     struct {
		Integrations string `env-default:"integrations"`
	}
}

func (t *temporalconf) ReadEnv() {
	cleanenv.ReadEnv(t)
}

func (t *temporalconf) GetConnectionString() string {
	return t.ServerHost + ":" + t.ServerPort
}

func (t *temporalconf) InitClient() {
	Logger.Info("Initializing Temporal Client ...", zap.String("host", t.ServerHost), zap.String("port", t.ServerPort))
	options := client.Options{
		HostPort: t.GetConnectionString(),
	}

	retryTemporal := func() error {
		client, err := client.Dial(options)
		if err != nil {
			return err
		}

		t.Client = client
		Logger.Info("Initializing Temporal Client ... Done")
		return nil
	}

	if err := retry.Do(
		retryTemporal,
		retry.Attempts(10),
		retry.Delay(1*time.Second),
	); err != nil {
		Logger.Fatal("Failed to initialize Temporal Client", zap.Error(err))
	}
}
