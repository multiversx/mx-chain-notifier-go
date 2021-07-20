package dispatcher

import (
	"strings"
	"sync"

	"github.com/google/uuid"
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

// NewSubscriptionMap initializes an empty map for subscriptions
func NewSubscriptionMap() *SubscriptionMap {
	return &SubscriptionMap{
		rwMut:         sync.RWMutex{},
		subscriptions: make(map[string][]Subscription),
	}
}

// MatchSubscribeEvent creates a subscription entry in the subscriptions map
// It assigns each SubscribeEvent a match level from the input provided
func (sm *SubscriptionMap) MatchSubscribeEvent(event SubscribeEvent) {
	if event.SubscriptionEntries == nil || len(event.SubscriptionEntries) == 0 {
		sm.appendSubscription(MatchAll, Subscription{
			DispatcherID: event.DispatcherID,
			MatchLevel:   MatchAll,
		})
		return
	}

	for _, subEntry := range event.SubscriptionEntries {
		matchLevel := sm.matchLevelFromInput(subEntry)
		subscription := Subscription{
			Address:      subEntry.Address,
			Identifier:   subEntry.Identifier,
			Topics:       subEntry.Topics,
			DispatcherID: event.DispatcherID,
			MatchLevel:   matchLevel,
		}
		sm.appendSubscription(subscription.Address, subscription)
	}
}

// Subscriptions returns the subscriptions map
func (sm *SubscriptionMap) Subscriptions() map[string][]Subscription {
	sm.rwMut.RLock()
	defer sm.rwMut.RUnlock()
	return sm.subscriptions
}

func (sm *SubscriptionMap) matchLevelFromInput(subEntry SubscriptionEntry) string {
	hasAddress := subEntry.Address != "" && strings.Contains(subEntry.Address, erdTag)
	hasIdentifier := subEntry.Identifier != ""
	hasTopics := len(subEntry.Topics) > 0

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
