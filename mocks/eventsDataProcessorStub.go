package mocks

import "github.com/multiversx/mx-chain-core-go/data/outport"

// EventsDataProcessorStub -
type EventsDataProcessorStub struct {
	SaveBlockCalled          func(outportBlock *outport.OutportBlock) error
	RevertIndexedBlockCalled func(blockData *outport.BlockData) error
	FinalizedBlockCalled     func(finalizedBlock *outport.FinalizedBlock) error
}

// SaveBlock -
func (stub *EventsDataProcessorStub) SaveBlock(outportBlock *outport.OutportBlock) error {
	if stub.SaveBlockCalled != nil {
		return stub.SaveBlockCalled(outportBlock)
	}

	return nil
}

// RevertIndexedBlock -
func (stub *EventsDataProcessorStub) RevertIndexedBlock(blockData *outport.BlockData) error {
	if stub.RevertIndexedBlockCalled != nil {
		return stub.RevertIndexedBlockCalled(blockData)
	}

	return nil
}

// FinalizedBlock -
func (stub *EventsDataProcessorStub) FinalizedBlock(finalizedBlock *outport.FinalizedBlock) error {
	if stub.FinalizedBlockCalled != nil {
		return stub.FinalizedBlockCalled(finalizedBlock)
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