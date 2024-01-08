package pong

import (
	"context"
	"log/slog"
)

// Pong represents an event between core domains
type Pong struct {
	Source    string
	Type      string
	RawParams []byte
}

// HandleFunc represents a function that can handle an event
type HandleFunc func(context.Context, Pong) error

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

func (c *Core) Pong() (string, error) {
	c.log.Info("Reached core Pong")
	return "PONG", nil
}
