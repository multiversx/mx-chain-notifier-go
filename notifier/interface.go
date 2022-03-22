package notifier

import (
	"github.com/ElrondNetwork/notifier-go/data"
)

// Publisher defines the behaviour of a publisher component which should be
// able to publish received events and broadcast them to channels
type Publisher interface {
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	IsInterfaceNil() bool
}

// TODO: evaluate adding Subscriber interface and split Hub interface if suitable
