package graph

import (
	"context"

	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph/generated"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph/model"
)

func (r *queryResolver) Empty(ctx context.Context) (*string, error) {
	return new(string), nil
}

func (r *subscriptionResolver) Subscribe(ctx context.Context, data *model.SubscribeData) (<-chan []*model.Event, error) {
	gqlDispatcher, _ := NewGraphqlDispatcher()
	r.hub.RegisterChan() <- gqlDispatcher
	//r.hub.Subscribe(event {data ...})

	dispatchChan := make(chan []*model.Event, 1)

	r.MapDispatchChan(gqlDispatcher.GetID(), dispatchChan)
	gqlDispatcher.SetDispatchChan(dispatchChan)

	go func() {
		<-ctx.Done()
		r.RemoveDispatchChan(gqlDispatcher.GetID())
	}()

	return dispatchChan, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
