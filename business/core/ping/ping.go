package ping

import (
	"context"
	"log/slog"
)

// Ping represents an event between core domains
type Ping struct {
	Source    string
	Type      string
	RawParams []byte
}

// HandleFunc represents a function that can handle an event
type HandleFunc func(context.Context, Ping) error

// Core manages the set of APIs for event access.
type Core struct {
	log      *slog.Logger
	handlers map[string]map[string][]HandleFunc
}

// NewCore constructs a Core for event API access.
func NewCore(log *slog.Logger) *Core {
	return &Core{
		log:      log,
		handlers: map[string]map[string][]HandleFunc{},
	}
}

func (c *Core) Ping() (string, error) {
	c.log.Info("Reached Core Ping")
	return "PING", nil
}
