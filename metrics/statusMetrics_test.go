package metrics_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatusMetrics(t *testing.T) {
	t.Parallel()

	sm := metrics.NewStatusMetrics()
	require.False(t, check.IfNil(sm))
}

func TestStatusMetrics_AddRequestData(t *testing.T) {
	t.Parallel()

	t.Run("one metric exists for an endpoint", func(t *testing.T) {
		t.Parallel()

		sm := metrics.NewStatusMetrics()

		testOp, testDuration := "block_events", 1*time.Second
		sm.AddRequest(testOp, testDuration)

		res := sm.GetAll()
		require.Equal(t, res[testOp], &data.EndpointMetricsResponse{
			NumRequests:       1,
			TotalResponseTime: testDuration,
		})
	})

	t.Run("multiple entries exist for an endpoint", func(t *testing.T) {
		t.Parallel()

		sm := metrics.NewStatusMetrics()

		testOp := "block_events"
		testDuration0, testDuration1, testDuration2 := 4*time.Millisecond, 20*time.Millisecond, 2*time.Millisecond
		sm.AddRequest(testOp, testDuration0)
		sm.AddRequest(testOp, testDuration1)
		sm.AddRequest(testOp, testDuration2)

		res := sm.GetAll()
		require.Equal(t, res[testOp], &data.EndpointMetricsResponse{
			NumRequests:       3,
			TotalResponseTime: testDuration0 + testDuration1 + testDuration2,
		})
	})

	t.Run("multiple entries for multiple endpoints", func(t *testing.T) {
		t.Parallel()

		sm := metrics.NewStatusMetrics()

		testOp0, testOp1 := "block_events", "scrs_events"

		testDuration0End0, testDuration1End0 := time.Second, 5*time.Second
		testDuration0End1, testDuration1End1 := time.Hour, 4*time.Hour

		sm.AddRequest(testOp0, testDuration0End0)
		sm.AddRequest(testOp0, testDuration1End0)

		sm.AddRequest(testOp1, testDuration0End1)
		sm.AddRequest(testOp1, testDuration1End1)

		res := sm.GetAll()

		require.Len(t, res, 2)
		require.Equal(t, res[testOp0], &data.EndpointMetricsResponse{
			NumRequests:       2,
			TotalResponseTime: testDuration0End0 + testDuration1End0,
		})
		require.Equal(t, res[testOp1], &data.EndpointMetricsResponse{
			NumRequests:       2,
			TotalResponseTime: testDuration0End1 + testDuration1End1,
		})
	})
}

func TestStatusMetrics_GetMetricsForPrometheus(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		sm := metrics.NewStatusMetrics()

		testOperation := "block_events"
		testDuration0, testDuration1, testDuration2 := 4*time.Millisecond, 20*time.Millisecond, 2*time.Millisecond
		sm.AddRequest(testOperation, testDuration0)
		sm.AddRequest(testOperation, testDuration1)
		sm.AddRequest(testOperation, testDuration2)

		res := sm.GetMetricsForPrometheus()

		expectedString := `# TYPE num_requests counter
num_requests{operation="block_events"} 3

# TYPE total_response_time counter
total_response_time{operation="block_events"} 26

`

		require.Equal(t, expectedString, res)

	})
}

func TestStatusMetrics_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	sm := metrics.NewStatusMetrics()

	numIterations := 500
	wg := sync.WaitGroup{}
	wg.Add(numIterations)

	for i := 0; i < numIterations; i++ {
		go func(index int) {
			switch index % 3 {
			case 0:
				sm.AddRequest(fmt.Sprintf("op_%d", index%5), time.Hour*time.Duration(index))
			case 1:
				_ = sm.GetAll()
			case 2:
				_ = sm.GetMetricsForPrometheus()
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestStatusMetrics_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	sm := metrics.NewStatusMetrics()
	assert.False(t, sm.IsInterfaceNil())
}
