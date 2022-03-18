package facade

import "github.com/ElrondNetwork/notifier-go/data"

type EventsHandler interface {
	HandlePushEvents(events data.BlockEvents)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	IsInterfaceNil() bool
}
