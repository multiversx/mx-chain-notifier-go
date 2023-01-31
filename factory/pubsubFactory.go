package factory

import (
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/disabled"
	"github.com/multiversx/mx-chain-notifier-go/rabbitmq"
)

// CreatePublisher creates publisher component
func CreatePublisher(apiType string, config *config.GeneralConfig) (rabbitmq.PublisherService, error) {
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
	rabbitClient, err := rabbitmq.NewRabbitMQClient(config.Url)
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
