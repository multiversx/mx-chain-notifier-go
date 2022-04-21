package dispatcher

import (
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/data"
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

// SubscriptionMapper defines a subscriptions manager component
type SubscriptionMapper struct {
	rwMut         sync.RWMutex
	subscriptions map[uuid.UUID][]data.Subscription
}

// NewSubscriptionMapper initializes an empty map for subscriptions
func NewSubscriptionMapper() *SubscriptionMapper {
	return &SubscriptionMapper{
		rwMut:         sync.RWMutex{},
		subscriptions: make(map[uuid.UUID][]data.Subscription),
	}
}

// MatchSubscribeEvent creates a subscription entry in the subscriptions map
// It assigns each SubscribeEvent a match level from the input provided
func (sm *SubscriptionMapper) MatchSubscribeEvent(event data.SubscribeEvent) {
	if event.SubscriptionEntries == nil || len(event.SubscriptionEntries) == 0 {
		sm.appendSubscription(data.Subscription{
			DispatcherID: event.DispatcherID,
			MatchLevel:   MatchAll,
			EventType:    common.PushBlockEvents,
		})
		log.Info("subscribed dispatcher",
			"dispatcherID", event.DispatcherID,
			"match level", MatchAll,
		)
		return
	}

	for _, subEntry := range event.SubscriptionEntries {
		matchLevel := sm.matchLevelFromInput(subEntry)
		eventType := getEventType(subEntry)
		subscription := data.Subscription{
			Address:      subEntry.Address,
			Identifier:   subEntry.Identifier,
			Topics:       subEntry.Topics,
			DispatcherID: event.DispatcherID,
			MatchLevel:   matchLevel,
			EventType:    eventType,
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
func (sm *SubscriptionMapper) Subscriptions() []data.Subscription {
	sm.rwMut.RLock()
	defer sm.rwMut.RUnlock()

	var subscriptions []data.Subscription
	for _, sub := range sm.subscriptions {
		subscriptions = append(subscriptions, sub...)
	}

	return subscriptions
}

func (sm *SubscriptionMapper) matchLevelFromInput(subEntry data.SubscriptionEntry) string {
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

func (sm *SubscriptionMapper) appendSubscription(sub data.Subscription) {
	sm.rwMut.Lock()
	defer sm.rwMut.Unlock()

	sm.subscriptions[sub.DispatcherID] = append(sm.subscriptions[sub.DispatcherID], sub)
}

func getEventType(subEntry data.SubscriptionEntry) string {
	if subEntry.EventType == common.FinalizedBlockEvents ||
		subEntry.EventType == common.RevertBlockEvents ||
		subEntry.EventType == common.BlockTxs ||
		subEntry.EventType == common.BlockScrs {
		return subEntry.EventType
	}

	return common.PushBlockEvents
}

// IsInterfaceNil returns true if there is no value under the interface
func (sm *SubscriptionMapper) IsInterfaceNil() bool {
	return sm == nil
}
