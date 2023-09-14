package data

import (
	"time"
)

// EndpointMetricsResponse defines the response for status metrics endpoint
type EndpointMetricsResponse struct {
	NumRequests       uint64        `json:"num_requests"`
	TotalResponseTime time.Duration `json:"total_response_time"`
}
