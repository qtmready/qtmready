package cmn

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var EventStream = &eventstream{}

type eventstream struct {
	*nats.Conn
	ServerURL string `env:"EVENTS_SERVERS_URL" env-default:"nats://event-stream:4222"`
}

func (e *eventstream) ReadEnv() {
	if err := cleanenv.ReadEnv(e); err != nil {
		Logger.Error("Failed to read environment variables", zap.Error(err))
	}
}

func (e *eventstream) InitConnection() {
	Logger.Info("Initializing Event Stream Client ...", zap.String("url", e.ServerURL))
	conn, err := nats.Connect(e.ServerURL, nats.MaxReconnects(5), nats.ReconnectWait(2*time.Second))
	if err != nil {
		Logger.Error("Failed to initialize Event Stream Client", zap.Error(err))
	}
	e.Conn = conn
	Logger.Info("Initializing Event Stream Client ... Done")
}
