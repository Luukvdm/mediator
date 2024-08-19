package mediator_test

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/luukvdm/mediator"
	"github.com/luukvdm/mediator/mocks"
)

type testBehavior struct {
	counter        int
	passedRequests []string
	handleFunc     func(ctx context.Context, l *slog.Logger, msg mediator.Message, next mediator.Handler) (any, error)
}

func (b *testBehavior) Handler(next mediator.Handler) mediator.Handler {
	return mediator.HandlerFunc(func(ctx context.Context, l *slog.Logger, msg mediator.Message) (any, error) {
		b.counter++
		name := reflect.TypeOf(msg).Name()
		b.passedRequests = append(b.passedRequests, name)
		if b.handleFunc == nil {
			return next.Handle(ctx, l, msg)
		}
		return b.handleFunc(ctx, l, msg, next)
	})
}

func TestChain_ResultPassesThroughBehavior(t *testing.T) {
	t.Parallel()

	cases := []struct {
		result string
		err    error
	}{
		{
			result: "test-123",
			err:    nil,
		},
		{
			result: "",
			err:    errors.New("request failed"),
		},
	}

	for _, c := range cases {
		ctx := context.Background()

		var next mediator.Handler
		b := mocks.NewMockBehavior(t)
		b.On("Handler", mock.AnythingOfType("mediator.HandlerFunc")).
			Run(func(args mock.Arguments) {
				next = args.Get(0).(mediator.Handler)
			}).
			Return(mediator.HandlerFunc(func(ctx context.Context, l *slog.Logger, msg mediator.Message) (any, error) {
				res, err := next.Handle(ctx, l, msg)
				assert.Equal(t, c.result, res)
				assert.Equal(t, c.err, err)
				return res, err
			}))

		m := mediator.New(mediator.WithBehaviors(b))

		req := mocks.NewMockRequest[string](t)
		req.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger")).Return(c.result, c.err)

		res, err := mediator.Send[string](ctx, m, req)
		assert.Equal(t, c.result, res)
		assert.Equal(t, c.err, err)
	}
}

func TestChain_BehaviorOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var counter int
	b1 := &testBehavior{
		handleFunc: func(ctx context.Context, l *slog.Logger, msg mediator.Message, next mediator.Handler) (any, error) {
			assert.Equal(t, 0, counter, "expected behavior to be in a different order")
			counter++
			return next.Handle(ctx, l, msg)
		},
	}
	b2 := &testBehavior{
		handleFunc: func(ctx context.Context, l *slog.Logger, msg mediator.Message, next mediator.Handler) (any, error) {
			assert.Equal(t, 1, counter, "expected behavior to be in a different order")
			counter++
			return next.Handle(ctx, l, msg)
		},
	}

	m := mediator.New(mediator.WithBehaviors(b1, b2))

	req := mocks.NewMockRequest[string](t)
	req.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger")).Return("test-123", nil)

	// use SendAnon to trigger the behavior chain
	_, err := mediator.Send[string](ctx, m, req)
	require.NoError(t, err)
}
