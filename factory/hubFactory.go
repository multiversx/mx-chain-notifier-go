package factory

import (
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/filters"
)

// CreateHub creates a common hub component
func CreateHub(apiType string) (dispatcher.Hub, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return &disabled.Hub{}, nil
	case common.WSAPIType:
		return createHub()
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createHub() (dispatcher.Hub, error) {
	args := hub.ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
	return hub.NewCommonHub(args)
}
