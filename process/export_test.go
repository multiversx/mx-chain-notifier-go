package process

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/stateChange"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// TryCheckProcessedWithRetry exports internal method for testing
func (eh *eventsHandler) TryCheckProcessedWithRetry(prefix, blockHash string) bool {
	return eh.tryCheckProcessedWithRetry(prefix, blockHash)
}

// HandlePushEvents -
func (eh *eventsHandler) HandlePushEvents(events data.BlockEvents) error {
	return eh.handlePushEvents(events)
}

// HandleBlockTxs -
func (eh *eventsHandler) HandleBlockTxs(blockTxs data.BlockTxs) {
	eh.handleBlockTxs(blockTxs)
}

// HandleBlockScrs -
func (eh *eventsHandler) HandleBlockScrs(blockScrs data.BlockScrs) {
	eh.handleBlockScrs(blockScrs)
}

// HandleBlockEventsWithOrder -
func (eh *eventsHandler) HandleBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) {
	eh.handleBlockEventsWithOrder(blockTxs)
}

// ShouldProcessSaveBlockEvents -
func (eh *eventsHandler) ShouldProcessSaveBlockEvents(blockHash string) bool {
	return eh.shouldProcessSaveBlockEvents(blockHash)
}

// GetLogEventsFromTransactionsPool exports internal method for testing
func (ei *eventsInterceptor) GetLogEventsFromTransactionsPool(logs []*transaction.LogData) []data.Event {
	return ei.getLogEventsFromTransactionsPool(logs)
}

// GetStateAccessesPerAccounts -
func (ei *eventsInterceptor) GetStateAccessesPerAccounts(eventsData *data.ArgsSaveBlockData) map[string]*stateChange.StateAccesses {
	return ei.getStateAccessesPerAccounts(eventsData, hex.EncodeToString(eventsData.HeaderHash), eventsData.TransactionsPool)
}

// GetStateAccessesPerAccountsV3 -
func (ei *eventsInterceptor) GetStateAccessesPerAccountsV3(eventsData *data.ArgsSaveBlockData) map[string]*stateChange.StateAccesses {
	return ei.getStateAccessesPerAccountsV3(eventsData, hex.EncodeToString(eventsData.HeaderHash), eventsData.TransactionsPool)
}

// BaseNilEventsDataCheks -
func BaseNilEventsDataCheks(eventsData *data.ArgsSaveBlockData) error {
	return baseNilEventsDataChecks(eventsData)
}

type txsWithOrderExportedFields struct {
	Hash   string
	Index  uint32
	TxType txType
}

// GetTxsWithOrder exports internal method for testing
func GetTxsWithOrder(transactionsPool *outport.TransactionPool) []txsWithOrderExportedFields {
	txsWithOrder := getTxsWithOrder(transactionsPool)
	txsWithExportedFields := make([]txsWithOrderExportedFields, len(txsWithOrder))
	for i, tx := range txsWithOrder {
		txsWithExportedFields[i] = txsWithOrderExportedFields{
			Hash:   tx.hash,
			Index:  tx.index,
			TxType: tx.txType,
		}
	}
	return txsWithExportedFields
}
