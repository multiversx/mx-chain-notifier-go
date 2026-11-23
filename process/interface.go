package process

import (
	"context"

	"github.com/multiversx/mx-chain-notifier-go/data"
)

// LockService defines the behaviour of a lock service component.
// It makes sure that a duplicated entry is not processed multiple times.
type LockService interface {
	IsEventProcessed(ctx context.Context, blockHash string) (bool, error)
	TryLock(ctx context.Context, key string) (bool, error)
	Unlock(ctx context.Context, key string) error
	HasConnection(ctx context.Context) bool
	IsInterfaceNil() bool
}

// Publisher defines the behaviour of a publisher component which should be
// able to publish received events and broadcast them to channels
type Publisher interface {
	Run() error
	Broadcast(events data.BlockEvents) error
	BroadcastRevert(event data.RevertBlock) error
	BroadcastFinalized(event data.FinalizedBlock) error
	BroadcastTxs(event data.BlockTxs) error
	BroadcastBlockEventsWithOrder(event data.BlockEventsWithOrder) error
	BroadcastScrs(event data.BlockScrs) error
	BroadcastStateAccesses(events data.BlockStateAccesses) error
	Close() error
	IsInterfaceNil() bool
}

// EventsHandler defines the behaviour of an events handler component
type EventsHandler interface {
	HandleSaveBlockEvents(allEvents data.ArgsSaveBlockData) error
	HandleRevertEvents(revertBlock data.RevertBlock) error
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) error
	IsInterfaceNil() bool
}

// EventsInterceptor defines the behaviour of an events interceptor component
type EventsInterceptor interface {
	ProcessBlockEvents(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error)
	ProcessBlockEventsV3(eventsData *data.ArgsSaveBlockData) ([]*data.InterceptorBlockData, error)
	IsInterfaceNil() bool
}

// WSClient defines what a websocket client should do
type WSClient interface {
	Close() error
}

// DataProcessor dines what a data indexer should do
type DataProcessor interface {
	SaveBlock(marshalledData []byte) error
	RevertIndexedBlock(marshalledData []byte) error
	FinalizedBlock(marshalledData []byte) error
	IsInterfaceNil() bool
}

// EventsFacadeHandler defines the behavior of a facade handler needed for events group
type EventsFacadeHandler interface {
	HandlePushEvents(events data.ArgsSaveBlockData) error
	HandleRevertEvents(revertBlock data.RevertBlock) error
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) error
	IsInterfaceNil() bool
}

// PublisherHandler defines the behavior of a publisher component
type PublisherHandler interface {
	Publish(events data.BlockEvents) error
	PublishRevert(revertBlock data.RevertBlock) error
	PublishFinalized(finalizedBlock data.FinalizedBlock) error
	PublishTxs(blockTxs data.BlockTxs) error
	PublishScrs(blockScrs data.BlockScrs) error
	PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) error
	PublishStateAccesses(stateAccesses data.BlockStateAccesses) error
	Close() error
	IsInterfaceNil() bool
}
