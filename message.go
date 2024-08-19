package mediator

import (
	"context"
	"log/slog"
	"reflect"
)

// MessageType is the type of [Message].
type MessageType int

const (
	// TypeRequest is the type for [Request] messages.
	TypeRequest MessageType = iota
	// TypeNotification is the type for [Notification] messages.
	TypeNotification
)

func (t MessageType) String() string {
	switch t {
	case TypeRequest:
		return "request"
	case TypeNotification:
		return "notification"
	default:
		return "unknown"
	}
}

type (
	// Message is a wrapper around the object that is passed through the [Pipeline].
	// The object is wrapped so the handlers in the [Pipeline] don't have to implement common functions
	// like getting the name or type of the object.
	Message interface {
		// String returns the name of the object wrapped in [Message].
		//
		// This functions uses the [reflect] package.
		// Reflect is only used when the String function is called,
		// and the result is cached for other calls in the future.
		String() string
		// Type returns the [MessageType] of this [Message].
		Type() MessageType
		// GetInner returns the object that is wrapped in the [Message] interface.
		GetInner() any
	}
	// RequestMessage extends the [Message] interface with the [Request] interface.
	// It is the request version of the [Message] interface.
	RequestMessage[T any] interface {
		Message
		Request[T]
		GetRequest() Request[T]
	}
	requestMessage[T any] struct {
		name string
		req  Request[T]
	}
	// NotificationMessage extends the [Message] interface with the [Notification] interface.
	// It is the notification version of the [Message] interface.
	NotificationMessage[T any] interface {
		Message
		Notification[T]
		GetNotification() Notification[T]
	}
	notificationMessage[T any] struct {
		name         string
		notification Notification[T]
	}
)

// RequestMessage implementation

func (r requestMessage[T]) GetInner() any {
	return r.req
}

func (r requestMessage[T]) Handle(ctx context.Context, l *slog.Logger) (T, error) {
	return r.req.Handle(ctx, l)
}

func (r requestMessage[T]) String() string {
	// don't do any unnecessary reflect calls
	if len(r.name) == 0 {
		r.name = reflect.TypeOf(r.req).Name()
	}
	return r.name
}

func (r requestMessage[T]) Type() MessageType {
	return TypeRequest
}

func (r requestMessage[T]) GetRequest() Request[T] {
	return r.req
}

// NewRequestMessage wraps a [Request] so it implements the [Message] and [RequestMessage] interfaces.
func NewRequestMessage[T any](req Request[T]) RequestMessage[T] {
	return requestMessage[T]{
		req: req,
	}
}

// NotificationMessage implementation

func (n notificationMessage[T]) GetInner() any {
	return n.notification
}

func (n notificationMessage[T]) String() string {
	// don't do any unnecessary reflect calls
	if len(n.name) == 0 {
		n.name = reflect.TypeOf(n.notification).Name()
	}
	return n.name
}

func (n notificationMessage[T]) Type() MessageType {
	return TypeNotification
}

func (n notificationMessage[T]) GetNotification() Notification[T] {
	return n.notification
}

// NewNotificationMessage wraps a [Notification] so it implements the [Message] and [NotificationMessage] interfaces.
func NewNotificationMessage[T any](notification Notification[T]) NotificationMessage[T] {
	return notificationMessage[T]{
		notification: notification,
	}
}
