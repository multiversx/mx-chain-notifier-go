package groups

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// EventsFacadeHandler defines the behavior of a facade handler needed for events group
type EventsFacadeHandler interface {
	HandlePushEvents(events data.BlockEvents)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
	IsInterfaceNil() bool
}

// HubFacadeHandler defines the behavior of a facade handler needed for hub group
type HubFacadeHandler interface {
	GetHub() dispatcher.Hub
	GetDispatchType() string
	IsInterfaceNil() bool
}
