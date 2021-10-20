package pubsub

type LockService interface {
	IsBlockProcessed(blockHash string) (bool, error)
	HasConnection() bool
}
