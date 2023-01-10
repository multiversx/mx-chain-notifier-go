package ws_test

import (
	"errors"
	"io"
	"testing"

	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher/ws"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testWriter implements io.WriteCloser
type testWriter struct{}

// Write implements io.Writer interface for testing
func (tw *testWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

// Close implements io.Closer interface for testing
func (tw *testWriter) Close() error {
	return nil
}

func createMockWSDispatcherArgs() ws.ArgsWSDispatcher {
	args := ws.ArgsWSDispatcher{}

	args.Hub = &mocks.HubStub{}
	args.Conn = &mocks.WSConnStub{}
	return args
}

func TestNewWebSocketDispatcher(t *testing.T) {
	t.Parallel()

	t.Run("nil hub handler", func(t *testing.T) {
		t.Parallel()

		args := createMockWSDispatcherArgs()
		args.Hub = nil

		wd, err := ws.NewTestWSDispatcher(args)
		require.Nil(t, wd)
		assert.Equal(t, ws.ErrNilHubHandler, err)
	})

	t.Run("nil ws conn", func(t *testing.T) {
		t.Parallel()

		args := createMockWSDispatcherArgs()
		args.Conn = nil

		wd, err := ws.NewTestWSDispatcher(args)
		require.Nil(t, wd)
		assert.Equal(t, ws.ErrNilWSConn, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockWSDispatcherArgs()
		wd, err := ws.NewTestWSDispatcher(args)
		require.Nil(t, err)
		require.NotNil(t, wd)
	})
}

func TestWritePump(t *testing.T) {
	t.Parallel()

	args := createMockWSDispatcherArgs()

	numCalls := uint32(0)

	wasCalled := false
	args.Conn = &mocks.WSConnStub{
		NextWriterCalled: func(message int) (io.WriteCloser, error) {
			wasCalled = true

			if numCalls > 0 {
				return nil, errors.New("new error")
			}

			numCalls++
			return &testWriter{}, nil
		},
	}

	wd, err := ws.NewTestWSDispatcher(args)
	require.Nil(t, err)

	events := []data.Event{
		{
			Address: "addr1",
		},
	}

	wd.PushEvents(events)
	wd.PushEvents(events)

	wd.WritePump()

	assert.True(t, wasCalled)
}

func TestReadPump(t *testing.T) {
	t.Parallel()

	args := createMockWSDispatcherArgs()

	numCalls := uint32(0)

	wasCalled := false
	args.Conn = &mocks.WSConnStub{
		ReadMessageCalled: func() (messageType int, p []byte, err error) {
			wasCalled = true

			if numCalls > 0 {
				return 0, nil, errors.New("new error")
			}

			numCalls++
			return 0, nil, nil
		},
	}

	wd, err := ws.NewTestWSDispatcher(args)
	require.Nil(t, err)

	wd.ReadPump()

	assert.True(t, wasCalled)
}
