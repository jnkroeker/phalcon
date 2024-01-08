package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"phalcon/app/pong/handlers"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	logger.Info("starting Pong service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// connect to NATS.io server
	nc, err := nats.Connect("nats.default.svc.cluster.local:4222")
	if err != nil {
		logger.Error("unable to connect to NATS server")
	}

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      logger,
		Nats:     nc,
	})

	api := http.Server{
		Addr:         "0.0.0.0:4000",
		Handler:      apiMux,
		ReadTimeout:  5,
		WriteTimeout: 10,
		IdleTimeout:  120,
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		logger.Info("service_two server error", err)
	case sig := <-shutdown:
		logger.Info("shutdown", "status", "service_two shutdown started", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
		}
	}
}
