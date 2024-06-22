package mediator

import (
	"context"
)

type (
	// Sender is used by the [Send] function.
	// It holds the [Pipeline] that messages pass through before they get handled.
	//
	// The interface is implemented by [Mediator].
	Sender interface {
		getPipeline() Pipeline
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
		Handle(ctx context.Context) (T, error)
	}
)

// Send a [Request] using a [Sender].
// This function uses reflect to decide the name of the request.
//
// The [Sender] interface is implemented by [Mediator].
func Send[T any](ctx context.Context, m Sender, req Request[T]) (T, error) {
	handler := m.getPipeline().Then(func(ctx context.Context, _ Message) (any, error) {
		return req.Handle(ctx)
	})
	resp, err := handler.Handle(ctx, newRequestMessage(req))
	respT := resp.(T)
	return respT, err
}
