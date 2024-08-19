package mediator_test

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/luukvdm/mediator"
)

type GetByID struct {
	id int
}

func (g GetByID) Handle(_ context.Context, _ *slog.Logger) (string, error) {
	return fmt.Sprintf("object with id %d", g.id), nil
}

func ExampleNewFake() {
	ctx := context.Background()

	req := GetByID{5}

	// create a new fake mediator and give it the response you want it to return
	// m := mediator.NewFakeMediator[string]("mocked response", nil)
	m := mediator.NewFake(func(_ context.Context, _ *slog.Logger, msg mediator.Message) (any, error) {
		// both notifications and requests can pass through this function
		if req, ok := msg.(mediator.RequestMessage[string]); ok {
			q := req.GetRequest().(GetByID)
			// if this where a test, this is the place to put asserts
			// assert.Equal(t, req, q)
			fmt.Printf("mock got object: %v\n", q.id)
		}
		return "mocked response", nil
	})

	// this would happen in the function that you are testing
	// using the fake mediator it's possible to test it without executing the actual request
	response, err := mediator.Send[string](ctx, m, req)
	if err != nil {
		return
	}

	fmt.Println(response)

	// Output:
	// mock got object: 5
	// mocked response
}
