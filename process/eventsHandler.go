package process

import (
	"context"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

var log = logger.GetOrCreate("process")

const (
	setRetryDuration       = time.Millisecond * 500
	reconnectRetryDuration = time.Second * 2
	minRetries             = 1
	revertKeyPrefix        = "revert_"
	finalizedKeyPrefix     = "finalized_"
	txsKeyPrefix           = "txs_"
	txsWithOrderKeyPrefix  = "txsWithOrder_"
	scrsKeyPrefix          = "scrs_"
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
func (eh *eventsHandler) HandlePushEvents(events data.BlockEvents) error {
	if events.Hash == "" {
		log.Debug("received empty hash", "event", common.PushLogsAndEvents,
			"will process", false,
		)
		return common.ErrReceivedEmptyEvents
	}

	shouldProcessEvents := true
	if eh.config.CheckDuplicates {
		shouldProcessEvents = eh.tryCheckProcessedWithRetry(common.PushLogsAndEvents, events.Hash)
	}

	if !shouldProcessEvents {
		log.Info("received duplicated events", "event", common.PushLogsAndEvents,
			"block hash", events.Hash,
			"will process", false,
		)
		return nil
	}

	if len(events.Events) == 0 {
		log.Warn("received empty events", "event", common.PushLogsAndEvents,
			"block hash", events.Hash,
			"will process", shouldProcessEvents,
		)
		events.Events = make([]data.Event, 0)
	} else {
		log.Info("received", "event", common.PushLogsAndEvents,
			"block hash", events.Hash,
			"will process", shouldProcessEvents,
		)
	}

	eh.publisher.Broadcast(events)
	return nil
}

// HandleRevertEvents will handle revents events received from observer
func (eh *eventsHandler) HandleRevertEvents(revertBlock data.RevertBlock) {
	if revertBlock.Hash == "" {
		log.Warn("received empty hash", "event", common.RevertBlockEvents,
			"will process", false,
		)
		return
	}

	shouldProcessRevert := true
	if eh.config.CheckDuplicates {
		shouldProcessRevert = eh.tryCheckProcessedWithRetry(common.RevertBlockEvents, revertBlock.Hash)
	}

	if !shouldProcessRevert {
		log.Info("received duplicated events", "event", common.RevertBlockEvents,
			"block hash", revertBlock.Hash,
			"will process", false,
		)
		return
	}

	log.Info("received", "event", common.RevertBlockEvents,
		"block hash", revertBlock.Hash,
		"will process", shouldProcessRevert,
	)

	eh.publisher.BroadcastRevert(revertBlock)
}

// HandleFinalizedEvents will handle finalized events received from observer
func (eh *eventsHandler) HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) {
	if finalizedBlock.Hash == "" {
		log.Warn("received empty hash", "event", common.FinalizedBlockEvents,
			"will process", false,
		)
		return
	}
	shouldProcessFinalized := true
	if eh.config.CheckDuplicates {
		shouldProcessFinalized = eh.tryCheckProcessedWithRetry(common.FinalizedBlockEvents, finalizedBlock.Hash)
	}

	if !shouldProcessFinalized {
		log.Info("received duplicated events", "event", common.FinalizedBlockEvents,
			"block hash", finalizedBlock.Hash,
			"will process", false,
		)
		return
	}

	log.Info("received", "event", common.FinalizedBlockEvents,
		"block hash", finalizedBlock.Hash,
		"will process", shouldProcessFinalized,
	)

	eh.publisher.BroadcastFinalized(finalizedBlock)
}

// HandleBlockTxs will handle txs events received from observer
func (eh *eventsHandler) HandleBlockTxs(blockTxs data.BlockTxs) {
	if blockTxs.Hash == "" {
		log.Warn("received empty hash", "event", common.BlockTxs,
			"will process", false,
		)
		return
	}
	shouldProcessTxs := true
	if eh.config.CheckDuplicates {
		shouldProcessTxs = eh.tryCheckProcessedWithRetry(common.BlockTxs, blockTxs.Hash)
	}

	if !shouldProcessTxs {
		log.Info("received duplicated events", "event", common.BlockTxs,
			"block hash", blockTxs.Hash,
			"will process", false,
		)
		return
	}

	if len(blockTxs.Txs) == 0 {
		log.Warn("received empty events", "event", common.BlockTxs,
			"block hash", blockTxs.Hash,
			"will process", shouldProcessTxs,
		)
	} else {
		log.Info("received", "event", common.BlockTxs,
			"block hash", blockTxs.Hash,
			"will process", shouldProcessTxs,
		)
	}

	eh.publisher.BroadcastTxs(blockTxs)
}

// HandleBlockScrs will handle scrs events received from observer
func (eh *eventsHandler) HandleBlockScrs(blockScrs data.BlockScrs) {
	if blockScrs.Hash == "" {
		log.Warn("received empty hash", "event", common.BlockScrs,
			"will process", false,
		)
		return
	}
	shouldProcessScrs := true
	if eh.config.CheckDuplicates {
		shouldProcessScrs = eh.tryCheckProcessedWithRetry(common.BlockScrs, blockScrs.Hash)
	}

	if !shouldProcessScrs {
		log.Info("received duplicated events", "event", common.BlockScrs,
			"block hash", blockScrs.Hash,
			"will process", false,
		)
		return
	}

	if len(blockScrs.Scrs) == 0 {
		log.Warn("received empty events", "event", common.BlockScrs,
			"block hash", blockScrs.Hash,
			"will process", shouldProcessScrs,
		)
	} else {
		log.Info("received", "event", common.BlockScrs,
			"block hash", blockScrs.Hash,
			"will process", shouldProcessScrs,
		)
	}

	eh.publisher.BroadcastScrs(blockScrs)
}

// HandleBlockEventsWithOrder will handle full block events received from observer
func (eh *eventsHandler) HandleBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
	if blockTxs.Hash == "" {
		log.Warn("received empty hash", "event", common.BlockEvents,
			"will process", false,
		)
		return
	}
	shouldProcessTxs := true
	if eh.config.CheckDuplicates {
		shouldProcessTxs = eh.tryCheckProcessedWithRetry(common.BlockEvents, blockTxs.Hash)
	}

	if !shouldProcessTxs {
		log.Info("received duplicated events", "event", common.BlockEvents,
			"block hash", blockTxs.Hash,
			"will process", false,
		)
		return
	}

	log.Info("received", "event", common.BlockEvents,
		"block hash", blockTxs.Hash,
		"will process", shouldProcessTxs,
	)

	eh.publisher.BroadcastBlockEventsWithOrder(blockTxs)
}

func (eh *eventsHandler) tryCheckProcessedWithRetry(id, blockHash string) bool {
	var err error
	var setSuccessful bool

	prefix := getPrefixLockerKey(id)
	key := prefix + blockHash

	for {
		setSuccessful, err = eh.locker.IsEventProcessed(context.Background(), key)

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

	log.Debug("locker", "event", id, "block hash", blockHash, "succeeded", setSuccessful)

	return setSuccessful
}

func getPrefixLockerKey(id string) string {
	// keep this matching for backwards compatibility
	switch id {
	case common.PushLogsAndEvents:
		return ""
	case common.RevertBlockEvents:
		return revertKeyPrefix
	case common.FinalizedBlockEvents:
		return finalizedKeyPrefix
	case common.BlockTxs:
		return txsKeyPrefix
	case common.BlockScrs:
		return scrsKeyPrefix
	case common.BlockEvents:
		return txsWithOrderKeyPrefix
	}

	return ""
}

// IsInterfaceNil returns true if there is no value under the interface
func (eh *eventsHandler) IsInterfaceNil() bool {
	return eh == nil
}
