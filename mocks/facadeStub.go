package mocks

import (
	"net/http"

	"github.com/multiversx/mx-chain-notifier-go/data"
)

// FacadeStub implements FacadeHandler interface
type FacadeStub struct {
	HandlePushEventsV2Called      func(events data.ArgsSaveBlockData) error
	HandlePushEventsV1Called      func(eventsData data.SaveBlockData) error
	HandleRevertEventsCalled      func(events data.RevertBlock)
	HandleFinalizedEventsCalled   func(events data.FinalizedBlock)
	ServeCalled                   func(w http.ResponseWriter, r *http.Request)
	GetConnectorUserAndPassCalled func() (string, string)
}

// HandlePushEventsV2 -
func (fs *FacadeStub) HandlePushEventsV2(events data.ArgsSaveBlockData) error {
	if fs.HandlePushEventsV2Called != nil {
		return fs.HandlePushEventsV2Called(events)
	}

	return nil
}

// HandlePushEventsV1 -
func (fs *FacadeStub) HandlePushEventsV1(events data.SaveBlockData) error {
	if fs.HandlePushEventsV1Called != nil {
		return fs.HandlePushEventsV1Called(events)
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

// IsInterfaceNil -
func (fs *FacadeStub) IsInterfaceNil() bool {
	return fs == nil
}
