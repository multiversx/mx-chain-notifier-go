package pubsub

import (
	"context"
	"encoding/json"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/go-redis/redis/v8"
)

type hubPublisher struct {
	dispatcher.Hub

	broadcast    chan []data.Event
	pubsubClient *redis.Client
	rendezvous   string

	ctx context.Context
}

// NewHubPublisher creates a new hubPublisher instance
func NewHubPublisher(
	ctx context.Context,
	config config.PubSubConfig,
	pubsubClient *redis.Client,
) *hubPublisher {
	return &hubPublisher{
		broadcast:    make(chan []data.Event),
		pubsubClient: pubsubClient,
		rendezvous:   config.Channel,
		ctx:          ctx,
	}
}

// Run is launched as a goroutine and listens for events on the exposed channels
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
	marshaledEvents, err := json.Marshal(events)
	if err != nil {
		log.Debug("could not marshal events", "err", err.Error())
		return
	}

	err = h.pubsubClient.Publish(h.ctx, h.rendezvous, marshaledEvents).Err()
	if err != nil {
		log.Debug("could not publish to channel", "err", err.Error())
	}
}
