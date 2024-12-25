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

var _ mediator.Pipeline = (*customPipeline)(nil)

type customPipeline struct {
	isCalled bool
}

func (c *customPipeline) Then(hf mediator.HandlerFunc) mediator.Handler {
	c.isCalled = true
	return hf
}

func TestPipeline_RequestCustom(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pipeline := &customPipeline{}
	m := mediator.New(mediator.WithRequestPipeline(pipeline))

	req := mocks.NewMockRequest[string](t)
	req.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger")).Once().Return("test-123", nil)

	_, err := mediator.Send[string](ctx, m, req)
	require.NoError(t, err)
	assert.True(t, pipeline.isCalled, "custom pipeline is not used")
}

func TestPipeline_NotificationCustom(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pipeline := &customPipeline{}
	m := mediator.New(mediator.WithNotificationPipeline(pipeline))

	event := "some-event"
	handler := mocks.NewMockNotificationHandler[string](t)
	handler.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger"), mock.Anything).Once().Return(nil)
	require.NoError(t, mediator.Subscribe[string](m, handler))

	err := mediator.Publish(ctx, m, event)
	require.NoError(t, err)
	assert.True(t, pipeline.isCalled, "custom pipeline is not used")
}
