package groups

import (
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// EventsFacadeHandler defines the behavior of a facade handler needed for events group
type EventsFacadeHandler interface {
	HandlePushEventsV2(events data.ArgsSaveBlockData) error
	HandlePushEventsV1(events data.SaveBlockData) error
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
	IsInterfaceNil() bool
}

// HubFacadeHandler defines the behavior of a facade handler needed for hub group
type HubFacadeHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	IsInterfaceNil() bool
}

// EmptyBlockCreatorContainer defines the behavior of a empty block creator container
type EmptyBlockCreatorContainer interface {
	Add(headerType core.HeaderType, creator block.EmptyBlockCreator) error
	Get(headerType core.HeaderType) (block.EmptyBlockCreator, error)
	IsInterfaceNil() bool
}

// EventsDataHandler defines the behaviour of an events data handler component
type EventsDataHandler interface {
	UnmarshallBlockDataV1(marshalledData []byte) (*data.SaveBlockData, error)
	UnmarshallBlockData(marshalledData []byte) (*data.ArgsSaveBlockData, error)
	IsInterfaceNil() bool
}
