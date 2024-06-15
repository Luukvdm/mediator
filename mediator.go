package mediator

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
		behaviorChain chain
		notifiers     map[any][]any
	}
	key[T any] struct{}

	// Option defines the method to customize [Mediator].
	Option  func(*options)
	options struct {
		behaviors []Behavior
	}
)

// WithBehaviors adds the given [Behavior] slice to the pipeline.
// The behaviors are used for both [Request] and [Notification] messages.
func WithBehaviors(behaviors ...Behavior) Option {
	return func(o *options) {
		o.behaviors = behaviors
	}
}

func (m *mediator) newNotifier(key any, notifier any) {
	m.notifiers[key] = append(m.notifiers[key], notifier)
}

func (m *mediator) getAllNotifiers() map[any][]any {
	return m.notifiers
}

func (m *mediator) getChain() chain {
	return m.behaviorChain
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

	return &mediator{
		behaviorChain: newChain(opts.behaviors...),
		notifiers:     make(map[any][]any),
	}
}
