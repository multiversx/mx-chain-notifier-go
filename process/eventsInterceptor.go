package process

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/outport"
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
	for hash, tx := range eventsData.TransactionsPool.Transactions {
		txs[hash] = tx.Transaction
	}
	txsWithOrder := eventsData.TransactionsPool.Transactions

	scrs := make(map[string]*smartContractResult.SmartContractResult)
	for hash, scr := range eventsData.TransactionsPool.SmartContractResults {
		scrs[hash] = scr.SmartContractResult
	}
	scrsWithOrder := eventsData.TransactionsPool.SmartContractResults

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
