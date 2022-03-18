package factory

import (
	"context"

	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/rabbitmq"
)

func CreatePubSubHandler(apiType common.APIType, config *config.GeneralConfig) (dispatcher.Hub, error) {
	// TODO: evaluate better ways for this scenarios
	switch apiType {
	case common.MessageQueueAPIType:
		return createRabbitMqPublisher(config.RabbitMQ)
	case common.WSAPIType:
		hubType := common.HubType(config.ConnectorApi.HubType)
		return CreateCommonHub(hubType)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createRabbitMqPublisher(config config.RabbitMQConfig) (dispatcher.Hub, error) {
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
