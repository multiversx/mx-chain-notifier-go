package process

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/data"
)

// ArgsEventsInterceptor defines the arguments needed for creating an events interceptor instance
type ArgsEventsInterceptor struct {
	PubKeyConverter core.PubkeyConverter
}

type eventsInterceptor struct {
	pubKeyConverter core.PubkeyConverter
}

// NewEventsInterceptor created a new eventsInterceptor instance
func NewEventsInterceptor(args ArgsEventsInterceptor) (*eventsInterceptor, error) {
	if check.IfNil(args.PubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &eventsInterceptor{
		pubKeyConverter: args.PubKeyConverter,
	}, nil
}

// ProcessBlockEvents will process block events data
func (ei *eventsInterceptor) ProcessBlockEvents(eventsData *data.ArgsSaveBlockData) *data.SaveBlockData {
	events := ei.getLogEventsFromTransactionsPool(eventsData.TransactionsPool.Logs)

	return &data.SaveBlockData{
		Hash:      hex.EncodeToString(eventsData.HeaderHash),
		Txs:       eventsData.TransactionsPool.Txs,
		Scrs:      eventsData.TransactionsPool.Scrs,
		LogEvents: events,
	}
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
