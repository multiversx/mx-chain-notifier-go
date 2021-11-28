package pubsub

import (
	"context"
	"encoding/json"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

type hubPublisher struct {
	dispatcher.Hub

	broadcast          chan data.BlockEvents
	broadcastRevert    chan data.RevertBlock
	broadcastFinalized chan data.FinalizedBlock

	pubsubClient RedisClient
	rendezvous   string

	ctx context.Context
}

// NewHubPublisher creates a new hubPublisher instance
func NewHubPublisher(
	ctx context.Context,
	config config.PubSubConfig,
	pubsubClient RedisClient,
) *hubPublisher {
	return &hubPublisher{
		broadcast:          make(chan data.BlockEvents),
		broadcastRevert:    make(chan data.RevertBlock),
		broadcastFinalized: make(chan data.FinalizedBlock),
		pubsubClient:       pubsubClient,
		rendezvous:         config.Channel,
		ctx:                ctx,
	}
}

// Run is launched as a goroutine and listens for events on the exposed channels
func (h *hubPublisher) Run() {
	for {
		select {
		case events := <-h.broadcast:
			h.publishToChannel(events)
		case revertEvent := <-h.broadcastRevert:
			_ = revertEvent
		case finalizedEvent := <-h.broadcastFinalized:
			_ = finalizedEvent
		}
	}
}

// BroadcastChan returns a receive-only channel on which events are pushed by producers
// Upon reading the channel, the hub publishes on the pubSub channel
func (h *hubPublisher) BroadcastChan() chan<- data.BlockEvents {
	return h.broadcast
}

// BroadcastRevertChan returns a receive-only channel on which revert events are pushed by producers
// Upon reading the channel, the hub does nothing
func (h *hubPublisher) BroadcastRevertChan() chan<- data.RevertBlock {
	return h.broadcastRevert
}

// BroadcastFinalizedChan returns a receive-only channel on which finalized events are pushed
// Upon reading the channel, the hub does nothing
func (h *hubPublisher) BroadcastFinalizedChan() chan<- data.FinalizedBlock {
	return h.broadcastFinalized
}

func (h *hubPublisher) publishToChannel(events data.BlockEvents) {
	marshaledEvents, err := json.Marshal(events)
	if err != nil {
		log.Debug("could not marshal events", "err", err.Error())
		return
	}

	err = h.pubsubClient.Publish(h.ctx, h.rendezvous, marshaledEvents)
	if err != nil {
		log.Debug("could not publish to channel", "err", err.Error())
	}
}
