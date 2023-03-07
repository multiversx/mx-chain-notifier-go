package process

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// logEvent defines a log event associated with corresponding tx hash
type logEvent struct {
	EventHandler nodeData.EventHandler
	TxHash       string
}

// ArgsEventsInterceptor defines the arguments needed for creating an events interceptor instance
type ArgsEventsInterceptor struct {
	PubKeyConverter core.PubkeyConverter
}

type eventsInterceptor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewEventsInterceptor creates a new eventsInterceptor instance
func NewEventsInterceptor(args ArgsEventsInterceptor) (*eventsInterceptor, error) {
	if check.IfNil(args.PubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &eventsInterceptor{
		pubKeyConverter: args.PubKeyConverter,
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
	txsWithOrder := make(map[string]*data.NotifierTransaction)
	for hash, tx := range eventsData.TransactionsPool.Txs {
		txs[hash] = tx.TransactionHandler
		txsWithOrder[hash] = &data.NotifierTransaction{
			Transaction:    tx.TransactionHandler,
			FeeInfo:        tx.FeeInfo,
			ExecutionOrder: tx.ExecutionOrder,
		}
	}

	scrs := make(map[string]*smartContractResult.SmartContractResult)
	scrsWithOrder := make(map[string]*data.NotifierSmartContractResult)
	for hash, scr := range eventsData.TransactionsPool.Scrs {
		scrs[hash] = scr.TransactionHandler
		scrsWithOrder[hash] = &data.NotifierSmartContractResult{
			SmartContractResult: scr.TransactionHandler,
			FeeInfo:             scr.FeeInfo,
			ExecutionOrder:      scr.ExecutionOrder,
		}
	}

	return &data.InterceptorBlockData{
		Hash:          hex.EncodeToString(eventsData.HeaderHash),
		Body:          eventsData.Body,
		Header:        eventsData.Header,
		Txs:           txs,
		TxsWithOrder:  txsWithOrder,
		Scrs:          scrs,
		ScrsWithOrder: scrsWithOrder,
		LogEvents:     events,
	}, nil
}

func (ei *eventsInterceptor) getLogEventsFromTransactionsPool(logs []*data.LogData) []data.Event {
	var logEvents []*logEvent
	for _, logData := range logs {
		if logData == nil {
			continue
		}
		if check.IfNil(logData.LogHandler) {
			continue
		}

		for _, eventHandler := range logData.LogHandler.GetLogEvents() {
			le := &logEvent{
				EventHandler: eventHandler,
				TxHash:       logData.TxHash,
			}

			logEvents = append(logEvents, le)
		}
	}

	if len(logEvents) == 0 {
		return nil
	}

	events := make([]data.Event, 0, len(logEvents))
	for _, event := range logEvents {
		if event == nil || check.IfNil(event.EventHandler) {
			continue
		}

		bech32Address := ei.pubKeyConverter.Encode(event.EventHandler.GetAddress())
		eventIdentifier := string(event.EventHandler.GetIdentifier())

		log.Debug("eventsInterceptor: received event from address",
			"address", bech32Address,
			"identifier", eventIdentifier,
		)

		events = append(events, data.Event{
			Address:    bech32Address,
			Identifier: eventIdentifier,
			Topics:     event.EventHandler.GetTopics(),
			Data:       event.EventHandler.GetData(),
			TxHash:     event.TxHash,
		})
	}

	return events
}

// IsInterfaceNil returns whether the interface is nil
func (ei *eventsInterceptor) IsInterfaceNil() bool {
	return ei == nil
}
