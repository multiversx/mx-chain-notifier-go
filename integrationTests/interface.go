package integrationTests

import (
	"net/http"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// FacadeHandler defines facade behaviour
type FacadeHandler interface {
	HandlePushEventsV2(events data.ArgsSaveBlockData) error
	HandlePushEventsV1(events data.SaveBlockData) error
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	IsInterfaceNil() bool
}

// PublisherHandler defines publisher behaviour
type PublisherHandler interface {
	Run()
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	Close() error
	IsInterfaceNil() bool
}

// ObserverConnector defines the observer connector behaviour
type ObserverConnector interface {
	PushEventsRequest(events *outport.OutportBlock) error
	RevertEventsRequest(events *outport.BlockData) error
	FinalizedEventsRequest(events *outport.FinalizedBlock) error
	Close() error
}
