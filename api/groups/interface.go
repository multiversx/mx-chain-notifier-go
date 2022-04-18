package groups

import (
	"net/http"

	"github.com/ElrondNetwork/notifier-go/data"
)

// EventsFacadeHandler defines the behavior of a facade handler needed for events group
type EventsFacadeHandler interface {
	HandlePushEvents(events data.BlockEvents)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	HandleTxsEvents(blockTxs data.BlockTxs)
	GetConnectorUserAndPass() (string, string)
	IsInterfaceNil() bool
}

// HubFacadeHandler defines the behavior of a facade handler needed for hub group
type HubFacadeHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	IsInterfaceNil() bool
}
