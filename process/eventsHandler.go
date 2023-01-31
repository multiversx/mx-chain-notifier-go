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
		log.Debug("received empty events block hash",
			"will process", false,
		)
		return common.ErrReceivedEmptyEvents
	}
	shouldProcessEvents := true
	if eh.config.CheckDuplicates {
		shouldProcessEvents = eh.tryCheckProcessedWithRetry(events.Hash)
	}

	if !shouldProcessEvents {
		log.Info("received duplicated events for block",
			"block hash", events.Hash,
			"processed", false,
		)
		return nil
	}

	if len(events.Events) == 0 {
		log.Info("received empty events for block",
			"block hash", events.Hash,
			"will process", shouldProcessEvents,
		)
		events.Events = make([]data.Event, 0)
	} else {
		log.Info("received events for block",
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

	if !shouldProcessRevert {
		log.Info("received duplicated revert event for block",
			"block hash", revertBlock.Hash,
			"processed", false,
		)
		return
	}

	log.Info("received revert event for block",
		"block hash", revertBlock.Hash,
		"will process", shouldProcessRevert,
	)

	eh.publisher.BroadcastRevert(revertBlock)
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

	if !shouldProcessFinalized {
		log.Info("received duplicated finalized event for block",
			"block hash", finalizedBlock.Hash,
			"processed", false,
		)
		return
	}

	log.Info("received finalized events for block",
		"block hash", finalizedBlock.Hash,
		"will process", shouldProcessFinalized,
	)

	eh.publisher.BroadcastFinalized(finalizedBlock)
}

// HandleBlockTxs will handle txs events received from observer
func (eh *eventsHandler) HandleBlockTxs(blockTxs data.BlockTxs) {
	if blockTxs.Hash == "" {
		log.Warn("received empty txs block hash",
			"will process", false,
		)
		return
	}
	shouldProcessTxs := true
	if eh.config.CheckDuplicates {
		txsKey := txsKeyPrefix + blockTxs.Hash
		shouldProcessTxs = eh.tryCheckProcessedWithRetry(txsKey)
	}

	if !shouldProcessTxs {
		log.Info("received duplicated txs event for block",
			"block hash", blockTxs.Hash,
			"processed", false,
		)
		return
	}

	if len(blockTxs.Txs) == 0 {
		log.Info("received empty txs event for block",
			"block hash", blockTxs.Hash,
			"will process", shouldProcessTxs,
		)
	} else {
		log.Info("received txs events for block",
			"block hash", blockTxs.Hash,
			"will process", shouldProcessTxs,
		)
	}

	eh.publisher.BroadcastTxs(blockTxs)
}

// HandleBlockScrs will handle scrs events received from observer
func (eh *eventsHandler) HandleBlockScrs(blockScrs data.BlockScrs) {
	if blockScrs.Hash == "" {
		log.Warn("received empty scrs block hash",
			"will process", false,
		)
		return
	}
	shouldProcessScrs := true
	if eh.config.CheckDuplicates {
		scrsKey := scrsKeyPrefix + blockScrs.Hash
		shouldProcessScrs = eh.tryCheckProcessedWithRetry(scrsKey)
	}

	if !shouldProcessScrs {
		log.Info("received duplicated scrs event for block",
			"block hash", blockScrs.Hash,
			"processed", false,
		)
		return
	}

	if len(blockScrs.Scrs) == 0 {
		log.Info("received empty scrs event for block",
			"block hash", blockScrs.Hash,
			"will process", shouldProcessScrs,
		)
	} else {
		log.Info("received scrs events for block",
			"block hash", blockScrs.Hash,
			"will process", shouldProcessScrs,
		)
	}

	eh.publisher.BroadcastScrs(blockScrs)
}

// HandleBlockEventsWithOrder will handle full block events received from observer
func (eh *eventsHandler) HandleBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
	if blockTxs.Hash == "" {
		log.Warn("received full block events with empty block hash",
			"will process", false,
		)
		return
	}
	shouldProcessTxs := true
	if eh.config.CheckDuplicates {
		txsKey := txsWithOrderKeyPrefix + blockTxs.Hash
		shouldProcessTxs = eh.tryCheckProcessedWithRetry(txsKey)
	}

	if !shouldProcessTxs {
		log.Info("received duplicated full block events for block",
			"block hash", blockTxs.Hash,
			"processed", false,
		)
		return
	}

	if len(blockTxs.Txs) == 0 {
		log.Info("received empty txs with order event for block",
			"block hash", blockTxs.Hash,
			"will process", shouldProcessTxs,
		)
	} else {
		log.Info("received txs events for block",
			"block hash", blockTxs.Hash,
			"will process", shouldProcessTxs,
		)
	}

	eh.publisher.BroadcastBlockEventsWithOrder(blockTxs)
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
