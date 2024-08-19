package mediator

import "log/slog"

type (
	fakeMediator struct {
		l          *slog.Logger
		handleFunc Handler
		notifiers  map[any][]any
	}
)

func (f fakeMediator) Then(_ HandlerFunc) Handler {
	return f.handleFunc
}

func (f fakeMediator) getPipeline() Pipeline {
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

// NewFake creates a [Mediator] object that can be used as a mock in tests.
//
// When sending [Request] through the fake mediator, it will pass it through the given [HandleFunc].
// This also works for [Notification], but the resp parameter will be ignored.
func NewFake(handleFunc HandlerFunc) Mediator {
	return fakeMediator{
		l:          slog.Default(),
		handleFunc: handleFunc,
		notifiers:  make(map[any][]any),
	}
}
