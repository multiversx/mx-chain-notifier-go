package mocks

// EventsDataProcessorStub -
type EventsDataProcessorStub struct {
	SaveBlockCalled          func(marshalledData []byte) error
	RevertIndexedBlockCalled func(marshalledData []byte) error
	FinalizedBlockCalled     func(marshalledData []byte) error
}

// SaveBlock -
func (stub *EventsDataProcessorStub) SaveBlock(marshalledData []byte) error {
	if stub.SaveBlockCalled != nil {
		return stub.SaveBlockCalled(marshalledData)
	}

	return nil
}

// RevertIndexedBlock -
func (stub *EventsDataProcessorStub) RevertIndexedBlock(marshalledData []byte) error {
	if stub.RevertIndexedBlockCalled != nil {
		return stub.RevertIndexedBlockCalled(marshalledData)
	}

	return nil
}

// FinalizedBlock -
func (stub *EventsDataProcessorStub) FinalizedBlock(marshalledData []byte) error {
	if stub.FinalizedBlockCalled != nil {
		return stub.FinalizedBlockCalled(marshalledData)
	}

	return nil
}

// Close -
func (stub *EventsDataProcessorStub) Close() error {
	return nil
}

// IsInterfaceNil -
func (stub *EventsDataProcessorStub) IsInterfaceNil() bool {
	return stub == nil
}
