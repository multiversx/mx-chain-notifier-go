package mocks

import "github.com/multiversx/mx-chain-notifier-go/data"

// EventsDataHandlerStub -
type EventsDataHandlerStub struct {
	UnmarshallBlockDataOldCalled  func(marshalledData []byte) (*data.SaveBlockData, error)
	UnmarshallBlockDataCalled     func(marshalledData []byte) (*data.ArgsSaveBlockData, error)
	UnmarshallRevertDataCalled    func(marshalledData []byte) (*data.RevertBlock, error)
	UnmarshallFinalizedDataCalled func(marshalledData []byte) (*data.FinalizedBlock, error)
}

// UnmarshallBlockDataOld -
func (stub *EventsDataHandlerStub) UnmarshallBlockDataOld(marshalledData []byte) (*data.SaveBlockData, error) {
	if stub.UnmarshallBlockDataOldCalled != nil {
		return stub.UnmarshallBlockDataOldCalled(marshalledData)
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

// UnmarshallRevertData -
func (stub *EventsDataHandlerStub) UnmarshallRevertData(marshalledData []byte) (*data.RevertBlock, error) {
	if stub.UnmarshallRevertDataCalled != nil {
		return stub.UnmarshallRevertDataCalled(marshalledData)
	}

	return &data.RevertBlock{}, nil
}

// UnmarshallFinalizedData -
func (stub *EventsDataHandlerStub) UnmarshallFinalizedData(marshalledData []byte) (*data.FinalizedBlock, error) {
	if stub.UnmarshallFinalizedDataCalled != nil {
		return stub.UnmarshallFinalizedDataCalled(marshalledData)
	}

	return &data.FinalizedBlock{}, nil
}

// IsInterfaceNil -
func (stub *EventsDataHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
