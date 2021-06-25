package dispatcher

import (
	"github.com/ElrondNetwork/notifier-go/filters"
	"github.com/google/uuid"
	"strings"
	"sync"
)

const (
	erdTag = "erd"
)

type SubscribeEvent struct {
	DispatcherID        uuid.UUID
	SubscriptionEntries []SubscriptionEntry `json:"subscriptionEntries"`
}

type SubscriptionEntry struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     []string `json:"topics"`
}

type Subscription struct {
	Address      string
	Identifier   string
	Topics       []string
	MatchLevel   string
	DispatcherID uuid.UUID
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
	if event.SubscriptionEntries == nil || len(event.SubscriptionEntries) == 0 {
		sm.appendSubscription(filters.MatchAll, Subscription{
			DispatcherID: event.DispatcherID,
			MatchLevel:   filters.MatchAll,
		})
		return
	}

	for _, subValues := range event.SubscriptionEntries {
		matchLevel := sm.matchLevelFromInput(subValues)
		subscription := Subscription{
			Address:      subValues.Address,
			Identifier:   subValues.Identifier,
			Topics:       subValues.Topics,
			DispatcherID: event.DispatcherID,
			MatchLevel:   matchLevel,
		}
		sm.appendSubscription(subscription.Address, subscription)
	}
}

func (sm *SubscriptionMap) Subscriptions() map[string][]Subscription {
	sm.rwMut.RLock()
	defer sm.rwMut.RUnlock()
	return sm.subscriptions
}

func (sm *SubscriptionMap) matchLevelFromInput(subValues SubscriptionEntry) string {
	hasAddress := subValues.Address != "" && strings.Contains(subValues.Address, erdTag)
	hasIdentifier := subValues.Identifier != ""
	hasTopics := len(subValues.Topics) > 0

	if hasAddress && hasIdentifier && hasTopics {
		return filters.MatchTopics
	}
	if hasAddress && hasIdentifier {
		return filters.MatchIdentifier
	}
	if hasAddress {
		return filters.MatchAddress
	}

	return filters.MatchAll
}

func (sm *SubscriptionMap) appendSubscription(key string, sub Subscription) {
	sm.rwMut.Lock()
	defer sm.rwMut.Unlock()
	sm.subscriptions[key] = append(sm.subscriptions[key], sub)
}
