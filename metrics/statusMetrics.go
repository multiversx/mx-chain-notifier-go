package metrics

import (
	"strings"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-notifier-go/data"
)

const (
	numRequestsPromMetric       = "num_requests"
	totalResponseTimePromMetric = "total_response_time"
)

type statusMetrics struct {
	operationMetrics    map[string]*data.EndpointMetricsResponse
	mutOperationMetrics sync.RWMutex
}

// NewStatusMetrics will return an instance of the statusMetrics
func NewStatusMetrics() *statusMetrics {
	return &statusMetrics{
		operationMetrics: make(map[string]*data.EndpointMetricsResponse),
	}
}

// AddRequest will add the received data to the metrics map
func (sm *statusMetrics) AddRequest(path string, duration time.Duration) {
	sm.mutOperationMetrics.Lock()
	defer sm.mutOperationMetrics.Unlock()

	currentData := sm.operationMetrics[path]
	if currentData == nil {
		newMetricsData := &data.EndpointMetricsResponse{
			NumRequests:       1,
			TotalResponseTime: duration,
		}

		sm.operationMetrics[path] = newMetricsData

		return
	}

	currentData.NumRequests++
	currentData.TotalResponseTime += duration
}

// GetAll returns the metrics map
func (sm *statusMetrics) GetAll() map[string]*data.EndpointMetricsResponse {
	sm.mutOperationMetrics.RLock()
	defer sm.mutOperationMetrics.RUnlock()

	return sm.getAll()
}

func (sm *statusMetrics) getAll() map[string]*data.EndpointMetricsResponse {
	newMap := make(map[string]*data.EndpointMetricsResponse)
	for key, value := range sm.operationMetrics {
		newMap[key] = value
	}

	return newMap
}

// GetMetricsForPrometheus returns the metrics in a prometheus format
func (sm *statusMetrics) GetMetricsForPrometheus() string {
	sm.mutOperationMetrics.RLock()
	defer sm.mutOperationMetrics.RUnlock()

	metricsMap := sm.getAll()

	stringBuilder := strings.Builder{}

	for endpointPath, endpointData := range metricsMap {
		stringBuilder.WriteString(requestsCounterMetric(numRequestsPromMetric, endpointPath, endpointData.NumRequests))
		stringBuilder.WriteString(requestsCounterMetric(totalResponseTimePromMetric, endpointPath, uint64(endpointData.TotalResponseTime.Milliseconds())))
	}

	return stringBuilder.String()
}

// IsInterfaceNil returns true if there is no value under the interface
func (sm *statusMetrics) IsInterfaceNil() bool {
	return sm == nil
}
