package behavior

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"

	"github.com/luukvdm/mediator"
)

const instrumentationName = "github.com/luukvdm/mediator/behavior"

type (
	// OtelTracer is a [mediator.Behavior] that adds tracing to the chain.
	//
	// The behavior creates a new span for every request that passes through it and adds it to the context.
	// It will also adjust the status and add an error attribute to the span if the resulting error of the request is not nil.
	OtelTracer struct {
		tracer trace.Tracer
	}

	// OtelTracerOption defines the method to customize [NewOtelTracer].
	OtelTracerOption  func(*otelTracerOptions)
	otelTracerOptions struct {
		provider trace.TracerProvider
	}
)

// WithTracerProvider overwrites the [tracer.TracerProvider] that the [OtelTracer] [mediator.Behavior] uses.
func WithTracerProvider(provider trace.TracerProvider) OtelTracerOption {
	return func(o *otelTracerOptions) {
		o.provider = provider
	}
}

// Handler runs the [OtelTracer] behavior.
func (b *OtelTracer) Handler(next mediator.Handler) mediator.Handler {
	return mediator.HandlerFunc(func(ctx context.Context, msg mediator.Message) (any, error) {
		spanCtx, span := b.tracer.Start(ctx, msg.String())
		defer span.End()

		resp, err := next.Handle(spanCtx, msg)
		if err != nil {
			span.SetAttributes(semconv.ExceptionType(msg.String()))
			span.SetAttributes(semconv.ExceptionMessage(err.Error()))
			span.RecordError(err)
			span.SetStatus(codes.Error, "an error occurred while processing this command")
		}
		return resp, err
	})
}

// NewOtelTracer creates a new [OtelTracer] [mediator.Behavior].
func NewOtelTracer(opt ...OtelTracerOption) mediator.Behavior {
	// default options
	opts := &otelTracerOptions{
		provider: otel.GetTracerProvider(),
	}
	for _, o := range opt {
		o(opts)
	}
	return &OtelTracer{
		tracer: opts.provider.Tracer(instrumentationName),
	}
}
