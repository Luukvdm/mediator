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

func TestPipeline_Custom(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	pipeline := &customPipeline{}
	m := mediator.New(mediator.WithPipeline(pipeline))

	req := mocks.NewMockRequest[string](t)
	req.EXPECT().Handle(ctx, mock.AnythingOfType("*slog.Logger")).Return("test-123", nil)

	_, err := mediator.Send[string](ctx, m, req)
	require.NoError(t, err)
	assert.True(t, pipeline.isCalled, "custom pipeline is not used")
}
