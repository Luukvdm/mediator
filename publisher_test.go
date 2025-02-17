package mediator_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/luukvdm/mediator"
	"github.com/luukvdm/mediator/mocks"
)

func TestSubscribe(t *testing.T) {
	t.Parallel()

	p := mediator.New()

	handler := mocks.NewMockNotificationHandler[string](t)
	err := mediator.Subscribe[string](p, handler)
	require.NoError(t, err)
}

func TestSubscribe_Multiple(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := mediator.New()
	myEvent := "test-123"

	for i := 0; i < 5; i++ {
		handler := mocks.NewMockNotificationHandler[string](t)
		handler.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger"), mock.Anything).
			Once().
			Return(nil)
		err := mediator.Subscribe[string](p, handler)
		require.NoError(t, err)
	}

	err := mediator.Publish(ctx, p, myEvent)
	require.NoError(t, err)
}

func TestPublish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := mediator.New()
	myEvent := "test-123"

	handler := mocks.NewMockNotificationHandler[string](t)
	handler.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger"), mock.Anything).
		Once().
		Return(nil)
	err := mediator.Subscribe[string](p, handler)
	require.NoError(t, err)

	err = mediator.Publish(ctx, p, myEvent)
	require.NoError(t, err)
}

func TestPublish_Errors(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := mediator.New()

	myEvent := "test"
	myErr := errors.New("fake error")
	handler := mocks.NewMockNotificationHandler[string](t)
	handler.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger"), mock.Anything).
		Twice().
		Return(myErr)
	err := mediator.Subscribe[string](p, handler)
	require.NoError(t, err)
	err = mediator.Subscribe[string](p, handler)
	require.NoError(t, err)

	err = mediator.Publish(ctx, p, myEvent)
	require.Error(t, err)
	require.Equal(t, errors.Join(myErr, myErr).Error(), err.Error())
}

func TestPublish_Parallel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := mediator.New(mediator.WithParallelNotifications())

	myEvent := "some-event"
	handlerCount := 5
	var wg sync.WaitGroup
	wg.Add(handlerCount)
	for range handlerCount {
		handler := mocks.NewMockNotificationHandler[string](t)
		handler.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger"), mock.Anything).
			Once().
			Run(func(_ mock.Arguments) {
				wg.Done()
			}).
			Return(nil)
		err := mediator.Subscribe(p, handler)
		require.NoError(t, err)
	}

	err := mediator.Publish(ctx, p, myEvent)
	require.NoError(t, err)

	wg.Wait()
}

func TestPublish_ParallelErrors(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := mediator.New()

	myEvent := "some-event"
	handlerCount := 5
	var wg sync.WaitGroup
	wg.Add(handlerCount)
	for i := range handlerCount {
		handler := mocks.NewMockNotificationHandler[string](t)
		handler.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger"), mock.Anything).
			Once().
			Run(func(_ mock.Arguments) {
				wg.Done()
			}).
			Return(fmt.Errorf("handler %d errored out", i))
		err := mediator.Subscribe(p, handler)
		require.NoError(t, err)
	}

	err := mediator.Publish(ctx, p, myEvent, mediator.WithParallelEnabled(true))
	require.Error(t, err)

	for i := range handlerCount {
		assert.ErrorContains(t, err, fmt.Sprintf("handler %d errored out", i))
	}

	wg.Wait()
}

func TestPublish_NoHandlers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := mediator.New()

	err := mediator.Publish(ctx, p, "some-event")
	require.NoError(t, err)
}

func TestPublish_WithLogger(t *testing.T) {
	t.Parallel()

	t.Run("publish_logger", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))
		p := mediator.New()

		myEvent := "some-event"
		handler := mediator.NewNotificationHandler(func(_ context.Context, l *slog.Logger, event string) error {
			assert.Equal(t, myEvent, event)
			assert.Exactly(t, logger, l)
			return nil
		})
		require.NoError(t, mediator.Subscribe(p, handler))
		err := mediator.Publish(ctx, p, myEvent, mediator.WithPublishLogger(logger))
		require.NoError(t, err)
	})

	t.Run("default_logger", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))
		p := mediator.New(mediator.WithLogger(logger))

		myEvent := "some-event"
		handler := mediator.NewNotificationHandler(func(_ context.Context, l *slog.Logger, event string) error {
			assert.Equal(t, myEvent, event)
			assert.Exactly(t, logger, l)
			return nil
		})
		require.NoError(t, mediator.Subscribe(p, handler))
		err := mediator.Publish(ctx, p, myEvent)
		require.NoError(t, err)
	})

}

func TestPublish_InlineSubscriber(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := mediator.New()

	myEvent := "some-event"
	var eventHandled bool
	handler := mediator.NewNotificationHandler(func(_ context.Context, _ *slog.Logger, event string) error {
		assert.Equal(t, myEvent, event)
		eventHandled = true
		return nil
	})
	err := mediator.Subscribe(p, handler)
	require.NoError(t, err)

	err = mediator.Publish(ctx, p, myEvent)
	require.NoError(t, err)

	assert.True(t, eventHandled, "inline handler should be notified about the event")
}

func TestPublish_BehaviorPersistence(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	behav := &testBehavior{}
	m := mediator.New(mediator.WithNotificationBehaviors(behav))
	myEvent := "test"

	handler := mocks.NewMockNotificationHandler[string](t)
	handler.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger"), mock.Anything).
		Return(nil)
	err := mediator.Subscribe[string](m, handler)
	require.NoError(t, err)

	rounds := 5
	for i := 0; i < rounds; i++ {
		err := mediator.Publish[string](ctx, m, myEvent)
		require.NoError(t, err)
	}
	assert.Equal(t, rounds, behav.counter, "publish seems to use copies of the behavior instead of reusing them (or not using them at all)")
}
