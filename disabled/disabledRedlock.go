package disabled

type disabledRedlockWrapper struct {
}

// NewDisabledRedlockWrapper created a new disabled Redlock wrapper
func NewDisabledRedlockWrapper() *disabledRedlockWrapper {
	return &disabledRedlockWrapper{}
}

// IsBlockProcessed returns true and nil
func (drw *disabledRedlockWrapper) IsBlockProcessed(blockHash string) (bool, error) {
	return true, nil
}

// HasConnection returns true
func (drw *disabledRedlockWrapper) HasConnection() bool {
	return true
}

// IsInterfaceNil returns true if there is no value under the interface
func (drw *disabledRedlockWrapper) IsInterfaceNil() bool {
	return false
}
