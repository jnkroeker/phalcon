package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	pingCore "phalcon/business/core/ping"
	"phalcon/foundation/web"
	"time"

	"github.com/nats-io/nats.go"
)

type Handlers struct {
	Ping *pingCore.Core
	Nats *nats.Conn
}

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *slog.Logger
	Nats     *nats.Conn
}

func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown)

	ping_handlers := Handlers{
		Ping: pingCore.NewCore(cfg.Log),
		Nats: cfg.Nats,
	}

	cfg.Nats.Subscribe("pong", func(m *nats.Msg) {
		cfg.Log.Info("received NATS message: %w", m)
		// mock call down into a python app
		time.Sleep(time.Second * 5)
		cfg.Nats.Publish("ping", []byte("PING"))
	})

	cfg.Nats.Publish("ping", []byte("PING"))

	app.Handle(http.MethodGet, "/ping", ping_handlers.Create)

	return app
}

func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ping, err := h.Ping.Ping()
	if err != nil {
		return fmt.Errorf("Error as result of calling Ping: %w", err)
	}

	h.Nats.Publish("ping", []byte("PING"))

	return web.Respond(ctx, w, ping, http.StatusOK)
}
