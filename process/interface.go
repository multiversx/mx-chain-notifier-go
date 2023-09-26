package process

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// LockService defines the behaviour of a lock service component.
// It makes sure that a duplicated entry is not processed multiple times.
type LockService interface {
	IsEventProcessed(ctx context.Context, blockHash string) (bool, error)
	HasConnection(ctx context.Context) bool
	IsInterfaceNil() bool
}

// Publisher defines the behaviour of a publisher component which should be
// able to publish received events and broadcast them to channels
type Publisher interface {
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	BroadcastTxs(event data.BlockTxs)
	BroadcastBlockEventsWithOrder(event data.BlockEventsWithOrder)
	BroadcastScrs(event data.BlockScrs)
	IsInterfaceNil() bool
}

// EventsHandler defines the behaviour of an events handler component
type EventsHandler interface {
	HandlePushEvents(events data.BlockEvents) error
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	HandleBlockTxs(blockTxs data.BlockTxs)
	HandleBlockScrs(blockScrs data.BlockScrs)
	HandleBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder)
	IsInterfaceNil() bool
}

// EventsInterceptor defines the behaviour of an events interceptor component
type EventsInterceptor interface {
	ProcessBlockEvents(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error)
	IsInterfaceNil() bool
}

// WSClient defines what a websocket client should do
type WSClient interface {
	Close() error
}

// DataProcessor dines what a data indexer should do
type DataProcessor interface {
	SaveBlock(outportBlock *outport.OutportBlock) error
	RevertIndexedBlock(blockData *data.RevertBlock) error
	FinalizedBlock(finalizedBlock *outport.FinalizedBlock) error
	IsInterfaceNil() bool
}

// EmptyBlockCreatorContainer defines the behavior of a empty block creator container
type EmptyBlockCreatorContainer interface {
	Add(headerType core.HeaderType, creator block.EmptyBlockCreator) error
	Get(headerType core.HeaderType) (block.EmptyBlockCreator, error)
	IsInterfaceNil() bool
}

// EventsFacadeHandler defines the behavior of a facade handler needed for events group
type EventsFacadeHandler interface {
	HandlePushEventsV2(events data.ArgsSaveBlockData) error
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	IsInterfaceNil() bool
}
