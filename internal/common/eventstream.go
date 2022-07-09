package common

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type eventstreamconf struct {
	*nats.Conn
	ServerURL string `env:"EVENTS_SERVERS_URL" env-default:"nats://event-stream:4222"`
}

func (e *eventstreamconf) ReadConf() {
	cleanenv.ReadEnv(e)
}

func (e *eventstreamconf) InitConnection() {
	Logger.Info("Initializing Event Stream Client ...", zap.String("url", e.ServerURL))
	conn, err := nats.Connect(e.ServerURL, nats.MaxReconnects(5), nats.ReconnectWait(2*time.Second))
	if err != nil {
		Logger.Fatal("Failed to initialize Event Stream Client", zap.Error(err))
	}
	e.Conn = conn
	Logger.Info("Initializing Event Stream Client ... Done")
}
