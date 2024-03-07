package process

import (
	"github.com/multiversx/mx-chain-core-go/data/outport"
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
func (ei *eventsInterceptor) GetLogEventsFromTransactionsPool(logs []*outport.LogData) []data.Event {
	return ei.getLogEventsFromTransactionsPool(logs)
}
