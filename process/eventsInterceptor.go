package process

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/notifier-go/data"
)

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
func (ei *eventsInterceptor) ProcessBlockEvents(eventsData *data.ArgsSaveBlockData) (*data.SaveBlockData, error) {
	if eventsData == nil {
		return nil, ErrNilBlockEvents
	}
	if eventsData.TransactionsPool == nil {
		return nil, ErrNilTransactionsPool
	}

	events := ei.getLogEventsFromTransactionsPool(eventsData.TransactionsPool.Logs)

	txs := make(map[string]transaction.Transaction)
	for hash, tx := range eventsData.TransactionsPool.Txs {
		txs[hash] = tx.GetInnerTx()
	}

	scrs := make(map[string]smartContractResult.SmartContractResult)
	for hash, scr := range eventsData.TransactionsPool.Scrs {
		scrs[hash] = scr.GetInnerScr()
	}

	return &data.SaveBlockData{
		Hash:      hex.EncodeToString(eventsData.HeaderHash),
		Txs:       txs,
		Scrs:      scrs,
		LogEvents: events,
	}, nil
}

func (ei *eventsInterceptor) getLogEventsFromTransactionsPool(logs []*data.LogData) []data.Event {
	var logEvents []*data.LogEvent
	for _, logData := range logs {
		if logData == nil {
			continue
		}
		if check.IfNil(logData.LogHandler) {
			continue
		}

		for _, eventHandler := range logData.LogHandler.GetLogEvents() {
			le := &data.LogEvent{
				EventHandler: eventHandler,
				TxHash:       hex.EncodeToString([]byte(logData.TxHash)),
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
