package notifier

import (
	"context"
	"time"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/redis"
)

// TODO: comments update

const (
	setnxRetryMs       = 500
	revertKeyPrefix    = "revert_"
	finalizedKeyPrefix = "finalized_"
)

type ArgsEventsHandler struct {
	Config    config.ConnectorApiConfig
	Locker    redis.LockService
	Publisher PublisherService
}

type eventsHandler struct {
	config    config.ConnectorApiConfig
	locker    redis.LockService
	publisher PublisherService
}

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
	// TODO: check configs
	if check.IfNil(args.Locker) {
		return ErrNilLockService
	}
	if check.IfNil(args.Publisher) {
		return ErrNilPublisherService
	}

	return nil
}

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
		setSuccessful, err = eh.locker.IsBlockProcessed(context.Background(), blockHash)

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
