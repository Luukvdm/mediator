package behavior

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/luukvdm/mediator"
)

// Slogger is a [mediator.Behavior] that adds logging to the chain.
//
// The logger behavior adds a `request` [slog.Attr] with the request name to the logger
// that is passed through it.
// It also logs the request after it is handled.
// This includes the time it took to handle the request and the error if it is not nil.
type Slogger struct {
	l *slog.Logger
}

// Handler runs the [Slogger] behavior.
func (b Slogger) Handler(next mediator.Handler) mediator.Handler {
	return mediator.HandlerFunc(func(ctx context.Context, l *slog.Logger, msg mediator.Message) (any, error) {
		l = l.With(msg.Type().String(), msg.String())

		start := time.Now()
		resp, err := next.Handle(ctx, l, msg)

		logArgs := []any{
			slog.Duration("elapsed", time.Since(start)),
		}

		if err != nil {
			logArgs = append(logArgs, slog.Any("error", err))
			l.ErrorContext(ctx, fmt.Sprintf("an error occurred while processing %s", msg.String()), logArgs...)
		} else {
			l.InfoContext(ctx, fmt.Sprintf("processed %s", msg.String()), logArgs...)
		}

		return resp, err
	})
}

// NewLogger creates a new [Slogger] [mediator.Behavior].
func NewLogger(l *slog.Logger) mediator.Behavior {
	return Slogger{
		l: l,
	}
}
