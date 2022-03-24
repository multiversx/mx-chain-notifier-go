package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// FacadeStub implements FacadeHandler interface
type FacadeStub struct {
	HandlePushEventsCalled        func(events data.BlockEvents)
	HandleRevertEventsCalled      func(events data.RevertBlock)
	HandleFinalizedEventsCalled   func(events data.FinalizedBlock)
	GetDispatchTypeCalled         func() string
	GetHubCalled                  func() dispatcher.Hub
	GetConnectorUserAndPassCalled func() (string, string)
}

// HandlePushEvents -
func (fs *FacadeStub) HandlePushEvents(events data.BlockEvents) {
	if fs.HandlePushEventsCalled != nil {
		fs.HandlePushEventsCalled(events)
	}
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

// GetDispatchType -
func (fs *FacadeStub) GetDispatchType() string {
	if fs.GetDispatchTypeCalled != nil {
		return fs.GetDispatchTypeCalled()
	}

	return ""
}

// GetHub -
func (fs *FacadeStub) GetHub() dispatcher.Hub {
	if fs.GetHubCalled != nil {
		return fs.GetHubCalled()
	}

	return nil
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
