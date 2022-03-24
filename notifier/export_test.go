package notifier

// TryCheckProcessedWithRetry exports internal method for testing
func (eh *eventsHandler) TryCheckProcessedWithRetry(blockHash string) bool {
	return eh.tryCheckProcessedWithRetry(blockHash)
}
