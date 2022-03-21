package facade

import "github.com/ElrondNetwork/notifier-go/data"

// EventsHandler defines the behavior of an events handler component.
// This will handle push events from observer node.
type EventsHandler interface {
	HandlePushEvents(events data.BlockEvents)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	IsInterfaceNil() bool
}
