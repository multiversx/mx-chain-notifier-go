package integrationTests

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// PublisherHandler defines publisher behaviour
type PublisherHandler interface {
	Run()
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	Close() error
	IsInterfaceNil() bool
}
