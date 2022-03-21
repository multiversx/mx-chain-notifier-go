package disabled

import "github.com/ElrondNetwork/notifier-go/data"

// Publisher defines a disabled publisher component
type Publisher struct{}

// Run does nothing
func (dp *Publisher) Run() {
}

// BroadcastChan returns a nil channel
func (dp *Publisher) BroadcastChan() chan<- data.BlockEvents {
	return nil
}

// BroadcastRevertChan returns a nil channel
func (dp *Publisher) BroadcastRevertChan() chan<- data.RevertBlock {
	return nil
}

// BroadcastFinalizedChan returns a nil channel
func (dp *Publisher) BroadcastFinalizedChan() chan<- data.FinalizedBlock {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (dp *Publisher) IsInterfaceNil() bool {
	return false
}
