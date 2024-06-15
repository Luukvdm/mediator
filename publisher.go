package mediator

import (
	"context"
	"errors"
)

type (
	// NotificationHandler can receive notifications from the [Publisher].
	// NotificationHandlers can be subscribed to a [Notification] using the [Subscribe] function.
	NotificationHandler[T any] interface {
		Handle(ctx context.Context, event T) error
	}

	// Publisher can be used to subscribe and publish notifications.
	// Publishing and subscribing is done through the [Publish] and [Subscribe] functions that take
	// Publisher as a parameter.
	//
	// Publisher also holds the [Pipeline] that messages pass through before they get handled.
	//
	// The interface is implemented by [Mediator].
	Publisher interface {
		getChain() chain
		getAllNotifiers() map[any][]any
		newNotifier(key any, notifier any)
	}

	// Notification is an event that can be published through the [Publisher].
	Notification[T any] interface{}
)

// newNotificationHandler creates the final 'behavior' for the [chain].
// It calls the handle func on the [NotificationHandler].
func newNotificationHandler[T Notification[any]](notification T, handler NotificationHandler[T]) Handler {
	return HandlerFunc(func(ctx context.Context, _ Message) (any, error) {
		return nil, handler.Handle(ctx, notification)
	})
}

// Subscribe to a [Notification] using [Publisher].
// When a [Notification] is published, every subscriber triggers the [Pipeline].
// So every subscriber in for the event makes the [Notification] go through the chain.
func Subscribe[T any](p Publisher, s NotificationHandler[T]) error {
	p.newNotifier(key[T]{}, s)
	return nil
}

// Publish a [Notification] using [Publisher].
func Publish[T Notification[any]](ctx context.Context, p Publisher, notification T) error {
	allHandlers := p.getAllNotifiers()
	handlers := allHandlers[key[T]{}]

	var errs []error
	for _, handler := range handlers {
		h, ok := handler.(NotificationHandler[T])
		if !ok {
			// This shouldn't happen, but catching it just in case to prevent possible panics
			errs = append(errs, errors.New("subscribers contain a broken handler that doesn't implement the NotificationHandler interface"))
		}
		handler := p.getChain().Then(newNotificationHandler(notification, h))

		_, err := handler.Handle(ctx, newNotificationMessage[T](notification))
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
