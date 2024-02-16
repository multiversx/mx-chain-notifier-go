package mocks

import (
	"net/http"

	"github.com/multiversx/mx-chain-notifier-go/data"
)

// FacadeStub implements FacadeHandler interface
type FacadeStub struct {
	HandlePushEventsCalled        func(events data.ArgsSaveBlockData) error
	HandleRevertEventsCalled      func(events data.RevertBlock)
	HandleFinalizedEventsCalled   func(events data.FinalizedBlock)
	ServeCalled                   func(w http.ResponseWriter, r *http.Request)
	GetConnectorUserAndPassCalled func() (string, string)
	GetMetricsCalled              func() map[string]*data.EndpointMetricsResponse
	GetMetricsForPrometheusCalled func() string
}

// HandlePushEvents -
func (fs *FacadeStub) HandlePushEvents(events data.ArgsSaveBlockData) error {
	if fs.HandlePushEventsCalled != nil {
		return fs.HandlePushEventsCalled(events)
	}

	return nil
}

// HandleRevertEvents -
func (fs *FacadeStub) HandleRevertEvents(events data.RevertBlock) {
	if fs.HandleRevertEventsCalled != nil {
		fs.HandleRevertEventsCalled(events)
	}
}

// HandleFinalizedEvents -
func (fs *FacadeStub) HandleFinalizedEvents(events data.FinalizedBlock) {
	if fs.HandleFinalizedEventsCalled != nil {
		fs.HandleFinalizedEventsCalled(events)
	}
}

// ServeHTTP -
func (fs *FacadeStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fs.ServeCalled != nil {
		fs.ServeCalled(w, r)
	}
}

// GetConnectorUserAndPass -
func (fs *FacadeStub) GetConnectorUserAndPass() (string, string) {
	if fs.GetConnectorUserAndPassCalled != nil {
		return fs.GetConnectorUserAndPassCalled()
	}

	return "", ""
}

// GetMetrics -
func (fs *FacadeStub) GetMetrics() map[string]*data.EndpointMetricsResponse {
	if fs.GetMetricsCalled != nil {
		return fs.GetMetricsCalled()
	}

	return nil
}

// GetMetricsForPrometheus -
func (fs *FacadeStub) GetMetricsForPrometheus() string {
	if fs.GetMetricsForPrometheusCalled != nil {
		return fs.GetMetricsForPrometheusCalled()
	}

	return ""
}

// IsInterfaceNil -
func (fs *FacadeStub) IsInterfaceNil() bool {
	return fs == nil
}
