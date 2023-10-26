package factory

import (
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
	"github.com/multiversx/mx-chain-notifier-go/process"
	"github.com/multiversx/mx-chain-notifier-go/rabbitmq"
)

// CreatePublisher creates publisher component
func CreatePublisher(
	apiType string,
	config config.MainConfig,
	marshaller marshal.Marshalizer,
	commonHub dispatcher.Hub,
) (process.Publisher, error) {
	switch apiType {
	case common.MessageQueuePublisherType:
		return createRabbitMqPublisher(config.RabbitMQ, marshaller)
	case common.WSPublisherType:
		return createWSPublisher(commonHub)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createRabbitMqPublisher(config config.RabbitMQConfig, marshaller marshal.Marshalizer) (rabbitmq.PublisherService, error) {
	rabbitClient, err := rabbitmq.NewRabbitMQClient(config.Url)
	if err != nil {
		return nil, err
	}

	rabbitMqPublisherArgs := rabbitmq.ArgsRabbitMqPublisher{
		Client:     rabbitClient,
		Config:     config,
		Marshaller: marshaller,
	}
	rabbitPublisher, err := rabbitmq.NewRabbitMqPublisher(rabbitMqPublisherArgs)
	if err != nil {
		return nil, err
	}

	publisher, err := process.NewPublisher(rabbitPublisher)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func createWSPublisher(commonHub dispatcher.Hub) (process.Publisher, error) {
	publisher, err := process.NewPublisher(commonHub)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
