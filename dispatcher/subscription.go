package dispatcher

import (
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/google/uuid"
)

var log = logger.GetOrCreate("subscription")

const (
	// MatchAll signals that all events will be matched
	MatchAll = "*"

	// MatchAddress signals that events will be filtered by (address)
	MatchAddress = "match:address"

	// MatchAddressIdentifier signals that events will be filtered by (address,identifier)
	MatchAddressIdentifier = "match:addressIdentifier"

	// MatchIdentifier signals that events will be filtered by (identifier)
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

type SubscriptionMapper struct {
	rwMut         sync.RWMutex
	subscriptions map[uuid.UUID][]Subscription
}

// NewSubscriptionMapper initializes an empty map for subscriptions
func NewSubscriptionMapper() *SubscriptionMapper {
	return &SubscriptionMapper{
		rwMut:         sync.RWMutex{},
		subscriptions: make(map[uuid.UUID][]Subscription),
	}
}

// MatchSubscribeEvent creates a subscription entry in the subscriptions map
// It assigns each SubscribeEvent a match level from the input provided
func (sm *SubscriptionMapper) MatchSubscribeEvent(event SubscribeEvent) {
	if event.SubscriptionEntries == nil || len(event.SubscriptionEntries) == 0 {
		sm.appendSubscription(Subscription{
			DispatcherID: event.DispatcherID,
			MatchLevel:   MatchAll,
		})
		log.Info("subscribed dispatcher",
			"dispatcherID", event.DispatcherID,
			"match level", MatchAll,
		)
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
		sm.appendSubscription(subscription)

		log.Info("added new subscription for dispatcher",
			"dispatcherID", event.DispatcherID,
			"match level", matchLevel,
		)
	}

	log.Info("subscribed dispatcher", "dispatcherID", event.DispatcherID)
}

// RemoveSubscriptions removes all subscriptions registered by a dispatcher
func (sm *SubscriptionMapper) RemoveSubscriptions(dispatcherID uuid.UUID) {
	sm.rwMut.Lock()
	defer sm.rwMut.Unlock()

	if _, ok := sm.subscriptions[dispatcherID]; ok {
		delete(sm.subscriptions, dispatcherID)
	}

	log.Info("unsubscribed dispatcher", "dispatcherID", dispatcherID)
}

// Subscriptions returns a slice reflecting the subscriptions present in the map
func (sm *SubscriptionMapper) Subscriptions() []Subscription {
	sm.rwMut.RLock()
	defer sm.rwMut.RUnlock()

	var subscriptions []Subscription
	for _, sub := range sm.subscriptions {
		subscriptions = append(subscriptions, sub...)
	}

	return subscriptions
}

func (sm *SubscriptionMapper) matchLevelFromInput(subEntry SubscriptionEntry) string {
	hasAddress := subEntry.Address != "" && strings.Contains(subEntry.Address, erdTag)
	hasIdentifier := subEntry.Identifier != ""
	hasTopics := len(subEntry.Topics) > 0

	if hasAddress && hasIdentifier && hasTopics {
		return MatchTopics
	}
	if hasAddress && hasIdentifier {
		return MatchAddressIdentifier
	}
	if hasIdentifier {
		return MatchIdentifier
	}
	if hasAddress {
		return MatchAddress
	}

	return MatchAll
}

func (sm *SubscriptionMapper) appendSubscription(sub Subscription) {
	sm.rwMut.Lock()
	defer sm.rwMut.Unlock()

	sm.subscriptions[sub.DispatcherID] = append(sm.subscriptions[sub.DispatcherID], sub)
}
