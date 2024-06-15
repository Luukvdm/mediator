package behavior_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/luukvdm/mediator"
	"github.com/luukvdm/mediator/behavior"
)

type fakeRequest struct {
	handleFunc func(_ context.Context, _ mediator.Message) (any, error)
}

func (n fakeRequest) GetInner() any {
	return n
}

func (n fakeRequest) String() string {
	return reflect.TypeOf(n).Name()
}

func (n fakeRequest) Type() mediator.MessageType {
	return mediator.TypeRequest
}

func (n fakeRequest) Handle(ctx context.Context, msg mediator.Message) (any, error) {
	if n.handleFunc == nil {
		return nil, nil
	}
	return n.handleFunc(ctx, msg)
}

func TestLogger_Handler_Successful(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, nil))

	handler := fakeRequest{}

	behav := behavior.NewLogger(l)
	_, err := behav.Handler(handler).Handle(ctx, handler)
	require.NoError(t, err)

	lines := bytes.Split(buf.Bytes(), []byte{'\n'})
	// split lines contains an empty line
	assert.Len(t, lines, 2, "the logger behavior should have only one log lines")

	var m map[string]any
	err = json.Unmarshal(lines[0], &m)
	require.NoError(t, err)

	assert.NotEmpty(t, m["time"])
	assert.Equal(t, "INFO", m["level"])
	assert.Equal(t, "fakeRequest", m["request"])
	assert.Equal(t, "processed fakeRequest", m["msg"])
	assert.NotEmpty(t, m["elapsed"])
}

func TestLogger_Handler_Error(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, nil))

	reqErr := errors.New("something went wrong")
	handler := fakeRequest{handleFunc: func(_ context.Context, _ mediator.Message) (any, error) {
		return nil, reqErr
	}}

	behav := behavior.NewLogger(l)
	_, err := behav.Handler(handler).Handle(ctx, handler)
	require.Error(t, err)

	lines := bytes.Split(buf.Bytes(), []byte{'\n'})
	// split lines contains an empty line
	assert.Len(t, lines, 2, "the logger behavior should have only one log line")

	var m map[string]any
	err = json.Unmarshal(lines[0], &m)
	require.NoError(t, err)

	assert.NotEmpty(t, m["time"])
	assert.Equal(t, "ERROR", m["level"])
	assert.Equal(t, reqErr.Error(), m["error"])
	assert.Equal(t, "fakeRequest", m["request"])
	assert.Equal(t, "an error occurred while processing fakeRequest", m["msg"])
	assert.NotEmpty(t, m["elapsed"])
}
