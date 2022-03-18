package notifier

import (
	"context"

	"github.com/ElrondNetwork/notifier-go/data"
)

// LockService defines the behaviour of a lock service component.
// It makes sure that a duplicated entry is not processed multiple times.
type LockService interface {
	// TODO: update this function name after proxy refactoring
	IsBlockProcessed(ctx context.Context, blockHash string) (bool, error)
	HasConnection(ctx context.Context) bool
	IsInterfaceNil() bool
}

// PublisherService defines the behaviour of a publisher component which should be
// able to publish received events and broadcast them to channels
type PublisherService interface {
	BroadcastChan() chan<- data.BlockEvents
	BroadcastRevertChan() chan<- data.RevertBlock
	BroadcastFinalizedChan() chan<- data.FinalizedBlock
	IsInterfaceNil() bool
}

// TODO: evaluate adding Subscriber interface and split Hub interface if suitable
