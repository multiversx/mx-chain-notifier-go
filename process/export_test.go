package process

import "github.com/multiversx/mx-chain-notifier-go/data"

// TryCheckProcessedWithRetry exports internal method for testing
func (eh *eventsHandler) TryCheckProcessedWithRetry(prefix, blockHash string) bool {
	return eh.tryCheckProcessedWithRetry(prefix, blockHash)
}

// GetLogEventsFromTransactionsPool exports internal method for testing
func (ei *eventsInterceptor) GetLogEventsFromTransactionsPool(logs []*data.LogData) []data.Event {
	return ei.getLogEventsFromTransactionsPool(logs)
}
