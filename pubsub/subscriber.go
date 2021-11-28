package pubsub

import (
	"context"
	"encoding/json"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

type hubSubscriber struct {
	notifierHub  dispatcher.Hub
	pubsubClient RedisClient
	rendezvous   string

	ctx context.Context
}

// NewHubSubscriber creates a new hubSubscriber instance
func NewHubSubscriber(
	ctx context.Context,
	config config.PubSubConfig,
	notifierHub dispatcher.Hub,
	pubsubClient RedisClient,
) *hubSubscriber {
	return &hubSubscriber{
		notifierHub:  notifierHub,
		pubsubClient: pubsubClient,
		rendezvous:   config.Channel,
		ctx:          ctx,
	}
}

// Subscribe launches a goroutine which listens on the pubsub message channel
func (s *hubSubscriber) Subscribe() {
	go s.subscribeToChannel()
}

func (s *hubSubscriber) subscribeToChannel() {
	subChannel := s.pubsubClient.SubscribeGetChan(s.ctx, s.rendezvous)
	
	for msg := range subChannel {
		var events data.BlockEvents

		err := json.Unmarshal([]byte(msg.Payload), &events)
		if err != nil {
			log.Debug("could not unmarshal message", "err", err.Error())
			continue
		}

		s.notifierHub.BroadcastChan() <- events
	}
}
