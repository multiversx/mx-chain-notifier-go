package mocks

import "github.com/ElrondNetwork/notifier-go/data"

// EventsInterceptorStub -
type EventsInterceptorStub struct {
	ProcessBlockEventsCalled func(eventsData *data.ArgsSaveBlockData) *data.SaveBlockData
}

// ProcessBlockEvents -
func (stub *EventsInterceptorStub) ProcessBlockEvents(eventsData *data.ArgsSaveBlockData) *data.SaveBlockData {
	if stub.ProcessBlockEventsCalled != nil {
		return stub.ProcessBlockEventsCalled(eventsData)
	}

	return nil
}

// IsInterfaceNil returns true if there is not value under the interface
func (stub *EventsInterceptorStub) IsInterfaceNil() bool {
	return stub == nil
}
