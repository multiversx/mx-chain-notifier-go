package mocks

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/google/uuid"
)

type SubscriptionMapperStub struct {
	MatchSubscribeEventCalled func(event data.SubscribeEvent)
	RemoveSubscriptionsCalled func(dispatcherID uuid.UUID)
	SubscriptionsCalled       func() []data.Subscription
	IsInterfaceNilCalled      func() bool
}

func (s *SubscriptionMapperStub) MatchSubscribeEvent(event data.SubscribeEvent) {
	panic("not implemented") // TODO: Implement
}

func (s *SubscriptionMapperStub) RemoveSubscriptions(dispatcherID uuid.UUID) {
	panic("not implemented") // TODO: Implement
}

func (s *SubscriptionMapperStub) Subscriptions() []data.Subscription {
	panic("not implemented") // TODO: Implement
}

func (s *SubscriptionMapperStub) IsInterfaceNil() bool {
	panic("not implemented") // TODO: Implement
}
