package rabbitmq

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/streadway/amqp"
)

const (
	emptyStr = ""
)

var log = logger.GetOrCreate("rabbitmq")

// ArgsRabbitMqPublisher defines the arguments needed for rabbitmq publisher creation
type ArgsRabbitMqPublisher struct {
	Client     RabbitMqClient
	Config     config.RabbitMQConfig
	Marshaller marshal.Marshalizer
	// MaxMessageSizeInBytes is the max size of a published message, 0 means the default one
	MaxMessageSizeInBytes uint64
}

type rabbitMqPublisher struct {
	client         RabbitMqClient
	marshaller     marshal.Marshalizer
	cfg            config.RabbitMQConfig
	maxMessageSize uint64
}

// NewRabbitMqPublisher creates a new rabbitMQ publisher instance
func NewRabbitMqPublisher(args ArgsRabbitMqPublisher) (*rabbitMqPublisher, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	rp := &rabbitMqPublisher{
		cfg:            args.Config,
		client:         args.Client,
		marshaller:     args.Marshaller,
		maxMessageSize: common.GetMaxPayloadSize(args.MaxMessageSizeInBytes),
	}

	err = rp.createExchanges()
	if err != nil {
		return nil, err
	}

	return rp, nil
}

func checkArgs(args ArgsRabbitMqPublisher) error {
	if check.IfNil(args.Client) {
		return ErrNilRabbitMqClient
	}
	if check.IfNil(args.Marshaller) {
		return common.ErrNilMarshaller
	}

	if args.Config.EventsExchange.Name == "" {
		return fmt.Errorf("%w for EventsExchange", ErrInvalidRabbitMqExchangeName)
	}
	if args.Config.EventsExchange.Type == "" {
		return fmt.Errorf("%w for EventsExchange", ErrInvalidRabbitMqExchangeType)
	}
	if args.Config.RevertEventsExchange.Name == "" {
		return fmt.Errorf("%w for RevertEventsExchange", ErrInvalidRabbitMqExchangeName)
	}
	if args.Config.RevertEventsExchange.Type == "" {
		return fmt.Errorf("%w for RevertEventsExchange", ErrInvalidRabbitMqExchangeType)
	}
	if args.Config.FinalizedEventsExchange.Name == "" {
		return fmt.Errorf("%w for FinalizedEventsExchange", ErrInvalidRabbitMqExchangeName)
	}
	if args.Config.FinalizedEventsExchange.Type == "" {
		return fmt.Errorf("%w for FinalizedEventsExchange", ErrInvalidRabbitMqExchangeType)
	}
	if args.Config.BlockTxsExchange.Name == "" {
		return fmt.Errorf("%w for BlockTxsExchange", ErrInvalidRabbitMqExchangeName)
	}
	if args.Config.BlockTxsExchange.Type == "" {
		return fmt.Errorf("%w for BlockTxsExchange", ErrInvalidRabbitMqExchangeType)
	}
	if args.Config.BlockScrsExchange.Name == "" {
		return fmt.Errorf("%w for BlockScrsExchange", ErrInvalidRabbitMqExchangeName)
	}
	if args.Config.BlockScrsExchange.Type == "" {
		return fmt.Errorf("%w for BlockScrsExchange", ErrInvalidRabbitMqExchangeType)
	}
	if args.Config.BlockEventsExchange.Name == "" {
		return fmt.Errorf("%w for BlockEventsExchange", ErrInvalidRabbitMqExchangeName)
	}
	if args.Config.BlockEventsExchange.Type == "" {
		return fmt.Errorf("%w for BlockEventsExchange", ErrInvalidRabbitMqExchangeType)
	}
	if args.Config.StateAccessesExchange.Name == "" {
		return fmt.Errorf("%w for StateAccessesExchange", ErrInvalidRabbitMqExchangeName)
	}
	if args.Config.StateAccessesExchange.Type == "" {
		return fmt.Errorf("%w for StateAccessesExchange", ErrInvalidRabbitMqExchangeType)
	}

	return nil
}

// checkAndCreateExchanges creates exchanges if they are not existing already
func (rp *rabbitMqPublisher) createExchanges() error {
	err := rp.createExchange(rp.cfg.EventsExchange)
	if err != nil {
		return err
	}
	err = rp.createExchange(rp.cfg.RevertEventsExchange)
	if err != nil {
		return err
	}
	err = rp.createExchange(rp.cfg.FinalizedEventsExchange)
	if err != nil {
		return err
	}
	err = rp.createExchange(rp.cfg.BlockTxsExchange)
	if err != nil {
		return err
	}
	err = rp.createExchange(rp.cfg.BlockScrsExchange)
	if err != nil {
		return err
	}
	err = rp.createExchange(rp.cfg.BlockEventsExchange)
	if err != nil {
		return err
	}
	err = rp.createExchange(rp.cfg.StateAccessesExchange)
	if err != nil {
		return err
	}

	return nil
}

func (rp *rabbitMqPublisher) createExchange(conf config.RabbitMQExchangeConfig) error {
	err := rp.client.ExchangeDeclare(conf.Name, conf.Type)
	if err != nil {
		return err
	}

	log.Info("checked and declared rabbitMQ exchange", "name", conf.Name, "type", conf.Type)

	return nil
}

// Publish will publish logs and events to rabbitmq
func (rp *rabbitMqPublisher) Publish(events data.BlockEvents) error {
	eventsBytes, err := rp.marshaller.Marshal(events)
	if err != nil {
		log.Error("could not marshal events", "err", err.Error())
		return err
	}

	err = rp.publishFanout(rp.cfg.EventsExchange.Name, eventsBytes)
	if err != nil {
		log.Error("failed to publish events to rabbitMQ", "err", err.Error())
	}

	return err
}

// PublishRevert will publish revert event to rabbitmq
func (rp *rabbitMqPublisher) PublishRevert(revertBlock data.RevertBlock) error {
	revertBlockBytes, err := rp.marshaller.Marshal(revertBlock)
	if err != nil {
		log.Error("could not marshal revert event", "err", err.Error())
		return err
	}

	err = rp.publishFanout(rp.cfg.RevertEventsExchange.Name, revertBlockBytes)
	if err != nil {
		log.Error("failed to publish revert event to rabbitMQ", "err", err.Error())
	}

	return err
}

// PublishFinalized will publish finalized event to rabbitmq
func (rp *rabbitMqPublisher) PublishFinalized(finalizedBlock data.FinalizedBlock) error {
	finalizedBlockBytes, err := rp.marshaller.Marshal(finalizedBlock)
	if err != nil {
		log.Error("could not marshal finalized event", "err", err.Error())
		return err
	}

	err = rp.publishFanout(rp.cfg.FinalizedEventsExchange.Name, finalizedBlockBytes)
	if err != nil {
		log.Error("failed to publish finalized event to rabbitMQ", "err", err.Error())
	}

	return err
}

// PublishTxs will publish txs event to rabbitmq
func (rp *rabbitMqPublisher) PublishTxs(blockTxs data.BlockTxs) error {
	txsBlockBytes, err := rp.marshaller.Marshal(blockTxs)
	if err != nil {
		log.Error("could not marshal block txs event", "err", err.Error())
		return err
	}

	err = rp.publishFanout(rp.cfg.BlockTxsExchange.Name, txsBlockBytes)
	if err != nil {
		log.Error("failed to publish block txs event to rabbitMQ", "err", err.Error())
	}

	return err
}

// PublishScrs will publish scrs event to rabbitmq
func (rp *rabbitMqPublisher) PublishScrs(blockScrs data.BlockScrs) error {
	scrsBlockBytes, err := rp.marshaller.Marshal(blockScrs)
	if err != nil {
		log.Error("could not marshal block scrs event", "err", err.Error())
		return err
	}

	err = rp.publishFanout(rp.cfg.BlockScrsExchange.Name, scrsBlockBytes)
	if err != nil {
		log.Error("failed to publish block scrs event to rabbitMQ", "err", err.Error())
	}

	return err
}

// PublishBlockEventsWithOrder will publish block events with order to rabbitmq
func (rp *rabbitMqPublisher) PublishBlockEventsWithOrder(blockTxs data.BlockEventsWithOrder) error {
	txsBlockBytes, err := rp.marshaller.Marshal(blockTxs)
	if err != nil {
		log.Error("could not marshal block txs event", "err", err.Error())
		return err
	}

	err = rp.publishFanout(rp.cfg.BlockEventsExchange.Name, txsBlockBytes)
	if err != nil {
		log.Error("failed to publish full block events to rabbitMQ", "err", err.Error())
	}

	return err
}

// PublishStateAccesses will publish block state accesses to rabbitmq
func (rp *rabbitMqPublisher) PublishStateAccesses(stateAccesses data.BlockStateAccesses) error {
	stateAccessesBytes, err := rp.marshaller.Marshal(stateAccesses)
	if err != nil {
		log.Error("could not marshal block state accesses", "err", err.Error())
		return err
	}

	err = rp.publishFanout(rp.cfg.StateAccessesExchange.Name, stateAccessesBytes)
	if err != nil {
		log.Error("failed to publish block state accesses to rabbitMQ", "err", err.Error())
	}

	return err
}

// CheckSaveBlockMessages checks that every message of a saved block fits in the max message size,
// so that the block is rejected before any of its messages is published
func (rp *rabbitMqPublisher) CheckSaveBlockMessages(messages data.SaveBlockMessages) error {
	messagesToCheck := []struct {
		exchange string
		message  interface{}
	}{
		{exchange: rp.cfg.EventsExchange.Name, message: messages.Events},
		{exchange: rp.cfg.BlockTxsExchange.Name, message: messages.Txs},
		{exchange: rp.cfg.BlockScrsExchange.Name, message: messages.Scrs},
		{exchange: rp.cfg.BlockEventsExchange.Name, message: messages.EventsWithOrder},
		{exchange: rp.cfg.StateAccessesExchange.Name, message: messages.StateAccesses},
	}

	for _, m := range messagesToCheck {
		payload, err := rp.marshaller.Marshal(m.message)
		if err != nil {
			return err
		}

		err = rp.checkMessageSize(m.exchange, payload)
		if err != nil {
			log.Error("block messages rejected, message too large", "err", err.Error())
			return err
		}
	}

	return nil
}

func (rp *rabbitMqPublisher) checkMessageSize(exchangeName string, payload []byte) error {
	if uint64(len(payload)) > rp.maxMessageSize {
		return fmt.Errorf("%w: exchange %s, size %d, max size %d", common.ErrPayloadTooLarge, exchangeName, len(payload), rp.maxMessageSize)
	}

	return nil
}

func (rp *rabbitMqPublisher) publishFanout(exchangeName string, payload []byte) error {
	err := rp.checkMessageSize(exchangeName, payload)
	if err != nil {
		return err
	}

	return rp.client.Publish(
		exchangeName,
		emptyStr,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			Body: payload,
		},
	)
}

// Close will trigger to close rabbitmq client
func (rp *rabbitMqPublisher) Close() error {
	rp.client.Close()
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rp *rabbitMqPublisher) IsInterfaceNil() bool {
	return rp == nil
}
