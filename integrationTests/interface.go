package integrationTests

import (
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// PublisherHandler defines publisher behaviour
type PublisherHandler interface {
	Run() error
	Broadcast(events data.BlockEvents) error
	BroadcastRevert(event data.RevertBlock) error
	BroadcastFinalized(event data.FinalizedBlock) error
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
