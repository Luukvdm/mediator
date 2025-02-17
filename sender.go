package mediator

import (
	"context"
	"log/slog"
)

type (
	// Sender is used by the [Send] function.
	// It holds the [Pipeline] that messages pass through before they get handled.
	//
	// The interface is implemented by [Mediator].
	Sender interface {
		getRequestPipeline() Pipeline
		getLogger() *slog.Logger
	}

	// Request is an object that can be sent through the [Mediator].
	//
	// Requests are often used in the CQRS design pattern.
	// In CQRS a request is a command or query that interacts with a store like a database or a secret store.
	// Parameters for the query would be passed in the request constructor.
	// And be handled through the mediator, so the Handle function can execute the query.
	Request[T any] interface {
		// Handle executes the [Request].
		// The [context.Context] parameter and the return values can be accessed and altered in the [Pipeline].
		Handle(ctx context.Context, l *slog.Logger) (T, error)
	}
)

// Send a [Request] using a [Sender].
// This function uses reflect to decide the name of the request.
//
// The [Sender] interface is implemented by [Mediator].
func Send[T any](ctx context.Context, m Sender, req Request[T]) (T, error) {
	return SendWithLogger(ctx, m.getLogger(), m, req)
}

// SendWithLogger a [Request] with a logger instance using a [Sender].
// This function is like [Send], but with a logger parameter.
// Passing a logger can be useful if you want to add attributes to the logger in the caller.
//
// The [Sender] interface is implemented by [Mediator].
func SendWithLogger[T any](ctx context.Context, l *slog.Logger, m Sender, req Request[T]) (T, error) {
	handler := m.getRequestPipeline().Then(func(ctx context.Context, l *slog.Logger, _ Message) (any, error) {
		return req.Handle(ctx, l)
	})
	resp, err := handler.Handle(ctx, l, NewRequestMessage(req))
	respT := resp.(T)
	return respT, err
}
