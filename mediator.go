package mediator

import (
	"log/slog"
	"sync"
)

type (
	// Mediator is mainly used through the [Send] and [Publish] functions.
	// It holds the [Pipeline] that messages pass through before they get handled.
	//
	// Mediator implements the [Publisher] and [Sender] interfaces.
	Mediator interface {
		Publisher
		Sender
	}
	mediator struct {
		l                    *slog.Logger
		requestPipeline      Pipeline
		notificationPipeline Pipeline
		defaultPublishOpts   *publishOptions
		notifiers            map[any][]any
		notifiersMu          sync.RWMutex
	}
	key[T any] struct{}
)

func (m *mediator) newNotifier(key any, notifier any) {
	m.notifiersMu.Lock()
	m.notifiers[key] = append(m.notifiers[key], notifier)
	m.notifiersMu.Unlock()
}

func (m *mediator) getAllNotifiers() map[any][]any {
	m.notifiersMu.RLock()
	defer m.notifiersMu.RUnlock()
	return m.notifiers
}

func (m *mediator) getRequestPipeline() Pipeline {
	return m.requestPipeline
}

func (m *mediator) getNotificationPipeline() Pipeline {
	return m.notificationPipeline
}

func (m *mediator) getLogger() *slog.Logger {
	return m.l
}

// getPublishOpts return the default publish options.
func (m *mediator) getDefaultPublishOpts() *publishOptions {
	return m.defaultPublishOpts
}

// New creates a new [Mediator].
// The [Mediator] can be customized with the [Option] slice parameter.
func New(opt ...Option) Mediator {
	// default options
	opts := &options{
		requestBehaviors:      []Behavior{},
		notificationBehaviors: []Behavior{},
		l:                     slog.Default(),
		parallelNotifications: false,
	}
	for _, o := range opt {
		o(opts)
	}
	if opts.requestPipeline == nil {
		opts.requestPipeline = newPipeline(opts.requestBehaviors...)
	}
	if opts.notificationPipeline == nil {
		opts.notificationPipeline = newPipeline(opts.notificationBehaviors...)
	}

	return &mediator{
		l:                    opts.l,
		requestPipeline:      opts.requestPipeline,
		notificationPipeline: opts.notificationPipeline,
		notifiers:            make(map[any][]any),
		defaultPublishOpts: &publishOptions{
			l:              opts.l,
			enableParallel: opts.parallelNotifications,
		},
	}
}
