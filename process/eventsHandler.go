package process

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

var log = logger.GetOrCreate("process")

const (
	setRetryDuration       = time.Millisecond * 500
	reconnectRetryDuration = time.Second * 2
	minRetries             = 1
	revertKeyPrefix        = "revert_"
	finalizedKeyPrefix     = "finalized_"
	v3BatchLockPrefix      = "v3batchlock_"

	rabbitmqMetricPrefix = "RabbitMQ"
	redisMetricPrefix    = "Redis"
)

// ArgsEventsHandler defines the arguments needed for an events handler
type ArgsEventsHandler struct {
	Locker               LockService
	Publisher            Publisher
	StatusMetricsHandler common.StatusMetricsHandler
	EventsInterceptor    EventsInterceptor
	CheckDuplicates      bool
}

type eventsHandler struct {
	locker            LockService
	publisher         Publisher
	metricsHandler    common.StatusMetricsHandler
	eventsInterceptor EventsInterceptor
	checkDuplicates   bool
}

// NewEventsHandler creates a new events handler component
func NewEventsHandler(args ArgsEventsHandler) (*eventsHandler, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &eventsHandler{
		locker:            args.Locker,
		publisher:         args.Publisher,
		metricsHandler:    args.StatusMetricsHandler,
		eventsInterceptor: args.EventsInterceptor,
		checkDuplicates:   args.CheckDuplicates,
	}, nil
}

func checkArgs(args ArgsEventsHandler) error {
	if check.IfNil(args.Locker) {
		return ErrNilLockService
	}
	if check.IfNil(args.Publisher) {
		return ErrNilPublisherService
	}
	if check.IfNil(args.StatusMetricsHandler) {
		return common.ErrNilStatusMetricsHandler
	}
	if check.IfNil(args.EventsInterceptor) {
		return ErrNilEventsInterceptor
	}

	return nil
}

// HandleSaveBlockEvents will handle save block events received from observer
func (eh *eventsHandler) HandleSaveBlockEvents(allEvents data.ArgsSaveBlockData) error {
	if check.IfNil(allEvents.Header) {
		return ErrNilBlockHeader
	}

	// V3 headers are handled by handleSaveBlockEventsV3, which dedupes per
	// execution-block hash instead of the outer proposed-header hash
	if allEvents.Header.IsHeaderV3() {
		return eh.handleSaveBlockEventsV3(allEvents)
	}

	blockHash := hex.EncodeToString(allEvents.HeaderHash)
	shouldProcessPushEvents := eh.shouldProcessSaveBlockEvents(blockHash)
	if !shouldProcessPushEvents {
		return nil
	}

	return eh.handleSaveBlockEventsLegacy(allEvents)
}

func (eh *eventsHandler) handleSaveBlockEventsLegacy(allEvents data.ArgsSaveBlockData) error {
	eventsData, err := eh.eventsInterceptor.ProcessBlockEvents(&allEvents)
	if err != nil {
		return err
	}

	headerTimeStamp := eventsData.Header.GetTimeStamp()
	headerTimeStampMs := allEvents.HeaderTimeStampMs
	shardID := eventsData.Header.GetShardID()
	nonce := eventsData.Header.GetNonce()

	return eh.handleSaveBlockEvents(
		eventsData,
		headerTimeStamp,
		headerTimeStampMs,
		shardID,
		nonce,
	)
}

func (eh *eventsHandler) handleSaveBlockEvents(
	eventsData *data.InterceptorBlockData,
	headerTimeStamp uint64,
	headerTimeStampMs uint64,
	shardID uint32,
	nonce uint64,
) error {
	if eventsData == nil {
		return ErrNilEventsInterceptor
	}
	if check.IfNil(eventsData.Header) {
		return ErrNilBlockHeader
	}

	t := time.Now()

	pushEvents := data.BlockEvents{
		Hash:        eventsData.Hash,
		ShardID:     shardID,
		TimeStamp:   headerTimeStamp,
		TimeStampMs: headerTimeStampMs,
		Events:      eventsData.LogEvents,
	}
	err := eh.handlePushEvents(pushEvents)
	if err != nil {
		return err
	}

	txs := data.BlockTxs{
		Hash: eventsData.Hash,
		Txs:  eventsData.Txs,
	}
	eh.handleBlockTxs(txs)

	scrs := data.BlockScrs{
		Hash: eventsData.Hash,
		Scrs: eventsData.Scrs,
	}
	eh.handleBlockScrs(scrs)

	txsWithOrder := data.BlockEventsWithOrder{
		Hash:        eventsData.Hash,
		ShardID:     shardID,
		TimeStamp:   headerTimeStamp,
		TimeStampMs: headerTimeStampMs,
		Txs:         eventsData.TxsWithOrder,
		Scrs:        eventsData.ScrsWithOrder,
		Events:      eventsData.LogEvents,
	}
	eh.handleBlockEventsWithOrder(txsWithOrder)

	stateAccesses := data.BlockStateAccesses{
		Hash:                     eventsData.Hash,
		ShardID:                  shardID,
		TimeStampMs:              headerTimeStampMs,
		Nonce:                    nonce,
		StateAccessesPerAccounts: eventsData.StateAccessesPerAccounts,
	}
	eh.handleStateAccesses(stateAccesses)

	log.Info("processed block events",
		"block hash", eventsData.Hash,
		"shard", shardID,
		"nonce", nonce,
		"num events", len(eventsData.LogEvents),
		"num txs", len(eventsData.Txs),
		"num scrs", len(eventsData.Scrs),
		"num state accesses accounts", len(eventsData.StateAccessesPerAccounts),
		"duration", time.Since(t),
	)

	return nil
}

func (eh *eventsHandler) handleSaveBlockEventsV3(allEvents data.ArgsSaveBlockData) error {
	if eh.checkDuplicates {
		lockKey := v3BatchLockPrefix + hex.EncodeToString(allEvents.HeaderHash)

		// temporary lock with defer for the execution results batch
		// this is needed to avoid concurrent triggers for partially processed batches
		acquired := eh.tryLockV3BatchWithRetry(lockKey)
		if !acquired {
			log.Debug("received duplicate v3 block events while already being processed, skipping",
				"lock key", lockKey,
			)
			return nil
		}
		defer eh.unlockV3Batch(lockKey)
	}

	executionResultsData, err := eh.eventsInterceptor.ProcessBlockEventsV3(&allEvents)
	if err != nil {
		return err
	}

	shardID := allEvents.Header.GetShardID()

	for _, executionResultData := range executionResultsData {
		shouldProcess := eh.shouldProcessSaveBlockEvents(executionResultData.Hash)
		if !shouldProcess {
			continue
		}

		timeStampSec := common.ConvertTimeStampMsToSec(executionResultData.TimeStampMs) // this is used for backwards compatibility
		err = eh.handleSaveBlockEvents(
			executionResultData,
			timeStampSec,
			executionResultData.TimeStampMs,
			shardID,
			executionResultData.Nonce,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// HandlePushEvents will handle push events received from observer
func (eh *eventsHandler) handlePushEvents(events data.BlockEvents) error {
	if events.Hash == "" {
		log.Debug("received empty hash", "event", common.PushLogsAndEvents,
			"will process", false,
		)
		return common.ErrReceivedEmptyEvents
	}

	if len(events.Events) == 0 {
		events.Events = make([]data.Event, 0)
	}

	log.Debug("received", "event", common.PushLogsAndEvents,
		"block hash", events.Hash,
		"num events", len(events.Events),
	)

	t := time.Now()
	eh.publisher.Broadcast(events)
	eh.metricsHandler.AddRequest(getRabbitOpID(common.PushLogsAndEvents), time.Since(t))
	return nil
}

func (eh *eventsHandler) shouldProcessSaveBlockEvents(blockHash string) bool {
	shouldProcessEvents := true
	if eh.checkDuplicates {
		shouldProcessEvents = eh.tryCheckProcessedWithRetry(common.PushLogsAndEvents, blockHash)
	}

	if !shouldProcessEvents {
		log.Debug("received duplicated push events, skipping",
			"block hash", blockHash,
		)

		return false
	}

	return true
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
	if eh.checkDuplicates {
		shouldProcessRevert = eh.tryCheckProcessedWithRetry(common.RevertBlockEvents, revertBlock.Hash)
	}

	if !shouldProcessRevert {
		log.Debug("received duplicated events, skipping", "event", common.RevertBlockEvents,
			"block hash", revertBlock.Hash,
		)
		return
	}

	log.Info("received", "event", common.RevertBlockEvents,
		"block hash", revertBlock.Hash,
		"shard", revertBlock.ShardID,
		"nonce", revertBlock.Nonce,
	)

	t := time.Now()
	eh.publisher.BroadcastRevert(revertBlock)
	eh.metricsHandler.AddRequest(getRabbitOpID(common.RevertBlockEvents), time.Since(t))
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
	if eh.checkDuplicates {
		shouldProcessFinalized = eh.tryCheckProcessedWithRetry(common.FinalizedBlockEvents, finalizedBlock.Hash)
	}

	if !shouldProcessFinalized {
		log.Debug("received duplicated events, skipping", "event", common.FinalizedBlockEvents,
			"block hash", finalizedBlock.Hash,
		)
		return
	}

	log.Debug("received", "event", common.FinalizedBlockEvents,
		"block hash", finalizedBlock.Hash,
	)

	t := time.Now()
	eh.publisher.BroadcastFinalized(finalizedBlock)
	eh.metricsHandler.AddRequest(getRabbitOpID(common.FinalizedBlockEvents), time.Since(t))
}

// handleBlockTxs will handle txs events received from observer
func (eh *eventsHandler) handleBlockTxs(blockTxs data.BlockTxs) {
	log.Debug("received", "event", common.BlockTxs,
		"block hash", blockTxs.Hash,
		"num txs", len(blockTxs.Txs),
	)

	t := time.Now()
	eh.publisher.BroadcastTxs(blockTxs)
	eh.metricsHandler.AddRequest(getRabbitOpID(common.BlockTxs), time.Since(t))
}

// handleBlockScrs will handle scrs events received from observer
func (eh *eventsHandler) handleBlockScrs(blockScrs data.BlockScrs) {
	log.Debug("received", "event", common.BlockScrs,
		"block hash", blockScrs.Hash,
		"num scrs", len(blockScrs.Scrs),
	)

	t := time.Now()
	eh.publisher.BroadcastScrs(blockScrs)
	eh.metricsHandler.AddRequest(getRabbitOpID(common.BlockScrs), time.Since(t))
}

// handleBlockEventsWithOrder will handle full block events received from observer
func (eh *eventsHandler) handleBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
	log.Debug("received", "event", common.BlockEvents,
		"block hash", blockTxs.Hash,
	)

	t := time.Now()
	eh.publisher.BroadcastBlockEventsWithOrder(blockTxs)
	eh.metricsHandler.AddRequest(getRabbitOpID(common.BlockEvents), time.Since(t))
}

func (eh *eventsHandler) handleStateAccesses(stateAccesses data.BlockStateAccesses) {
	log.Debug("received state accesses",
		"block hash", stateAccesses.Hash,
		"nonce", stateAccesses.Nonce,
		"stateAccesesPerAccounts num", len(stateAccesses.StateAccessesPerAccounts),
	)

	t := time.Now()
	eh.publisher.BroadcastStateAccesses(stateAccesses)
	eh.metricsHandler.AddRequest(getRabbitOpID(common.BlockStateAccesses), time.Since(t))
}

func (eh *eventsHandler) tryCheckProcessedWithRetry(id, blockHash string) bool {
	var err error
	var setSuccessful bool

	prefix := getPrefixLockerKey(id)
	key := prefix + blockHash

	for {
		t := time.Now()
		setSuccessful, err = eh.locker.IsEventProcessed(context.Background(), key)
		eh.metricsHandler.AddRequest(getRedisOpID(id), time.Since(t))

		if err == nil {
			break
		}

		hasConnection := eh.locker.HasConnection(context.Background())
		log.Error("failed to check event in locker", "error", err.Error(), "has connection", hasConnection)
		if !hasConnection {
			time.Sleep(reconnectRetryDuration)
		} else {
			time.Sleep(setRetryDuration)
		}
	}

	log.Trace("locker", "event", id, "block hash", blockHash, "succeeded", setSuccessful)

	return setSuccessful
}

func (eh *eventsHandler) tryLockV3BatchWithRetry(key string) bool {
	var err error
	var acquired bool

	for {
		acquired, err = eh.locker.TryLock(context.Background(), key)
		if err == nil {
			break
		}

		hasConnection := eh.locker.HasConnection(context.Background())
		log.Error("failed to acquire v3 batch lock", "error", err.Error(), "has connection", hasConnection)
		if !hasConnection {
			time.Sleep(reconnectRetryDuration)
		} else {
			time.Sleep(setRetryDuration)
		}
	}

	return acquired
}

func (eh *eventsHandler) unlockV3Batch(key string) {
	err := eh.locker.Unlock(context.Background(), key)
	if err != nil {
		log.Error("failed to release v3 batch lock", "error", err.Error(), "key", key)
	}
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
	}

	return ""
}

func getRabbitOpID(operation string) string {
	return fmt.Sprintf("%s-%s", rabbitmqMetricPrefix, operation)
}

func getRedisOpID(operation string) string {
	return fmt.Sprintf("%s-%s", redisMetricPrefix, operation)
}

// IsInterfaceNil returns true if there is no value under the interface
func (eh *eventsHandler) IsInterfaceNil() bool {
	return eh == nil
}
