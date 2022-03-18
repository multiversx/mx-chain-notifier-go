package groups

import (
	"context"

	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

// TODO: move this after further refactoring
// 		 add separate Publisher interface

// LockService defines the behaviour of a lock service component.
// It makes sure that a duplicated entry is not processed multiple times.
type LockService interface {
	// TODO: update this function name after proxy refactoring
	IsBlockProcessed(ctx context.Context, blockHash string) (bool, error)
	HasConnection(ctx context.Context) bool
	IsInterfaceNil() bool
}

type EventsFacadeHandler interface {
	HandlePushEvents(events data.BlockEvents)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
}

type HubFacadeHandler interface {
	GetHub() dispatcher.Hub
	GetDispatchType() string
	//GetGraphQLHandler()
}
