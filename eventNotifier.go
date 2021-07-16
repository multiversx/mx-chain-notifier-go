package notifier

import (
	"github.com/ElrondNetwork/elrond-go/core/statistics"
	nodeData "github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/indexer"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/notifier-go/proxy/client"
)

type eventNotifier struct {
	isNilNotifier bool
	httpClient    client.HttpClient
	marshalizer   marshal.Marshalizer
}

type EventNotifierArgs struct {
	HttpClient  client.HttpClient
	Marshalizer marshal.Marshalizer
}

func NewEventNotifier(args EventNotifierArgs) (*eventNotifier, error) {
	return &eventNotifier{
		isNilNotifier: false,
		httpClient:    args.HttpClient,
		marshalizer:   args.Marshalizer,
	}, nil
}

func (en *eventNotifier) SaveBlock(args *indexer.ArgsSaveBlockData) {
}

func (en *eventNotifier) Close() error {
	return nil
}

func (en *eventNotifier) RevertIndexedBlock(header nodeData.HeaderHandler, body nodeData.BodyHandler) {
}

func (en *eventNotifier) SaveRoundsInfo(rf []*indexer.RoundInfo) {
}

func (en *eventNotifier) SaveValidatorsRating(indexID string, validatorsRatingInfo []*indexer.ValidatorRatingInfo) {
}

func (en *eventNotifier) SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) {
}

func (en *eventNotifier) UpdateTPS(tpsBenchmark statistics.TPSBenchmark) {
}

func (en *eventNotifier) SaveAccounts(timestamp uint64, accounts []state.UserAccountHandler) {
}

func (en *eventNotifier) SetTxLogsProcessor(txLogsProc process.TransactionLogProcessorDatabase) {
}

func (en *eventNotifier) IsNilIndexer() bool {
	return en.isNilNotifier
}

func (en *eventNotifier) IsInterfaceNil() bool {
	return en == nil
}
