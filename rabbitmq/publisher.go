package rabbitmq

import (
	"encoding/json"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/streadway/amqp"
)

const (
	emptyStr = ""
)

var log = logger.GetOrCreate("rabbitmq")

// ArgsRabbitMqPublisher defines the arguments needed for rabbitmq publisher creation
type ArgsRabbitMqPublisher struct {
	Client RabbitMqClient
	Config config.RabbitMQConfig
}

type rabbitMqPublisher struct {
	client RabbitMqClient
	cfg    config.RabbitMQConfig

	broadcast          chan data.BlockEvents
	broadcastRevert    chan data.RevertBlock
	broadcastFinalized chan data.FinalizedBlock
}

// NewRabbitMqPublisher creates a new rabbitMQ publisher instance
func NewRabbitMqPublisher(args ArgsRabbitMqPublisher) (*rabbitMqPublisher, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	return &rabbitMqPublisher{
		broadcast:          make(chan data.BlockEvents),
		broadcastRevert:    make(chan data.RevertBlock),
		broadcastFinalized: make(chan data.FinalizedBlock),
		cfg:                args.Config,
		client:             args.Client,
	}, nil
}

func checkArgs(args ArgsRabbitMqPublisher) error {
	if check.IfNil(args.Client) {
		return ErrNilRabbitMqClient
	}

	return nil
}

// Run is launched as a goroutine and listens for events on the exposed channels
func (rp *rabbitMqPublisher) Run() {
	for {
		select {
		case events := <-rp.broadcast:
			rp.publishToExchanges(events)
		case revertBlock := <-rp.broadcastRevert:
			rp.publishRevertToExchange(revertBlock)
		case finalizedBlock := <-rp.broadcastFinalized:
			rp.publishFinalizedToExchange(finalizedBlock)
		}
	}
}

// BroadcastChan returns a receive-only channel on which events are pushed by producers
// Upon reading the channel, the hub publishes on the configured rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastChan() chan<- data.BlockEvents {
	return rp.broadcast
}

// BroadcastRevertChan returns a receive-only channel on which revert events are pushed by producers
// Upon reading the channel, the hub publishes on the configured rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastRevertChan() chan<- data.RevertBlock {
	return rp.broadcastRevert
}

// BroadcastFinalizedChan returns a receive-only channel on which finalized events are pushed
// Upon reading the channel, the hub publishes on the configured rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastFinalizedChan() chan<- data.FinalizedBlock {
	return rp.broadcastFinalized
}

func (rp *rabbitMqPublisher) publishToExchanges(events data.BlockEvents) {
	if rp.cfg.EventsExchange != "" {
		eventsBytes, err := json.Marshal(events)
		if err != nil {
			log.Error("could not marshal events", "err", err.Error())
			return
		}

		err = rp.publishFanout(rp.cfg.EventsExchange, eventsBytes)
		if err != nil {
			log.Error("failed to publish events to rabbitMQ", "err", err.Error())
		}
	}
}

func (rp *rabbitMqPublisher) publishRevertToExchange(revertBlock data.RevertBlock) {
	if rp.cfg.RevertEventsExchange != "" {
		revertBlockBytes, err := json.Marshal(revertBlock)
		if err != nil {
			log.Error("could not marshal revert event", "err", err.Error())
			return
		}

		err = rp.publishFanout(rp.cfg.RevertEventsExchange, revertBlockBytes)
		if err != nil {
			log.Error("failed to publish revert event to rabbitMQ", "err", err.Error())
		}
	}
}

func (rp *rabbitMqPublisher) publishFinalizedToExchange(finalizedBlock data.FinalizedBlock) {
	if rp.cfg.FinalizedEventsExchange != "" {
		finalizedBlockBytes, err := json.Marshal(finalizedBlock)
		if err != nil {
			log.Error("could not marshal finalized event", "err", err.Error())
			return
		}

		err = rp.publishFanout(rp.cfg.FinalizedEventsExchange, finalizedBlockBytes)
		if err != nil {
			log.Error("failed to publish finalized event to rabbitMQ", "err", err.Error())
		}
	}
}

func (rp *rabbitMqPublisher) publishFanout(exchangeName string, payload []byte) error {
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

// IsInterfaceNil returns true if there is no value under the interface
func (rp *rabbitMqPublisher) IsInterfaceNil() bool {
	return rp == nil
}
