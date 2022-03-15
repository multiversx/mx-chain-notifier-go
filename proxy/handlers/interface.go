package handlers

// TODO: move this after further refactoring
// 		 add separate Publisher interface

// LockService defines the behaviour of a lock service component.
// It makes sure that a duplicated entry is not processed multiple times.
type LockService interface {
	// TODO: update this function name after proxy refactoring
	IsBlockProcessed(blockHash string) (bool, error)
	HasConnection() bool
	IsInterfaceNil() bool
}
