package mocks

// PayloadHandlerStub -
type PayloadHandlerStub struct {
	ProcessPayloadCalled func(payload []byte, topic string, version uint32) error
	CloseCalled          func() error
}

// ProcessPayload -
func (ph *PayloadHandlerStub) ProcessPayload(payload []byte, topic string, version uint32) error {
	if ph.ProcessPayloadCalled != nil {
		return ph.ProcessPayloadCalled(payload, topic, version)
	}
	return nil
}

// Close -
func (ph *PayloadHandlerStub) Close() error {
	if ph.CloseCalled != nil {
		return ph.CloseCalled()
	}
	return nil
}

// IsInterfaceNil -
func (ph *PayloadHandlerStub) IsInterfaceNil() bool {
	return ph == nil
}
