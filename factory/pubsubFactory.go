package factory

import (
	"context"

	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
)

// CreatePublisher creates publisher/subscriber component
func CreatePublisher(apiType common.APIType, config *config.GeneralConfig) (rabbitmq.PublisherService, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return createRabbitMqPublisher(config.RabbitMQ)
	case common.WSAPIType:
		return &disabled.Publisher{}, nil
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createRabbitMqPublisher(config config.RabbitMQConfig) (rabbitmq.PublisherService, error) {
	rabbitClient, err := rabbitmq.NewRabbitMQClient(context.Background(), config.Url)
	if err != nil {
		return nil, err
	}

	rabbitMqPublisherArgs := rabbitmq.ArgsRabbitMqPublisher{
		Client: rabbitClient,
		Config: config,
	}
	rabbitPublisher, err := rabbitmq.NewRabbitMqPublisher(rabbitMqPublisherArgs)
	if err != nil {
		return nil, err
	}

	return rabbitPublisher, nil
}
