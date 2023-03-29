package mocks

import "github.com/multiversx/mx-chain-notifier-go/data"

// EventsDataHandlerStub -
type EventsDataHandlerStub struct {
	UnmarshallBlockDataV1Called func(marshalledData []byte) (*data.SaveBlockData, error)
	UnmarshallBlockDataCalled   func(marshalledData []byte) (*data.ArgsSaveBlockData, error)
}

// UnmarshallBlockDataV1 -
func (stub *EventsDataHandlerStub) UnmarshallBlockDataV1(marshalledData []byte) (*data.SaveBlockData, error) {
	if stub.UnmarshallBlockDataV1Called != nil {
		return stub.UnmarshallBlockDataV1Called(marshalledData)
	}

	return &data.SaveBlockData{}, nil
}

// UnmarshallBlockData -
func (stub *EventsDataHandlerStub) UnmarshallBlockData(marshalledData []byte) (*data.ArgsSaveBlockData, error) {
	if stub.UnmarshallBlockDataCalled != nil {
		return stub.UnmarshallBlockDataCalled(marshalledData)
	}

	return &data.ArgsSaveBlockData{}, nil
}

// IsInterfaceNil -
func (stub *EventsDataHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
