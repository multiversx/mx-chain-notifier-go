package mocks

import "context"

// LockerStub implements LockService interface
type LockerStub struct {
	IsEventProcessedCalled func(ctx context.Context, blockHash string) (bool, error)
	TryLockCalled          func(ctx context.Context, key string) (bool, error)
	UnlockCalled           func(ctx context.Context, key string) error
	HasConnectionCalled    func(ctx context.Context) bool
}

// IsEventProcessed -
func (ls *LockerStub) IsEventProcessed(ctx context.Context, blockHash string) (bool, error) {
	if ls.IsEventProcessedCalled != nil {
		return ls.IsEventProcessedCalled(ctx, blockHash)
	}

	return false, nil
}

// TryLock -
func (ls *LockerStub) TryLock(ctx context.Context, key string) (bool, error) {
	if ls.TryLockCalled != nil {
		return ls.TryLockCalled(ctx, key)
	}

	return true, nil
}

// Unlock -
func (ls *LockerStub) Unlock(ctx context.Context, key string) error {
	if ls.UnlockCalled != nil {
		return ls.UnlockCalled(ctx, key)
	}

	return nil
}

// HasConnection -
func (ls *LockerStub) HasConnection(ctx context.Context) bool {
	if ls.HasConnectionCalled != nil {
		return ls.HasConnectionCalled(ctx)
	}

	return false
}

// IsInterfaceNil -
func (ls *LockerStub) IsInterfaceNil() bool {
	return ls == nil
}
