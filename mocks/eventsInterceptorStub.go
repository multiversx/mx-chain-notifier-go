package mocks

import "github.com/multiversx/mx-chain-notifier-go/data"

// EventsInterceptorStub -
type EventsInterceptorStub struct {
	ProcessBlockEventsCalled func(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error)
}

// ProcessBlockEvents -
func (stub *EventsInterceptorStub) ProcessBlockEvents(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error) {
	if stub.ProcessBlockEventsCalled != nil {
		return stub.ProcessBlockEventsCalled(eventsData)
	}

	return nil, nil
}

// IsInterfaceNil returns true if there is not value under the interface
func (stub *EventsInterceptorStub) IsInterfaceNil() bool {
	return stub == nil
}
