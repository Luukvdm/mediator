package mediator

import (
	"context"
)

type (
	// Handler is an interface used by [Behavior].
	Handler interface {
		Handle(ctx context.Context, msg Message) (any, error)
	}
	// HandlerFunc for [Behavior].
	HandlerFunc func(ctx context.Context, msg Message) (any, error)
	// Behavior is a middleware for the [Mediator].
	// Behaviors are wrapped around the handling of a [Notification] or [Request].
	// And executed by the [Pipeline].
	Behavior interface {
		Handler(next Handler) Handler
	}

	chain interface {
		Then(h Handler) Handler
	}
	// Pipeline helps with combining/ chaining multiple [Behavior] instances.
	Pipeline struct {
		behaviors []Behavior
	}
)

// Handle runs the [Handle] function. It is required so that [HandlerFunc] implements the [Handler] interface.
func (f HandlerFunc) Handle(ctx context.Context, msg Message) (any, error) {
	return f(ctx, msg)
}

// Then finalizes the chain build from the [Pipeline] by adding the final [Handler] h to the chain.
func (c Pipeline) Then(h Handler) Handler {
	for i := range c.behaviors {
		next := c.behaviors[len(c.behaviors)-1-i]
		h = next.Handler(h)
	}
	return h
}

// newChain creates a new chain for the given [Behavior] slice.
func newChain(behaviors ...Behavior) chain {
	return Pipeline{behaviors: append(([]Behavior)(nil), behaviors...)}
}
