package notifier

import (
	"encoding/hex"
	"fmt"
	"strings"

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
	pushEventEndpoint       = "/events/push"
	revertEventsEndpoint    = "/events/revert"
	finalizedEventsEndpoint = "/events/finalized"
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

// SaveBlock converts block data in order to be pushed to subscribers
func (en *eventNotifier) SaveBlock(args *indexer.ArgsSaveBlockData) error {
	log.Debug("SaveBlock called at block", "block hash", args.HeaderHash)
	if args.TransactionsPool == nil {
		return ErrNilTransactionsPool
	}

	log.Debug("checking if block has logs", "num logs", len(args.TransactionsPool.Logs))
	var logEvents []nodeData.EventHandler
	for _, handler := range args.TransactionsPool.Logs {
		if !handler.IsInterfaceNil() {
			logEvents = append(logEvents, handler.GetLogEvents()...)
		}
	}

	if len(logEvents) == 0 {
		return nil
	}

	var events []data.Event
	for _, eventHandler := range logEvents {
		if !eventHandler.IsInterfaceNil() {
			bech32Address := en.pubKeyConverter.Encode(eventHandler.GetAddress())
			eventIdentifier := string(eventHandler.GetIdentifier())

			topics := eventHandler.GetTopics()
			topicsAsHex := make([]string, 0, len(topics))
			for _, topic := range topics {
				topicsAsHex = append(topicsAsHex, hex.EncodeToString(topic))
			}

			log.Debug("received event from address",
				"address", bech32Address,
				"identifier", eventIdentifier,
				"topics", strings.Join(topicsAsHex, ", "),
				"data", eventHandler.GetData(),
			)

			events = append(events, data.Event{
				Address:    bech32Address,
				Identifier: eventIdentifier,
				Topics:     topics,
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
		return fmt.Errorf("%w in eventNotifier.SaveBlock while posting event data", err)
	}

	return nil
}

// RevertIndexedBlock converts revert data in order to be pushed to subscribers
func (en *eventNotifier) RevertIndexedBlock(header nodeData.HeaderHandler, _ nodeData.BodyHandler) error {
	blockHash, err := core.CalculateHash(en.marshalizer, en.hasher, header)
	if err != nil {
		return fmt.Errorf("%w in eventNotifier.RevertIndexedBlock while computing the block hash", err)
	}

	revertBlock := data.RevertBlock{
		Hash:  hex.EncodeToString(blockHash),
		Nonce: header.GetNonce(),
		Round: header.GetRound(),
		Epoch: header.GetEpoch(),
	}

	err = en.httpClient.Post(revertEventsEndpoint, revertBlock, nil)
	if err != nil {
		return fmt.Errorf("%w in eventNotifier.RevertIndexedBlock while posting event data", err)
	}

	return nil
}

// FinalizedBlock converts finalized block data in order to push it to subscribers
func (en *eventNotifier) FinalizedBlock(headerHash []byte) error {
	finalizedBlock := data.FinalizedBlock{
		Hash: hex.EncodeToString(headerHash),
	}

	err := en.httpClient.Post(finalizedEventsEndpoint, finalizedBlock, nil)
	if err != nil {
		return fmt.Errorf("%w in eventNotifier.FinalizedBlock while posting event data", err)
	}

	return nil
}

// SaveRoundsInfo returns nil
func (en *eventNotifier) SaveRoundsInfo(_ []*indexer.RoundInfo) error {
	return nil
}

// SaveValidatorsRating returns nil
func (en *eventNotifier) SaveValidatorsRating(_ string, _ []*indexer.ValidatorRatingInfo) error {
	return nil
}

// SaveValidatorsPubKeys returns nil
func (en *eventNotifier) SaveValidatorsPubKeys(_ map[uint32][][]byte, _ uint32) error {
	return nil
}

// SaveAccounts does nothing
func (en *eventNotifier) SaveAccounts(_ uint64, _ []nodeData.UserAccountHandler) error {
	return nil
}

// IsInterfaceNil returns whether the interface is nil
func (en *eventNotifier) IsInterfaceNil() bool {
	return en == nil
}

func (en *eventNotifier) Close() error {
	return nil
}
