package process

import "github.com/ElrondNetwork/notifier-go/data"

// TryCheckProcessedWithRetry exports internal method for testing
func (eh *eventsHandler) TryCheckProcessedWithRetry(blockHash string) bool {
	return eh.tryCheckProcessedWithRetry(blockHash)
}

// GetLogEventsFromTransactionsPool exports internal method for testing
func (ei *eventsInterceptor) GetLogEventsFromTransactionsPool(logs []*data.LogData) []data.Event {
	return ei.getLogEventsFromTransactionsPool(logs)
}
