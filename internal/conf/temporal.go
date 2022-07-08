package conf

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ilyakaznacheev/cleanenv"
	tc "go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type temporal struct {
	ServerHost string `env:"TEMPORAL_HOST" env-default:"temporal"`
	ServerPort string `env:"TEMPORAL_PORT" env-default:"7233"`
	Client     tc.Client
	Queues     struct {
		Webhooks string `env-default:"webhooks"`
	}
}

var Temporal temporal

func (t *temporal) ReadConf() {
	cleanenv.ReadEnv(t)
}

func (t *temporal) GetConnectionString() string {
	return t.ServerHost + ":" + t.ServerPort
}

func (t *temporal) InitClient() {
	Logger.Info(
		"Initializing Temporal Client ...",
		zap.String("host", Temporal.ServerHost),
		zap.String("port", Temporal.ServerPort),
	)
	options := tc.Options{
		HostPort: t.GetConnectionString(),
	}

	retryTemporal := func() error {
		client, err := tc.Dial(options)
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
