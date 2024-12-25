package mediator

import "log/slog"

type (
	// PublishOption defines the method to customize [Publisher].
	PublishOption  func(*publishOptions)
	publishOptions struct {
		l              *slog.Logger
		enableParallel bool
	}
)

// WithPublishLogger is used to add a custom logger instance to the publisher.
func WithPublishLogger(l *slog.Logger) PublishOption {
	return func(o *publishOptions) {
		o.l = l
	}
}

// WithParallelEnabled when enabled runs notification handlers in parallel.
// This currently creates a new goroutine for every handler.
// Can be enabled by default using [WithParallelNotifications].
func WithParallelEnabled(enabled bool) PublishOption {
	return func(o *publishOptions) {
		o.enableParallel = enabled
	}
}
