package process

import (
	"context"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
)

var log = logger.GetOrCreate("process")

const (
	setRetryDuration       = time.Millisecond * 500
	reconnectRetryDuration = time.Second * 2
	minRetries             = 1
	revertKeyPrefix        = "revert_"
	finalizedKeyPrefix     = "finalized_"
)

// ArgsEventsHandler defines the arguments needed for an events handler
type ArgsEventsHandler struct {
	Config    config.ConnectorApiConfig
	Locker    LockService
	Publisher Publisher
}

type eventsHandler struct {
	config    config.ConnectorApiConfig
	locker    LockService
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
		config:    args.Config,
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
	if events.Hash == "" {
		log.Warn("received empty events block hash",
			"will process", false,
		)
		return
	}
	shouldProcessEvents := true
	if eh.config.CheckDuplicates {
		shouldProcessEvents = eh.tryCheckProcessedWithRetry(events.Hash)
	}

	if events.Events != nil && shouldProcessEvents {
		log.Info("received events for block",
			"block hash", events.Hash,
			"will process", shouldProcessEvents,
		)
		eh.publisher.Broadcast(events)
		return
	}

	log.Info("received duplicated events for block",
		"block hash", events.Hash,
		"processed", false,
	)
}

// HandleRevertEvents will handle revents events received from observer
func (eh *eventsHandler) HandleRevertEvents(revertBlock data.RevertBlock) {
	if revertBlock.Hash == "" {
		log.Warn("received empty revert block hash",
			"will process", false,
		)
		return
	}

	shouldProcessRevert := true
	if eh.config.CheckDuplicates {
		revertKey := revertKeyPrefix + revertBlock.Hash
		shouldProcessRevert = eh.tryCheckProcessedWithRetry(revertKey)
	}

	if shouldProcessRevert {
		log.Info("received revert event for block",
			"block hash", revertBlock.Hash,
			"will process", shouldProcessRevert,
		)
		eh.publisher.BroadcastRevert(revertBlock)
		return
	}

	log.Info("received duplicated revert event for block",
		"block hash", revertBlock.Hash,
		"processed", false,
	)
}

// HandleFinalizedEvents will handle finalized events received from observer
func (eh *eventsHandler) HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) {
	if finalizedBlock.Hash == "" {
		log.Warn("received empty finalized block hash",
			"will process", false,
		)
		return
	}
	shouldProcessFinalized := true
	if eh.config.CheckDuplicates {
		finalizedKey := finalizedKeyPrefix + finalizedBlock.Hash
		shouldProcessFinalized = eh.tryCheckProcessedWithRetry(finalizedKey)
	}

	if shouldProcessFinalized {
		log.Info("received finalized events for block",
			"block hash", finalizedBlock.Hash,
			"will process", shouldProcessFinalized,
		)
		eh.publisher.BroadcastFinalized(finalizedBlock)
		return
	}

	log.Info("received duplicated finalized event for block",
		"block hash", finalizedBlock.Hash,
		"processed", false,
	)
}

func (eh *eventsHandler) tryCheckProcessedWithRetry(blockHash string) bool {
	var err error
	var setSuccessful bool

	for {
		setSuccessful, err = eh.locker.IsEventProcessed(context.Background(), blockHash)

		if err == nil {
			break
		}

		log.Error("failed to check event in locker", "error", err.Error())
		if !eh.locker.HasConnection(context.Background()) {
			log.Error("failure connecting to locker service")

			time.Sleep(reconnectRetryDuration)
		} else {
			time.Sleep(setRetryDuration)
		}
	}

	if err != nil {
		log.Error("failed to check event in locker", "error", err.Error())
		return false
	}

	if !setSuccessful {
		log.Debug("did not succeed to set event in locker")
		return false
	}

	return true
}

// IsInterfaceNil returns true if there is no value under the interface
func (eh *eventsHandler) IsInterfaceNil() bool {
	return eh == nil
}
