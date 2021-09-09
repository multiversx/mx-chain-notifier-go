package pubsub

import (
	"context"
	"encoding/json"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/go-redis/redis/v8"
)

var log = logger.GetOrCreate("notifier/publisherHub")

type hubPublisher struct {
	dispatcher.Hub

	broadcast    chan []data.Event
	pubsubClient *redis.Client
	rendezvous   string

	ctx context.Context
}

// NewHubPublisher creates a new hubPublisher instance
func NewHubPublisher(ctx context.Context, config config.PubSubConfig) *hubPublisher {
	return &hubPublisher{
		broadcast:    make(chan []data.Event),
		pubsubClient: CreatePubsubClient(config),
		rendezvous:   config.Channel,
		ctx:          ctx,
	}
}

func (h *hubPublisher) Run() {
	for {
		select {
		case events := <-h.broadcast:
			h.publishToChannel(events)
		}
	}
}

// BroadcastChan returns a receive-only channel on which events are pushed by producers
// Upon reading the channel, the hub publishes on the pubSub channel
func (h *hubPublisher) BroadcastChan() chan<- []data.Event {
	return h.broadcast
}

func (h *hubPublisher) publishToChannel(events []data.Event) {
	marshalledEvents, err := json.Marshal(events)
	if err != nil {
		log.Debug("could not marshal events", "err", err.Error())
		return
	}

	err = h.pubsubClient.Publish(h.ctx, h.rendezvous, marshalledEvents).Err()
	if err != nil {
		log.Debug("could not publish to channel", "err", err.Error())
	}
}
