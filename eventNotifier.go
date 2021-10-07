package notifier

import (
	"encoding/hex"
	"github.com/ElrondNetwork/elrond-go-core/core"
	nodeData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/proxy/client"
)

var log = logger.GetOrCreate("outport/eventNotifier")

const (
	pushEventEndpoint = "/events/push"
)

type eventNotifier struct {
	isNilNotifier   bool
	httpClient      client.HttpClient
	marshalizer     marshal.Marshalizer
	hasher          hashing.Hasher
	pubKeyConverter core.PubkeyConverter
}

type EventNotifierArgs struct {
	HttpClient      client.HttpClient
	Marshalizer     marshal.Marshalizer
	Hasher          hashing.Hasher
	PubKeyConverter core.PubkeyConverter
}

// NewEventNotifier creates a new instance of the eventNotifier
// It implements all methods of process.Indexer
func NewEventNotifier(args EventNotifierArgs) (*eventNotifier, error) {
	return &eventNotifier{
		isNilNotifier:   false,
		httpClient:      args.HttpClient,
		marshalizer:     args.Marshalizer,
		hasher:          args.Hasher,
		pubKeyConverter: args.PubKeyConverter,
	}, nil
}

func (en *eventNotifier) SaveBlock(args *indexer.ArgsSaveBlockData) {
	log.Debug("SaveBlock called at block", "block hash", args.HeaderHash)
	if args.TransactionsPool == nil {
		return
	}

	log.Debug("checking if block has logs", "num logs", len(args.TransactionsPool.Logs))
	var logEvents []nodeData.EventHandler
	for _, handler := range args.TransactionsPool.Logs {
		if !handler.IsInterfaceNil() {
			logEvents = append(logEvents, handler.GetLogEvents()...)
		}
	}

	if len(logEvents) == 0 {
		return
	}

	var events []data.Event
	for _, eventHandler := range logEvents {
		if !eventHandler.IsInterfaceNil() {
			bech32Address := en.pubKeyConverter.Encode(eventHandler.GetAddress())
			eventIdentifier := string(eventHandler.GetIdentifier())

			log.Debug("received event from address",
				"address", bech32Address,
				"identifier", eventIdentifier,
			)

			events = append(events, data.Event{
				Address:    bech32Address,
				Identifier: eventIdentifier,
				Topics:     eventHandler.GetTopics(),
				Data:       eventHandler.GetData(),
			})
		}
	}

	log.Debug("extracted events from block logs", "num events", len(events))

	blockEvents := data.BlockEvents{
		Hash:   hex.EncodeToString(args.HeaderHash),
		Events: events,
	}

	err := en.httpClient.Post(pushEventEndpoint, blockEvents, nil)
	if err != nil {
		log.Error("error while posting event data", "err", err.Error())
	}
}

func (en *eventNotifier) RevertIndexedBlock(header nodeData.HeaderHandler, body nodeData.BodyHandler) {
	blockHash, err := core.CalculateHash(en.marshalizer, en.hasher, header)
	if err != nil {
		log.Error("could not compute block hash", "err", err.Error())
		return
	}

	revertBlock := data.RevertBlock{
		Hash:  hex.EncodeToString(blockHash),
		Nonce: header.GetNonce(),
		Round: header.GetRound(),
		Epoch: header.GetEpoch(),
	}

	_ = revertBlock
}

func (en *eventNotifier) SaveRoundsInfo(rf []*indexer.RoundInfo) {
}

func (en *eventNotifier) SaveValidatorsRating(indexID string, validatorsRatingInfo []*indexer.ValidatorRatingInfo) {
}

func (en *eventNotifier) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {
}

func (en *eventNotifier) SaveAccounts(timestamp uint64, accounts []nodeData.UserAccountHandler) {
}

func (en *eventNotifier) IsNilIndexer() bool {
	return en.isNilNotifier
}

func (en *eventNotifier) IsInterfaceNil() bool {
	return en == nil
}

func (en *eventNotifier) Close() error {
	return nil
}
