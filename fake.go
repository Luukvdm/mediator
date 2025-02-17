package mediator

import "log/slog"

type (
	fakeMediator struct {
		l                  *slog.Logger
		handleFunc         Handler
		notifiers          map[any][]any
		defaultPublishOpts *publishOptions
	}
)

func (f fakeMediator) Then(_ HandlerFunc) Handler {
	return f.handleFunc
}

func (f fakeMediator) getRequestPipeline() Pipeline {
	return f
}
func (f fakeMediator) getNotificationPipeline() Pipeline {
	return f
}

func (f fakeMediator) getLogger() *slog.Logger {
	return f.l
}

func (f fakeMediator) getAllNotifiers() map[any][]any {
	return f.notifiers
}

func (f fakeMediator) newNotifier(key any, notifier any) {
	f.notifiers[key] = append(f.notifiers[key], notifier)
}

func (f fakeMediator) getDefaultPublishOpts() *publishOptions {
	return f.defaultPublishOpts
}

// NewFake creates a [Mediator] object that can be used as a mock in tests.
//
// When sending [Request] through the fake mediator, it will pass it through the given [HandleFunc].
// This also works for [Notification], but the resp parameter will be ignored.
func NewFake(handleFunc HandlerFunc) Mediator {
	l := slog.Default()
	return fakeMediator{
		l:          l,
		handleFunc: handleFunc,
		notifiers:  make(map[any][]any),
		defaultPublishOpts: &publishOptions{
			l:              l,
			enableParallel: false,
		},
	}
}
