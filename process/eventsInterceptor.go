package process

import (
	"encoding/hex"
	"sort"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/stateChange"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

type txWithOrder struct {
	hash  string
	index uint32
}

// logEvent defines a log event associated with corresponding tx hash
type logEvent struct {
	EventHandler nodeData.EventHandler
	TxHash       string
}

// ArgsEventsInterceptor defines the arguments needed for creating an events interceptor instance
type ArgsEventsInterceptor struct {
	PubKeyConverter      core.PubkeyConverter
	WithReadStateChanges bool
}

type eventsInterceptor struct {
	pubKeyConverter      core.PubkeyConverter
	withReadStateChanges bool
}

// NewEventsInterceptor creates a new eventsInterceptor instance
func NewEventsInterceptor(args ArgsEventsInterceptor) (*eventsInterceptor, error) {
	if check.IfNil(args.PubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &eventsInterceptor{
		pubKeyConverter:      args.PubKeyConverter,
		withReadStateChanges: args.WithReadStateChanges,
	}, nil
}

// ProcessBlockEvents will process block events data
func (ei *eventsInterceptor) ProcessBlockEvents(eventsData *data.ArgsSaveBlockData) (*data.InterceptorBlockData, error) {
	if eventsData == nil {
		return nil, ErrNilBlockEvents
	}
	if eventsData.TransactionsPool == nil {
		return nil, ErrNilTransactionsPool
	}
	if eventsData.Body == nil {
		return nil, ErrNilBlockBody
	}
	if eventsData.Header == nil {
		return nil, ErrNilBlockHeader
	}

	events := ei.getLogEventsFromTransactionsPool(eventsData.TransactionsPool.Logs)

	txs := make(map[string]*transaction.Transaction)
	for hash, tx := range eventsData.TransactionsPool.Transactions {
		txs[hash] = tx.Transaction
	}
	txsWithOrder := eventsData.TransactionsPool.Transactions

	scrs := make(map[string]*smartContractResult.SmartContractResult)
	for hash, scr := range eventsData.TransactionsPool.SmartContractResults {
		scrs[hash] = scr.SmartContractResult
	}
	scrsWithOrder := eventsData.TransactionsPool.SmartContractResults

	stateAccessesPerAccounts := ei.getStateAccessesPerAccounts(eventsData)

	return &data.InterceptorBlockData{
		Hash:                     hex.EncodeToString(eventsData.HeaderHash),
		Body:                     eventsData.Body,
		Header:                   eventsData.Header,
		Txs:                      txs,
		TxsWithOrder:             txsWithOrder,
		Scrs:                     scrs,
		ScrsWithOrder:            scrsWithOrder,
		LogEvents:                events,
		StateAccessesPerAccounts: stateAccessesPerAccounts,
	}, nil
}

func getTxsWithOrder(transactionsPool *outport.TransactionPool) []txWithOrder {
	txsWithOrderMap := make(map[string]uint32)

	for txHash, txInfo := range transactionsPool.Transactions {
		txsWithOrderMap[txHash] = txInfo.ExecutionOrder
	}
	for txHash, txInfo := range transactionsPool.SmartContractResults {
		txsWithOrderMap[txHash] = txInfo.ExecutionOrder
	}
	for txHash, txInfo := range transactionsPool.Rewards {
		txsWithOrderMap[txHash] = txInfo.ExecutionOrder
	}
	for txHash, txInfo := range transactionsPool.InvalidTxs {
		txsWithOrderMap[txHash] = txInfo.ExecutionOrder
	}

	txsWithOrder := make([]txWithOrder, 0, len(txsWithOrderMap))
	for txHash, index := range txsWithOrderMap {
		txsWithOrder = append(txsWithOrder, txWithOrder{
			hash:  txHash,
			index: index,
		})
	}

	sort.Slice(txsWithOrder, func(i, j int) bool {
		return txsWithOrder[i].index < txsWithOrder[j].index
	})

	return txsWithOrder
}

func (ei *eventsInterceptor) getStateAccessesPerAccounts(eventsData *data.ArgsSaveBlockData) map[string]*stateChange.StateAccesses {
	if eventsData.StateAccesses == nil {
		log.Warn("getStateAccessesPerAccounts failed: will return empty state accesses per accounts",
			"block hash", eventsData.HeaderHash,
			"error", ErrNilStateAccesses,
		)

		return make(map[string]*stateChange.StateAccesses)
	}

	stateAccessesPerTxs := eventsData.StateAccesses

	logStateAccessesPerTxs(stateAccessesPerTxs)

	// txs hashes with order
	txsWithOrder := getTxsWithOrder(eventsData.TransactionsPool)

	stateAccessesPerAccounts := make(map[string]*stateChange.StateAccesses)
	for _, txInfo := range txsWithOrder {
		txHash, err := hex.DecodeString(txInfo.hash)
		if err != nil {
			log.Error("failed to decode tx hash", "txHash", txInfo.hash)
			continue
		}

		stateAccessesPerTx, ok := stateAccessesPerTxs[string(txHash)]
		if !ok {
			log.Warn("did not find state accesses for tx", "txHash", txInfo.hash)
			continue
		}

		for _, stateAccess := range stateAccessesPerTx.StateAccess {
			if stateAccess.Type == stateChange.Read && !ei.withReadStateChanges {
				continue
			}

			accKey := hex.EncodeToString(stateAccess.MainTrieKey)
			_, ok := stateAccessesPerAccounts[accKey]
			if !ok {
				stateAccessesPerAccounts[accKey] = &stateChange.StateAccesses{
					StateAccess: make([]*stateChange.StateAccess, 0),
				}
			}

			stateAccessesPerAccounts[accKey].StateAccess = append(stateAccessesPerAccounts[accKey].StateAccess, stateAccess)
		}
	}

	log.Trace("getStateAccessesPerAccounts",
		"num stateAccessesPerAccounts", len(stateAccessesPerAccounts),
	)

	return stateAccessesPerAccounts
}

func logStateAccessesPerTxs(stateAccesses map[string]*stateChange.StateAccesses) {
	if log.GetLevel() > logger.LogTrace {
		return
	}

	log.Trace("getStateAccessesPerAccounts",
		"num stateAccessesPerTxs", len(stateAccesses),
	)

	for txHash, sts := range stateAccesses {
		log.Trace("stateAccessesPerTx",
			"txHash", txHash,
		)

		for _, st := range sts.StateAccess {
			log.Trace("st",
				"actionType", st.GetType(),
				"operation", st.GetOperation(),
			)
		}
	}
}

func (ei *eventsInterceptor) getLogEventsFromTransactionsPool(logs []*outport.LogData) []data.Event {
	var logEvents []*logEvent
	for _, logData := range logs {
		if logData == nil {
			continue
		}
		if check.IfNilReflect(logData.Log) {
			continue
		}

		for _, event := range logData.Log.Events {

			le := &logEvent{
				EventHandler: event,
				TxHash:       logData.TxHash,
			}

			logEvents = append(logEvents, le)
		}
	}

	if len(logEvents) == 0 {
		return make([]data.Event, 0)
	}

	events := make([]data.Event, 0, len(logEvents))
	for _, event := range logEvents {
		if event == nil || check.IfNil(event.EventHandler) {
			continue
		}

		bech32Address, err := ei.pubKeyConverter.Encode(event.EventHandler.GetAddress())
		if err != nil {
			log.Error("eventsInterceptor: failed to decode event address", "error", err)
			continue
		}
		eventIdentifier := string(event.EventHandler.GetIdentifier())

		log.Debug("eventsInterceptor: received event from address",
			"address", bech32Address,
			"identifier", eventIdentifier,
		)

		topics := event.EventHandler.GetTopics()
		if topics == nil {
			topics = make([][]byte, 0)
		}

		eventData := event.EventHandler.GetData()
		if eventData == nil {
			eventData = make([]byte, 0)
		}

		events = append(events, data.Event{
			Address:    bech32Address,
			Identifier: eventIdentifier,
			Topics:     topics,
			Data:       eventData,
			TxHash:     event.TxHash,
		})
	}

	return events
}

// IsInterfaceNil returns whether the interface is nil
func (ei *eventsInterceptor) IsInterfaceNil() bool {
	return ei == nil
}
