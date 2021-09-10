package pubsub

import (
	"context"
	"encoding/json"

	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/go-redis/redis/v8"
)

type hubSubscriber struct {
	notifierHub  dispatcher.Hub
	pubsubClient *redis.Client
	rendezvous   string

	ctx context.Context
}

// NewHubSubscriber creates a new hubSubscriber instance
func NewHubSubscriber(ctx context.Context, config config.PubSubConfig, notifierHub dispatcher.Hub) *hubSubscriber {
	return &hubSubscriber{
		notifierHub:  notifierHub,
		pubsubClient: CreatePubsubClient(config),
		rendezvous:   config.Channel,
		ctx:          ctx,
	}
}

func (s *hubSubscriber) Subscribe() {
	sub := s.pubsubClient.Subscribe(s.ctx, s.rendezvous)
	channel := sub.Channel()

	var events []data.Event
	for msg := range channel {
		event := data.Event{}

		err := json.Unmarshal([]byte(msg.Payload), &event)
		if err != nil {
			log.Debug("could not unmarshal message", "err", err.Error())
			continue
		}

		events = append(events, event)
	}

	s.notifierHub.BroadcastChan() <- events
}
