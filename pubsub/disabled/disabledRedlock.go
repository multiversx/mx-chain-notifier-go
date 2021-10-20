package disabled

type disabledRedlockWrapper struct {
}

func NewDisabledRedlockWrapper() *disabledRedlockWrapper {
	return &disabledRedlockWrapper{}
}

func (drw *disabledRedlockWrapper) IsBlockProcessed(blockHash string) (bool, error) {
	return true, nil
}

func (drw *disabledRedlockWrapper) HasConnection() bool {
	return true
}
