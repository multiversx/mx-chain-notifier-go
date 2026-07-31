package rabbitmq

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

const testRetryInterval = time.Millisecond * 10

var errDialFailed = errors.New("dial failed")

// createTestClient creates a client which is not connected to any rabbitmq instance, with
// a dial function which always fails and signals each attempt on the returned channel
func createTestClient() (*rabbitMqClient, *int32, chan struct{}) {
	numDialCalls := int32(0)
	dialCalled := make(chan struct{}, 1)

	rc := &rabbitMqClient{
		url:           "amqp://guest:guest@localhost:5672/",
		retryInterval: testRetryInterval,
		closeChan:     make(chan struct{}),
		dialFunc: func(url string) (*amqp.Connection, error) {
			atomic.AddInt32(&numDialCalls, 1)

			select {
			case dialCalled <- struct{}{}:
			default:
			}

			return nil, errDialFailed
		},
	}

	return rc, &numDialCalls, dialCalled
}

func requireReturns(t *testing.T, description string, handler func()) {
	done := make(chan struct{})
	go func() {
		handler()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second * 10):
		require.Fail(t, description+" did not return")
	}
}

func requireDialAttempted(t *testing.T, dialCalled chan struct{}) {
	select {
	case <-dialCalled:
	case <-time.After(time.Second * 10):
		require.Fail(t, "no dial attempt was made")
	}
}

func TestNewRabbitMQClient(t *testing.T) {
	t.Parallel()

	t.Run("invalid url, should fail", func(t *testing.T) {
		t.Parallel()

		rc, err := NewRabbitMQClient("invalid url")
		require.Error(t, err)
		require.Nil(t, rc)
	})
}

func TestRabbitMqClient_ExchangeDeclare(t *testing.T) {
	t.Parallel()

	t.Run("no channel available, should fail", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		err := rc.ExchangeDeclare("allevents", "fanout")
		require.Equal(t, ErrClientClosed, err)
		require.Empty(t, rc.exchanges)
	})
}

func TestRabbitMqClient_SaveExchangeDeclaration(t *testing.T) {
	t.Parallel()

	t.Run("should save each exchange only once", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		rc.saveExchangeDeclaration("allevents", "fanout")
		rc.saveExchangeDeclaration("revert", "fanout")
		rc.saveExchangeDeclaration("allevents", "fanout")

		expectedExchanges := []exchangeDeclaration{
			{name: "allevents", kind: "fanout"},
			{name: "revert", kind: "fanout"},
		}
		require.Equal(t, expectedExchanges, rc.exchanges)
	})

	t.Run("should update the exchange type", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		rc.saveExchangeDeclaration("allevents", "fanout")
		rc.saveExchangeDeclaration("allevents", "direct")

		expectedExchanges := []exchangeDeclaration{
			{name: "allevents", kind: "direct"},
		}
		require.Equal(t, expectedExchanges, rc.exchanges)
	})
}

func TestRabbitMqClient_RedeclareExchanges(t *testing.T) {
	t.Parallel()

	t.Run("no saved exchange, should not fail", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		err := rc.redeclareExchanges()
		require.Nil(t, err)
	})
}

func TestRabbitMqClient_Publish(t *testing.T) {
	t.Parallel()

	t.Run("closed client, should fail without trying to reconnect", func(t *testing.T) {
		t.Parallel()

		rc, numDialCalls, _ := createTestClient()
		rc.Close()

		err := rc.Publish("allevents", "", true, false, amqp.Publishing{})
		require.Equal(t, ErrClientClosed, err)
		require.Equal(t, int32(0), atomic.LoadInt32(numDialCalls))
	})

	t.Run("no channel available, should fail instead of panicking", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		err := rc.Publish("allevents", "", true, false, amqp.Publishing{})
		require.Equal(t, ErrClientClosed, err)
	})
}

func TestRabbitMqClient_Reconnect(t *testing.T) {
	t.Parallel()

	t.Run("should retry until the client is closed", func(t *testing.T) {
		t.Parallel()

		rc, numDialCalls, dialCalled := createTestClient()

		done := make(chan struct{})
		go func() {
			rc.Reconnect()
			close(done)
		}()

		requireDialAttempted(t, dialCalled)
		rc.Close()

		select {
		case <-done:
		case <-time.After(time.Second * 10):
			require.Fail(t, "Reconnect did not stop after the client was closed")
		}

		require.True(t, atomic.LoadInt32(numDialCalls) >= 1)
	})

	t.Run("already closed client, should not dial", func(t *testing.T) {
		t.Parallel()

		rc, numDialCalls, _ := createTestClient()
		rc.Close()

		requireReturns(t, "Reconnect", rc.Reconnect)
		require.Equal(t, int32(0), atomic.LoadInt32(numDialCalls))
	})
}

func TestRabbitMqClient_ReopenChannel(t *testing.T) {
	t.Parallel()

	// a channel cannot be opened on a closed connection, so retrying to open the channel
	// would loop forever. It has to fall back on re-establishing the connection.
	t.Run("connection is down, should reconnect", func(t *testing.T) {
		t.Parallel()

		rc, _, dialCalled := createTestClient()

		done := make(chan struct{})
		go func() {
			rc.ReopenChannel()
			close(done)
		}()

		requireDialAttempted(t, dialCalled)
		rc.Close()

		select {
		case <-done:
		case <-time.After(time.Second * 10):
			require.Fail(t, "ReopenChannel did not stop after the client was closed")
		}
	})

	t.Run("already closed client, should not dial", func(t *testing.T) {
		t.Parallel()

		rc, numDialCalls, _ := createTestClient()
		rc.Close()

		requireReturns(t, "ReopenChannel", rc.ReopenChannel)
		require.Equal(t, int32(0), atomic.LoadInt32(numDialCalls))
	})
}

func TestRabbitMqClient_RecoverConnection(t *testing.T) {
	t.Parallel()

	t.Run("connection is down, should reconnect", func(t *testing.T) {
		t.Parallel()

		rc, _, dialCalled := createTestClient()

		done := make(chan struct{})
		go func() {
			rc.recoverConnection()
			close(done)
		}()

		requireDialAttempted(t, dialCalled)
		rc.Close()

		select {
		case <-done:
		case <-time.After(time.Second * 10):
			require.Fail(t, "recoverConnection did not stop after the client was closed")
		}
	})
}

func TestRabbitMqClient_ConnectionIsDown(t *testing.T) {
	t.Parallel()

	rc, _, _ := createTestClient()
	require.True(t, rc.connectionIsDown())
}

func TestRabbitMqClient_OpenChannel(t *testing.T) {
	t.Parallel()

	t.Run("no connection available, should fail", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		err := rc.tryOpenChannel()
		require.Equal(t, ErrClientClosed, err)
	})
}

func TestRabbitMqClient_Connect(t *testing.T) {
	t.Parallel()

	t.Run("dial fails, should return error", func(t *testing.T) {
		t.Parallel()

		rc, numDialCalls, _ := createTestClient()

		err := rc.tryConnect()
		require.Equal(t, errDialFailed, err)
		require.Equal(t, int32(1), atomic.LoadInt32(numDialCalls))
	})
}

func TestRabbitMqClient_WaitBeforeRetry(t *testing.T) {
	t.Parallel()

	t.Run("open client, should wait for the retry interval", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		startTime := time.Now()
		require.False(t, rc.waitBeforeRetry())
		require.True(t, time.Since(startTime) >= testRetryInterval)
	})

	t.Run("closed client, should not wait", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()
		rc.retryInterval = time.Minute
		rc.Close()

		startTime := time.Now()
		require.True(t, rc.waitBeforeRetry())
		require.True(t, time.Since(startTime) < time.Second)
	})
}

func TestRabbitMqClient_Close(t *testing.T) {
	t.Parallel()

	t.Run("should be idempotent", func(t *testing.T) {
		t.Parallel()

		rc, _, _ := createTestClient()

		require.False(t, rc.isClosed())

		rc.Close()
		rc.Close()

		require.True(t, rc.isClosed())
	})
}

func TestRabbitMqClient_ErrChannels(t *testing.T) {
	t.Parallel()

	rc, _, _ := createTestClient()

	require.Nil(t, rc.ConnErrChan())
	require.Nil(t, rc.CloseErrChan())

	connErrCh := make(chan *amqp.Error, notifyChanBufferSize)
	chanErr := make(chan *amqp.Error, notifyChanBufferSize)
	rc.connErrCh = connErrCh
	rc.chanErr = chanErr

	require.Equal(t, connErrCh, rc.ConnErrChan())
	require.Equal(t, chanErr, rc.CloseErrChan())
}

func TestRabbitMqClient_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	rc, _, dialCalled := createTestClient()

	numGoRoutines := 20
	wg := sync.WaitGroup{}
	wg.Add(numGoRoutines)

	for i := 0; i < numGoRoutines; i++ {
		go func(idx int) {
			defer wg.Done()

			switch idx % 5 {
			case 0:
				_ = rc.Publish("allevents", "", true, false, amqp.Publishing{})
			case 1:
				rc.Reconnect()
			case 2:
				rc.ReopenChannel()
			case 3:
				_ = rc.ExchangeDeclare("allevents", "fanout")
			case 4:
				_ = rc.ConnErrChan()
			}
		}(i)
	}

	requireDialAttempted(t, dialCalled)
	rc.Close()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second * 10):
		require.Fail(t, "not all operations returned after the client was closed")
	}
}

func TestRabbitMqClient_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var rc *rabbitMqClient
	require.True(t, rc.IsInterfaceNil())

	rc, _, _ = createTestClient()
	require.False(t, rc.IsInterfaceNil())
}
