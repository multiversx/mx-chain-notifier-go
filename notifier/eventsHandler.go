package notifier

import (
	"context"
	"time"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/redis"
)

const (
	setnxRetryMs       = 500
	revertKeyPrefix    = "revert_"
	finalizedKeyPrefix = "finalized_"
)

// ArgsEventsHandler defines the arguments needed for an events handler
type ArgsEventsHandler struct {
	Config    config.ConnectorApiConfig
	Locker    redis.LockService
	Publisher Publisher
}

type eventsHandler struct {
	config    config.ConnectorApiConfig
	locker    redis.LockService
	publisher Publisher
}

// NewEventsHandler creates a new events handler component
func NewEventsHandler(args ArgsEventsHandler) (*eventsHandler, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &eventsHandler{
		locker:    args.Locker,
		publisher: args.Publisher,
	}, nil
}

func checkArgs(args ArgsEventsHandler) error {
	if check.IfNil(args.Locker) {
		return ErrNilLockService
	}
	if check.IfNil(args.Publisher) {
		return ErrNilPublisherService
	}

	return nil
}

// HandlePushEvents will handle push events received from observer
func (eh *eventsHandler) HandlePushEvents(events data.BlockEvents) {
	shouldProcessEvents := true
	if eh.config.CheckDuplicates {
		shouldProcessEvents = eh.tryCheckProcessedOrRetry(events.Hash)
	}

	if events.Events != nil && shouldProcessEvents {
		log.Info("received events for block",
			"block hash", events.Hash,
			"shouldProcess", shouldProcessEvents,
		)
		eh.publisher.BroadcastChan() <- events
	} else {
		log.Info("received duplicated events for block",
			"block hash", events.Hash,
			"ignoring", true,
		)
	}
}

// HandleRevertEvents will handle revents events received from observer
func (eh *eventsHandler) HandleRevertEvents(revertBlock data.RevertBlock) {
	shouldProcessRevert := true
	if eh.config.CheckDuplicates {
		revertKey := revertKeyPrefix + revertBlock.Hash
		shouldProcessRevert = eh.tryCheckProcessedOrRetry(revertKey)
	}

	if shouldProcessRevert {
		log.Info("received revert event for block",
			"block hash", revertBlock.Hash,
			"shouldProcess", shouldProcessRevert,
		)
		eh.publisher.BroadcastRevertChan() <- revertBlock
	} else {
		log.Info("received duplicated revert event for block",
			"block hash", revertBlock.Hash,
			"ignoring", true,
		)
	}
}

// HandleFinalizedEvents will handle finalized events received from observer
func (eh *eventsHandler) HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) {
	shouldProcessFinalized := true
	if eh.config.CheckDuplicates {
		finalizedKey := finalizedKeyPrefix + finalizedBlock.Hash
		shouldProcessFinalized = eh.tryCheckProcessedOrRetry(finalizedKey)
	}

	if shouldProcessFinalized {
		log.Info("received finalized events for block",
			"block hash", finalizedBlock.Hash,
			"shouldProcess", shouldProcessFinalized,
		)
		eh.publisher.BroadcastFinalizedChan() <- finalizedBlock
	} else {
		log.Info("received duplicated finalized event for block",
			"block hash", finalizedBlock.Hash,
			"ignoring", true,
		)
	}
}

func (eh *eventsHandler) tryCheckProcessedOrRetry(blockHash string) bool {
	var err error
	var setSuccessful bool

	for {
		setSuccessful, err = eh.locker.IsEventProcessed(context.Background(), blockHash)

		if err != nil {
			if !eh.locker.HasConnection(context.Background()) {
				log.Error("failure connecting to redis")
				time.Sleep(time.Millisecond * setnxRetryMs)
				continue
			}
		}

		break
	}

	return err == nil && setSuccessful
}

// IsInterfaceNil returns true if there is no value under the interface
func (eh *eventsHandler) IsInterfaceNil() bool {
	return eh == nil
}
