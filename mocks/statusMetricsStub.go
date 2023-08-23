package mocks

import (
	"time"

	"github.com/multiversx/mx-chain-notifier-go/data"
)

// StatusMetricsStub -
type StatusMetricsStub struct {
	AddRequestCalled              func(path string, duration time.Duration)
	GetAllCalled                  func() map[string]*data.EndpointMetricsResponse
	GetMetricsForPrometheusCalled func() string
}

// AddRequest -
func (s *StatusMetricsStub) AddRequest(path string, duration time.Duration) {
	if s.AddRequestCalled != nil {
		s.AddRequestCalled(path, duration)
	}
}

// GetAll -
func (s *StatusMetricsStub) GetAll() map[string]*data.EndpointMetricsResponse {
	if s.GetAllCalled != nil {
		return s.GetAllCalled()
	}

	return nil
}

// GetMetricsForPrometheus -
func (s *StatusMetricsStub) GetMetricsForPrometheus() string {
	if s.GetMetricsForPrometheusCalled != nil {
		return s.GetMetricsForPrometheusCalled()
	}

	return ""
}

// IsInterfaceNil -
func (s *StatusMetricsStub) IsInterfaceNil() bool {
	return s == nil
}
