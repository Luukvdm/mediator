package mediator

import (
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
		pipeline    Pipeline
		notifiers   map[any][]any
		notifiersMu sync.RWMutex
	}
	key[T any] struct{}

	// Option defines the method to customize [Mediator].
	Option  func(*options)
	options struct {
		behaviors []Behavior
		pipeline  Pipeline
	}
)

// WithBehaviors adds the given [Behavior] slice to the pipeline.
// The behaviors are used for both [Request] and [Notification] messages.
func WithBehaviors(behaviors ...Behavior) Option {
	return func(o *options) {
		o.behaviors = behaviors
	}
}

// WithPipeline overwrites the default pipeline with the given implementation.
// If this option is set, other pipeline options like [WithBehaviors] are ignored.
func WithPipeline(pipeline Pipeline) Option {
	return func(o *options) {
		o.pipeline = pipeline
	}
}

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

func (m *mediator) getPipeline() Pipeline {
	return m.pipeline
}

// New creates a new [Mediator].
// The [Mediator] can be customized with the [Option] slice parameter.
func New(opt ...Option) Mediator {
	// default options
	opts := &options{
		behaviors: []Behavior{},
	}
	for _, o := range opt {
		o(opts)
	}
	if opts.pipeline == nil {
		opts.pipeline = newPipeline(opts.behaviors...)
	}

	return &mediator{
		pipeline:  opts.pipeline,
		notifiers: make(map[any][]any),
	}
}
