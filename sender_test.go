package mediator_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/luukvdm/mediator"
	"github.com/luukvdm/mediator/mocks"
)

func TestSend(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	behav := &testBehavior{}
	m := mediator.New(mediator.WithBehaviors(behav))

	msg := "test123"
	req := mocks.NewMockRequest[string](t)
	req.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger")).Return(msg, nil)

	resp, err := mediator.Send[string](ctx, m, req)
	require.NoError(t, err)
	assert.Equal(t, msg, resp)
	// check if the behavior was used
	assert.Equal(t, 1, behav.counter, "middleware called unexpected amount of times")
}

func TestSend_BehaviorPersistence(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	behav := &testBehavior{}
	m := mediator.New(mediator.WithBehaviors(behav))

	req := mocks.NewMockRequest[string](t)
	req.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger")).Return("test-123", nil)

	rounds := 5
	for i := 0; i < rounds; i++ {
		_, err := mediator.Send[string](ctx, m, req)
		require.NoError(t, err)
	}
	assert.Equal(t, rounds, behav.counter, "send seems to use copies of the behavior instead of reusing them (or not using them at all)")
}
