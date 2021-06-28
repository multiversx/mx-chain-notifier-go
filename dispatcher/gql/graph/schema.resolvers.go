package graph

import (
	"context"

	"github.com/ElrondNetwork/notifier-go/dispatcher"

	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph/generated"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph/model"
)

func (r *queryResolver) Empty(ctx context.Context) (*string, error) {
	return new(string), nil
}

func (r *subscriptionResolver) Subscribe(
	ctx context.Context,
	subscriptionEntries []*model.SubscriptionEntry,
) (<-chan []*model.Event, error) {
	gqlDispatcher, _ := NewGraphqlDispatcher()
	r.hub.RegisterChan() <- gqlDispatcher

	subscribeEvent := dispatcher.SubscribeEvent{
		DispatcherID:        gqlDispatcher.GetID(),
		SubscriptionEntries: subscriptionEntriesNoPtr(subscriptionEntries),
	}
	r.hub.Subscribe(subscribeEvent)

	dispatchChan := make(chan []*model.Event, 1)

	r.MapDispatchChan(gqlDispatcher.GetID(), dispatchChan)
	gqlDispatcher.SetDispatchChan(dispatchChan)

	go func() {
		<-ctx.Done()
		r.hub.UnregisterChan() <- gqlDispatcher
		r.RemoveDispatchChan(gqlDispatcher.GetID())
	}()

	return dispatchChan, nil
}

func subscriptionEntriesNoPtr(
	subscriptionEntries []*model.SubscriptionEntry,
) []dispatcher.SubscriptionEntry {
	var subEntriesNoPtr []dispatcher.SubscriptionEntry
	for _, entry := range subscriptionEntries {
		var topicsNoPtr []string
		for _, topic := range entry.Topics {
			topicsNoPtr = append(topicsNoPtr, *topic)
		}
		subEntriesNoPtr = append(subEntriesNoPtr, dispatcher.SubscriptionEntry{
			Address:    *entry.Address,
			Identifier: *entry.Identifier,
			Topics:     topicsNoPtr,
		})
	}
	return subEntriesNoPtr
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
