package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	pb "phalcon/business/core/pb"
	pongCore "phalcon/business/core/pong"
	"phalcon/foundation/web"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type Handlers struct {
	Pong *pongCore.Core
}

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *slog.Logger
	Nats     *nats.Conn
}

func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown)

	pong_handlers := Handlers{
		Pong: pongCore.NewCore(cfg.Log),
	}

	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cfg.Log.Error("error dialing grpc host: %w", err)
	}

	client := pb.NewOutliersClient(conn)

	cfg.Nats.Subscribe("ping", func(m *nats.Msg) {
		cfg.Log.Info("received NATS message: %w", m)
		// mock call down into a python app
		// time.Sleep(time.Second * 5)

		req := pb.OutliersRequest{
			Metrics: dummyData(),
		}

		resp, err := client.Detect(context.Background(), &req)
		if err != nil {
			cfg.Log.Error("error calling Detect: %w", err)
		}
		cfg.Nats.Publish("pong", []byte("PONG"))

		cfg.Log.Info("outliers at: %v", resp.Indicies)
	})

	app.Handle(http.MethodPost, "/pong", pong_handlers.Create)

	return app
}

func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	pong, err := h.Pong.Pong()
	if err != nil {
		return fmt.Errorf("Error as result of call to Pong: %w", err)
	}

	return web.Respond(ctx, w, pong, http.StatusCreated)
}

func dummyData() []*pb.Metric {
	const size = 1000
	out := make([]*pb.Metric, size)
	for i := 0; i < size; i++ {
		m := pb.Metric{
			Name: "CPU",
			// normally we're below 40% CPU utilization
			Value: rand.Float64() * 40,
		}
		out[i] = &m
	}
	// Create some outliers
	out[7].Value = 97.3
	out[113].Value = 92.1
	out[835].Value = 93.2
	return out
}
