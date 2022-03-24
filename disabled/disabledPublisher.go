package disabled

import (
	"github.com/ElrondNetwork/notifier-go/data"
)

// Publisher defines a disabled publisher component
type Publisher struct{}

// Run does nothing
func (dp *Publisher) Run() {
}

// Broadcast does nothing
func (dp *Publisher) Broadcast(events data.BlockEvents) {
}

// BroadcastRevert does nothing
func (dp *Publisher) BroadcastRevert(event data.RevertBlock) {
}

// BroadcastFinalized does nothing
func (dp *Publisher) BroadcastFinalized(event data.FinalizedBlock) {
}

// Close returns nil
func (dp *Publisher) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (dp *Publisher) IsInterfaceNil() bool {
	return false
}
