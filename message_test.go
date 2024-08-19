package mediator_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/luukvdm/mediator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestMessage(t *testing.T) {
	t.Parallel()

	name := "some-gopher"
	req := NewSearchGopherQuery(name)
	msg := mediator.NewRequestMessage(req)

	assert.Equal(t, req, msg.GetInner())
	assert.Equal(t, "searchGopher", msg.String())
	assert.Equal(t, mediator.TypeRequest, msg.Type())
	assert.IsType(t, searchGopher{}, msg.GetRequest())

	res, err := msg.Handle(context.Background(), slog.Default())
	require.NoError(t, err)

	assert.Equal(t, name, res.Name)
}

func TestNotificationMessage(t *testing.T) {
	t.Parallel()

	notification := mediator.Notification[searchGopher](searchGopher{})
	msg := mediator.NewNotificationMessage[searchGopher](notification)

	assert.Equal(t, notification, msg.GetInner())
	assert.Equal(t, "searchGopher", msg.String())
	assert.Equal(t, mediator.TypeNotification, msg.Type())
	assert.Equal(t, searchGopher{}, msg.GetNotification())
}
