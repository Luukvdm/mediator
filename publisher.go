package mediator

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

type (
	// NotificationHandler can receive notifications from the [Publisher].
	// NotificationHandlers can be subscribed to a [Notification] using the [Subscribe] function.
	NotificationHandler[T any] interface {
		Handle(ctx context.Context, l *slog.Logger, event T) error
	}

	notificationHandler[T any] struct {
		handleFunc func(ctx context.Context, l *slog.Logger, event T) error
	}

	// Publisher can be used to subscribe and publish notifications.
	// Publishing and subscribing is done through the [Publish] and [Subscribe] functions that take
	// Publisher as a parameter.
	//
	// Publisher also holds the [Pipeline] that messages pass through before they get handled.
	//
	// The interface is implemented by [Mediator].
	Publisher interface {
		getNotificationPipeline() Pipeline
		getLogger() *slog.Logger
		getAllNotifiers() map[any][]any
		newNotifier(key any, notifier any)
		getDefaultPublishOpts() *publishOptions
	}

	// Notification is an event that can be published through the [Publisher].
	Notification[T any] interface{}
)

func (nh notificationHandler[T]) Handle(ctx context.Context, l *slog.Logger, event T) error {
	return nh.handleFunc(ctx, l, event)
}

// NewNotificationHandler is a utility function for creating a [NotificationHandler] without having to define a type.
// This is especially useful when writing tests.
func NewNotificationHandler[T any](handleFunc func(ctx context.Context, l *slog.Logger, event T) error) NotificationHandler[T] {
	return notificationHandler[T]{
		handleFunc: handleFunc,
	}
}

// Subscribe to a [Notification] using [Publisher].
// When a [Notification] is published, every subscriber triggers the [Pipeline].
// So every subscriber in for the event makes the [Notification] go through the chain.
func Subscribe[T any](p Publisher, s NotificationHandler[T]) error {
	p.newNotifier(key[T]{}, s)
	return nil
}

// Publish a [Notification] using [Publisher].
//
// The [Publisher] interface is implemented by [Mediator].
func Publish[T Notification[any]](ctx context.Context, p Publisher, notification T, options ...PublishOption) error {
	return publish(ctx, p, notification, options...)
}

// PublishWithLogger publishes a [Notification] using [Publisher].
// This function is like [Publish], but with a logger parameter.
// Passing a logger can be useful if you want to add attributes to the logger in the caller.
//
// The [Publisher] interface is implemented by [Mediator].
//
// Deprecated: Use [Publish] with the [WithPublishLogger] publish option instead.
func PublishWithLogger[T Notification[any]](ctx context.Context, l *slog.Logger, p Publisher, notification T) error {
	return publish(ctx, p, notification, WithPublishLogger(l))
}

func publish[T Notification[any]](ctx context.Context, p Publisher, notification T, options ...PublishOption) error {
	opts := p.getDefaultPublishOpts()
	// overwrite default options with given options
	for _, o := range options {
		o(opts)
	}

	var handlers []NotificationHandler[T]
	for _, h := range p.getAllNotifiers()[key[T]{}] {
		handlers = append(handlers, h.(NotificationHandler[T]))
	}

	if len(handlers) == 0 {
		return nil
	}

	pl := p.getNotificationPipeline()

	var err error
	if opts.enableParallel {
		err = runHandlersParallel(ctx, opts.l, pl, notification, handlers)
	} else {
		err = runHandlersSerial(ctx, opts.l, pl, notification, handlers)
	}

	return err
}

func runHandlersSerial[T Notification[any]](ctx context.Context, l *slog.Logger, pl Pipeline, notification T, handlers []NotificationHandler[T]) error {
	var errs []error

	for _, h := range handlers {
		handlerPl := pl.Then(func(ctx context.Context, l *slog.Logger, _ Message) (any, error) {
			return nil, h.Handle(ctx, l, notification)
		})

		_, err := handlerPl.Handle(ctx, l, NewNotificationMessage[T](notification))
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func runHandlersParallel[T Notification[any]](ctx context.Context, l *slog.Logger, pl Pipeline, notification T, handlers []NotificationHandler[T]) error {
	var wg sync.WaitGroup
	wg.Add(len(handlers))

	var errs []error
	var errMu sync.Mutex

	for _, h := range handlers {
		// TODO: could add a goroutine pool option
		go func() {
			handlerPl := pl.Then(func(ctx context.Context, l *slog.Logger, _ Message) (any, error) {
				return nil, h.Handle(ctx, l, notification)
			})

			_, err := handlerPl.Handle(ctx, l, NewNotificationMessage[T](notification))
			if err != nil {
				errMu.Lock()
				errs = append(errs, err)
				errMu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
