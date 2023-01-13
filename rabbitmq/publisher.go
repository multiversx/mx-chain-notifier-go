package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
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

	broadcast             chan data.BlockEvents
	broadcastRevert       chan data.RevertBlock
	broadcastFinalized    chan data.FinalizedBlock
	broadcastTxs          chan data.BlockTxs
	broadcastTxsWithOrder chan data.BlockTxsWithOrder
	broadcastScrs         chan data.BlockScrs

	cancelFunc func()
	closeChan  chan struct{}
}

// NewRabbitMqPublisher creates a new rabbitMQ publisher instance
func NewRabbitMqPublisher(args ArgsRabbitMqPublisher) (*rabbitMqPublisher, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	rp := &rabbitMqPublisher{
		broadcast:             make(chan data.BlockEvents),
		broadcastRevert:       make(chan data.RevertBlock),
		broadcastFinalized:    make(chan data.FinalizedBlock),
		broadcastTxs:          make(chan data.BlockTxs),
		broadcastScrs:         make(chan data.BlockScrs),
		broadcastTxsWithOrder: make(chan data.BlockTxsWithOrder),
		cfg:                   args.Config,
		client:                args.Client,
		closeChan:             make(chan struct{}),
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

	if args.Config.EventsExchange.Name == "" {
		return ErrInvalidRabbitMqExchangeName
	}
	if args.Config.EventsExchange.Type == "" {
		return ErrInvalidRabbitMqExchangeType
	}
	if args.Config.RevertEventsExchange.Name == "" {
		return ErrInvalidRabbitMqExchangeName
	}
	if args.Config.RevertEventsExchange.Type == "" {
		return ErrInvalidRabbitMqExchangeType
	}
	if args.Config.FinalizedEventsExchange.Name == "" {
		return ErrInvalidRabbitMqExchangeName
	}
	if args.Config.FinalizedEventsExchange.Type == "" {
		return ErrInvalidRabbitMqExchangeType
	}
	if args.Config.BlockTxsExchange.Name == "" {
		return ErrInvalidRabbitMqExchangeName
	}
	if args.Config.BlockTxsExchange.Type == "" {
		return ErrInvalidRabbitMqExchangeType
	}
	if args.Config.BlockScrsExchange.Name == "" {
		return ErrInvalidRabbitMqExchangeName
	}
	if args.Config.BlockScrsExchange.Type == "" {
		return ErrInvalidRabbitMqExchangeType
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
	err = rp.createExchange(rp.cfg.BlockTxsWithOrderExchange)
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

// Run is launched as a goroutine and listens for events on the exposed channels
func (rp *rabbitMqPublisher) Run() {
	var ctx context.Context
	ctx, rp.cancelFunc = context.WithCancel(context.Background())

	go rp.run(ctx)
}

func (rp *rabbitMqPublisher) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("RabbitMQ publisher is stopping...")
			rp.client.Close()
		case events := <-rp.broadcast:
			rp.publishToExchanges(events)
		case revertBlock := <-rp.broadcastRevert:
			rp.publishRevertToExchange(revertBlock)
		case finalizedBlock := <-rp.broadcastFinalized:
			rp.publishFinalizedToExchange(finalizedBlock)
		case blockTxs := <-rp.broadcastTxs:
			rp.publishTxsToExchange(blockTxs)
		case blockScrs := <-rp.broadcastScrs:
			rp.publishScrsToExchange(blockScrs)
		case blockTxs := <-rp.broadcastTxsWithOrder:
			rp.publishTxsWithOrderToExchange(blockTxs)
		case err := <-rp.client.ConnErrChan():
			if err != nil {
				log.Error("rabbitMQ connection failure", "err", err.Error())
				rp.client.Reconnect()
			}
		case err := <-rp.client.CloseErrChan():
			if err != nil {
				log.Error("rabbitMQ channel failure", "err", err.Error())
				rp.client.ReopenChannel()
			}
		}
	}
}

// Broadcast will handle the block events pushed by producers and sends them to rabbitMQ channel
func (rp *rabbitMqPublisher) Broadcast(events data.BlockEvents) {
	select {
	case rp.broadcast <- events:
	case <-rp.closeChan:
	}
}

// BroadcastRevert will handle the revert event pushed by producers and sends them to rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastRevert(events data.RevertBlock) {
	select {
	case rp.broadcastRevert <- events:
	case <-rp.closeChan:
	}
}

// BroadcastFinalized will handle the finalized event pushed by producers and sends them to rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastFinalized(events data.FinalizedBlock) {
	select {
	case rp.broadcastFinalized <- events:
	case <-rp.closeChan:
	}
}

// BroadcastTxs will handle the txs event pushed by producers and sends them to rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastTxs(events data.BlockTxs) {
	select {
	case rp.broadcastTxs <- events:
	case <-rp.closeChan:
	}
}

// BroadcastScrs will handle the scrs event pushed by producers and sends them to rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastScrs(events data.BlockScrs) {
	select {
	case rp.broadcastScrs <- events:
	case <-rp.closeChan:
	}
}

// BroadcastTxsWithOrder will handle the txs event pushed by producers and sends them to rabbitMQ channel
func (rp *rabbitMqPublisher) BroadcastTxsWithOrder(events data.BlockTxsWithOrder) {
	select {
	case rp.broadcastTxsWithOrder <- events:
	case <-rp.closeChan:
	}
}

func (rp *rabbitMqPublisher) publishToExchanges(events data.BlockEvents) {
	eventsBytes, err := json.Marshal(events)
	if err != nil {
		log.Error("could not marshal events", "err", err.Error())
		return
	}

	err = rp.publishFanout(rp.cfg.EventsExchange.Name, eventsBytes)
	if err != nil {
		log.Error("failed to publish events to rabbitMQ", "err", err.Error())
	}
}

func (rp *rabbitMqPublisher) publishRevertToExchange(revertBlock data.RevertBlock) {
	revertBlockBytes, err := json.Marshal(revertBlock)
	if err != nil {
		log.Error("could not marshal revert event", "err", err.Error())
		return
	}

	err = rp.publishFanout(rp.cfg.RevertEventsExchange.Name, revertBlockBytes)
	if err != nil {
		log.Error("failed to publish revert event to rabbitMQ", "err", err.Error())
	}
}

func (rp *rabbitMqPublisher) publishFinalizedToExchange(finalizedBlock data.FinalizedBlock) {
	finalizedBlockBytes, err := json.Marshal(finalizedBlock)
	if err != nil {
		log.Error("could not marshal finalized event", "err", err.Error())
		return
	}

	err = rp.publishFanout(rp.cfg.FinalizedEventsExchange.Name, finalizedBlockBytes)
	if err != nil {
		log.Error("failed to publish finalized event to rabbitMQ", "err", err.Error())
	}
}

func (rp *rabbitMqPublisher) publishTxsToExchange(blockTxs data.BlockTxs) {
	txsBlockBytes, err := json.Marshal(blockTxs)
	if err != nil {
		log.Error("could not marshal block txs event", "err", err.Error())
		return
	}

	err = rp.publishFanout(rp.cfg.BlockTxsExchange.Name, txsBlockBytes)
	if err != nil {
		log.Error("failed to publish block txs event to rabbitMQ", "err", err.Error())
	}
}

func (rp *rabbitMqPublisher) publishScrsToExchange(blockScrs data.BlockScrs) {
	scrsBlockBytes, err := json.Marshal(blockScrs)
	if err != nil {
		log.Error("could not marshal block scrs event", "err", err.Error())
		return
	}

	err = rp.publishFanout(rp.cfg.BlockScrsExchange.Name, scrsBlockBytes)
	if err != nil {
		log.Error("failed to publish block scrs event to rabbitMQ", "err", err.Error())
	}
}

func (rp *rabbitMqPublisher) publishTxsWithOrderToExchange(blockTxs data.BlockTxsWithOrder) {
	txsBlockBytes, err := json.Marshal(blockTxs)
	if err != nil {
		log.Error("could not marshal block txs event", "err", err.Error())
		return
	}

	err = rp.publishFanout(rp.cfg.BlockTxsWithOrderExchange.Name, txsBlockBytes)
	if err != nil {
		log.Error("failed to publish block txs with order event to rabbitMQ", "err", err.Error())
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

// Close will close the channels
func (rp *rabbitMqPublisher) Close() error {
	if rp.cancelFunc != nil {
		rp.cancelFunc()
	}

	close(rp.closeChan)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rp *rabbitMqPublisher) IsInterfaceNil() bool {
	return rp == nil
}
