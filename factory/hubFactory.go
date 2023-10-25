package factory

import (
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/disabled"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

// CreateHub creates a common hub component
func CreateHub(apiType string) (dispatcher.Hub, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return &disabled.Hub{}, nil
	case common.WSPublisherType:
		return createHub()
	default:
		return nil, common.ErrInvalidAPIType
	}
}
