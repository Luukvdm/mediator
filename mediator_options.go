package mediator

import "log/slog"

type (
	// Option defines the method to customize [Mediator].
	Option  func(*options)
	options struct {
		l                     *slog.Logger
		requestBehaviors      []Behavior
		requestPipeline       Pipeline
		notificationBehaviors []Behavior
		notificationPipeline  Pipeline
		parallelNotifications bool
	}
)

// WithLogger adds a custom logger that [Mediator] should use.
// This is also used as the default for request and notification middleware.
func WithLogger(l *slog.Logger) Option {
	return func(o *options) {
		o.l = l
	}
}

// WithRequestBehaviors adds the given [Behavior] slice to the [Request] [Pipeline].
func WithRequestBehaviors(behaviors ...Behavior) Option {
	return func(o *options) {
		o.requestBehaviors = behaviors
	}
}

// WithRequestPipeline overwrites the default [Pipeline] with the given implementation.
// If this option is set, other pipeline options like [WithRequestBehaviors] are ignored.
func WithRequestPipeline(pipeline Pipeline) Option {
	return func(o *options) {
		o.requestPipeline = pipeline
	}
}

// WithNotificationBehaviors adds behaviors to the [Notification] [Pipeline].
func WithNotificationBehaviors(behaviors ...Behavior) Option {
	return func(o *options) {
		o.notificationBehaviors = behaviors
	}
}

// WithNotificationPipeline overwrites the default [Pipeline] with the given implementation.
// If this option is set, other pipeline options like [WithNotificationBehaviors] are ignored.
func WithNotificationPipeline(pipeline Pipeline) Option {
	return func(o *options) {
		o.notificationPipeline = pipeline
	}
}

// WithParallelNotifications enables calling notification handlers in parallel by default.
// Can be overwritten by [WithParallelEnabled].
func WithParallelNotifications() Option {
	return func(o *options) {
		o.parallelNotifications = true
	}
}
