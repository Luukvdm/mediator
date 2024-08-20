package mediator_test

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/luukvdm/mediator"
)

type (
	searchGopher struct {
		name string
	}
	GopherCreatedEvent struct {
		gopher Gopher
	}
	gopherCounter struct {
		count int
	}
	Gopher struct {
		Name          string
		Color         string
		CutenessLevel int
	}
	ExampleLogger struct{}
)

func (g searchGopher) Handle(_ context.Context, _ *slog.Logger) (Gopher, error) {
	return Gopher{
		Name:          g.name,
		Color:         "blue",
		CutenessLevel: 100,
	}, nil
}

func NewSearchGopherQuery(name string) mediator.Request[Gopher] {
	return searchGopher{name: name}
}

func (g gopherCounter) Handle(_ context.Context, _ *slog.Logger, event GopherCreatedEvent) error {
	g.count++
	fmt.Printf("counted gopher nr %d named %s\n", g.count, event.gopher.Name)
	return nil
}

func NewCreatedGophersCounter() mediator.NotificationHandler[GopherCreatedEvent] {
	return gopherCounter{}
}

func (b ExampleLogger) Handler(next mediator.Handler) mediator.Handler {
	return mediator.HandlerFunc(func(ctx context.Context, l *slog.Logger, msg mediator.Message) (any, error) {
		resp, err := next.Handle(ctx, l, msg)

		fmt.Printf("Request: %s request=%v\n", msg.String(), msg.GetInner())
		return resp, err
	})
}

// NewLogger creates a new [Logger] [mediator.Behavior].
func NewExampleLogger() mediator.Behavior {
	return ExampleLogger{}
}

func ExampleSend() {
	ctx := context.Background()

	// create a new mediator with a behavior that logs the request
	m := mediator.New(mediator.WithRequestBehaviors(NewExampleLogger()))

	// create a request and handle it through the mediator
	req := NewSearchGopherQuery("Gus")
	gopher, err := mediator.Send(ctx, m, req)
	if err != nil {
		return
	}
	fmt.Printf("found a gopher: %v\n", gopher)

	// Output:
	// Request: searchGopher request={Gus}
	// found a gopher: {Gus blue 100}
}

func ExamplePublish() {
	ctx := context.Background()

	// create a new mediator
	m := mediator.New()

	// create a new notification handler and subscribe it
	myHandler := NewCreatedGophersCounter()
	err := mediator.Subscribe(m, myHandler)
	if err != nil {
		return
	}

	// create a new notification
	created := Gopher{Name: "Gus", Color: "blue", CutenessLevel: 100}
	notification := GopherCreatedEvent{gopher: created}

	// publish it
	err = mediator.Publish(ctx, m, notification)
	if err != nil {
		return
	}

	// Output:
	// counted gopher nr 1 named Gus
}
