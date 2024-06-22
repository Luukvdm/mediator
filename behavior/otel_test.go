package behavior_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/luukvdm/mediator"
	"github.com/luukvdm/mediator/behavior"
)

func TestTracer_Handler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	createBehav := func() (mediator.Behavior, *tracetest.InMemoryExporter) {
		inmemoryExp := tracetest.NewInMemoryExporter()
		traceProvider := sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
			sdkTrace.WithSpanProcessor(sdkTrace.NewSimpleSpanProcessor(inmemoryExp)))
		otel.SetTracerProvider(traceProvider)

		behav := behavior.NewOtelTracer(behavior.WithTracerProvider(traceProvider))
		return behav, inmemoryExp
	}

	t.Run("Attributes", func(t *testing.T) {
		t.Parallel()

		behav, inmemoryExp := createBehav()

		var reqCtx context.Context
		handler := fakeRequest{handleFunc: func(ctx context.Context, _ mediator.Message) (any, error) {
			reqCtx = ctx

			span := trace.SpanFromContext(reqCtx)
			assert.NotEmpty(t, span, "the tracer behavior should always create a span")
			assert.True(t, span.IsRecording())
			assert.True(t, span.SpanContext().HasTraceID())
			assert.True(t, span.SpanContext().HasSpanID())

			return nil, nil
		}}

		_, err := behav.Handler(handler).Handle(ctx, handler)
		require.NoError(t, err)

		span := trace.SpanFromContext(reqCtx)
		require.NotNil(t, span, "no span was created by the behavior")
		assert.False(t, span.IsRecording())

		expSpans := inmemoryExp.GetSpans()
		require.Len(t, expSpans, 1, "no span was created by the behavior")
		mySpan := expSpans[0]
		assert.Equal(t, "fakeRequest", mySpan.Name)
	})

	t.Run("RecordErrors", func(t *testing.T) {
		t.Parallel()

		behav, inmemoryExp := createBehav()

		reqErr := errors.New("something went wrong")
		var reqCtx context.Context
		handler := fakeRequest{handleFunc: func(ctx context.Context, _ mediator.Message) (any, error) {
			reqCtx = ctx

			span := trace.SpanFromContext(reqCtx)
			assert.NotEmpty(t, span, "the tracer behavior should always create a span")
			assert.True(t, span.IsRecording())
			assert.True(t, span.SpanContext().HasTraceID())
			assert.True(t, span.SpanContext().HasSpanID())

			return nil, reqErr
		}}

		_, err := behav.Handler(handler).Handle(ctx, handler)
		require.Error(t, err)

		expSpans := inmemoryExp.GetSpans()
		require.Len(t, expSpans, 1, "the behavior created an unexpected amount of spans")
		mySpan := expSpans[0]
		assert.Equal(t, "fakeRequest", mySpan.Name)
		assert.Equal(t, codes.Error, mySpan.Status.Code)
		assert.Equal(t, "an error occurred while processing this command", mySpan.Status.Description)
		assert.Equal(t, reqErr.Error(), mySpan.Attributes[1].Value.AsString())
	})
}
