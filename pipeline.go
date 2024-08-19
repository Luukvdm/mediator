package mediator

import (
	"context"
	"log/slog"
)

type (
	// Handler is an interface used by [Behavior].
	Handler interface {
		Handle(ctx context.Context, l *slog.Logger, msg Message) (any, error)
	}

	// HandlerFunc for [Behavior].
	HandlerFunc func(ctx context.Context, l *slog.Logger, msg Message) (any, error)

	// Behavior is a middleware for the [Mediator].
	// Behaviors are wrapped around the handling of a [Notification] or [Request].
	// And executed by the [Pipeline].
	Behavior interface {
		Handler(next Handler) Handler
	}

	// Pipeline helps with combining/ chaining multiple [Behavior] instances.
	Pipeline interface {
		Then(hf HandlerFunc) Handler
	}

	// pipeline is the default [Pipeline] implementation.
	// It simply chains behaviors.
	pipeline struct {
		behaviors []Behavior
	}
)

// Handle runs the [Handle] function. It is required so that [HandlerFunc] implements the [Handler] interface.
func (f HandlerFunc) Handle(ctx context.Context, l *slog.Logger, msg Message) (any, error) {
	return f(ctx, l, msg)
}

// Then creates a handler chain from the [Pipeline].
// [Handler] h is the final piece of the chain and the message being handled.
func (c pipeline) Then(hf HandlerFunc) Handler {
	var h Handler = hf
	for i := range c.behaviors {
		next := c.behaviors[len(c.behaviors)-1-i]
		h = next.Handler(h)
	}
	return h
}

// newPipeline creates a new pipeline for the given [Behavior] slice.
func newPipeline(behaviors ...Behavior) Pipeline {
	return pipeline{behaviors: append(([]Behavior)(nil), behaviors...)}
}
