package dispatcher

import (
	"strings"
	"sync"
)

const (
	erdTag = "erd"
)

const (
	// MatchAll signals that all events will be matched
	MatchAll = "*"

	// MatchAddress signals that events will be filtered by (address)
	MatchAddress = "match:address"

	// MatchIdentifier signals that events will be filtered by (address,identifier)
	MatchIdentifier = "match:identifier"

	// MatchTopics signals that events will be filtered by (address,identifier,[topics_pattern])
	MatchTopics = "match:topics"
)

type SubscribeEvent struct {
	Caller             EventDispatcher
	SubscriptionValues []SubscriptionValues `json:"subscriptionValues"`
}

type SubscriptionValues struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     []string `json:"topics"`
}

type Subscription struct {
	Address    string
	Identifier string
	Topics     []string
	MatchLevel string
	Subscriber EventDispatcher
}

type SubscriptionMap struct {
	rwMut         sync.RWMutex
	subscriptions map[string][]Subscription
}

func NewSubscriptionMap() *SubscriptionMap {
	return &SubscriptionMap{
		rwMut:         sync.RWMutex{},
		subscriptions: make(map[string][]Subscription),
	}
}

func (sm *SubscriptionMap) MatchSubscribeEvent(event SubscribeEvent) {
	if event.SubscriptionValues == nil || len(event.SubscriptionValues) == 0 {
		sm.appendSubscription(MatchAll, Subscription{
			Subscriber: event.Caller,
			MatchLevel: MatchAll,
		})
		return
	}

	for _, subValues := range event.SubscriptionValues {
		matchLevel := sm.matchLevelFromInput(subValues)
		subscription := Subscription{
			Address:    subValues.Address,
			Identifier: subValues.Identifier,
			Topics:     subValues.Topics,
			Subscriber: event.Caller,
			MatchLevel: matchLevel,
		}
		sm.appendSubscription(subscription.Address, subscription)
	}
}

func (sm *SubscriptionMap) Subscriptions() map[string][]Subscription {
	sm.rwMut.RLock()
	defer sm.rwMut.RUnlock()
	return sm.subscriptions
}

func (sm *SubscriptionMap) matchLevelFromInput(subValues SubscriptionValues) string {
	hasAddress := subValues.Address != "" && strings.Contains(subValues.Address, erdTag)
	hasIdentifier := subValues.Identifier != ""
	hasTopics := len(subValues.Topics) > 0

	if hasAddress && hasIdentifier && hasTopics {
		return MatchTopics
	}
	if hasAddress && hasIdentifier {
		return MatchIdentifier
	}
	if hasAddress {
		return MatchAddress
	}

	return MatchAll
}

func (sm *SubscriptionMap) appendSubscription(key string, sub Subscription) {
	sm.rwMut.Lock()
	defer sm.rwMut.Unlock()
	sm.subscriptions[key] = append(sm.subscriptions[key], sub)
}
