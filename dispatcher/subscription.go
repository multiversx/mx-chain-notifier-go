package dispatcher

import (
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/data"
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

	// MaxSubscriptionsPerDispatcher is the maximum number of distinct subscriptions a single dispatcher can hold
	MaxSubscriptionsPerDispatcher = 100
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
// Identical subscriptions are stored only once and a subscribe event that would make the
// dispatcher exceed MaxSubscriptionsPerDispatcher is rejected entirely
func (sm *SubscriptionMapper) MatchSubscribeEvent(event data.SubscribeEvent) {
	if len(event.SubscriptionEntries) > MaxSubscriptionsPerDispatcher {
		log.Warn("rejected subscribe event: too many subscription entries",
			"dispatcherID", event.DispatcherID,
			"num entries", len(event.SubscriptionEntries),
			"max", MaxSubscriptionsPerDispatcher,
		)
		return
	}

	newSubscriptions := make([]data.Subscription, 0, len(event.SubscriptionEntries))
	if len(event.SubscriptionEntries) == 0 {
		newSubscriptions = append(newSubscriptions, data.Subscription{
			DispatcherID: event.DispatcherID,
			MatchLevel:   MatchAll,
			EventType:    common.PushLogsAndEvents,
		})
	}

	for _, subEntry := range event.SubscriptionEntries {
		newSubscriptions = append(newSubscriptions, data.Subscription{
			Address:      subEntry.Address,
			Identifier:   subEntry.Identifier,
			Topics:       subEntry.Topics,
			DispatcherID: event.DispatcherID,
			MatchLevel:   sm.matchLevelFromInput(subEntry),
			EventType:    getEventType(subEntry),
		})
	}

	sm.rwMut.Lock()
	defer sm.rwMut.Unlock()

	existing := sm.subscriptions[event.DispatcherID]
	toAdd := make([]data.Subscription, 0, len(newSubscriptions))
	for _, sub := range newSubscriptions {
		if containsSubscription(existing, sub) || containsSubscription(toAdd, sub) {
			continue
		}
		toAdd = append(toAdd, sub)
	}

	if len(existing)+len(toAdd) > MaxSubscriptionsPerDispatcher {
		log.Warn("rejected subscribe event: too many subscriptions for dispatcher",
			"dispatcherID", event.DispatcherID,
			"num subscriptions", len(existing)+len(toAdd),
			"max", MaxSubscriptionsPerDispatcher,
		)
		return
	}

	sm.subscriptions[event.DispatcherID] = append(existing, toAdd...)

	log.Info("subscribed dispatcher",
		"dispatcherID", event.DispatcherID,
		"num subscriptions", len(existing)+len(toAdd),
	)
}

func containsSubscription(subscriptions []data.Subscription, sub data.Subscription) bool {
	for _, s := range subscriptions {
		if s.EventType == sub.EventType &&
			s.MatchLevel == sub.MatchLevel &&
			s.Address == sub.Address &&
			s.Identifier == sub.Identifier &&
			slices.Equal(s.Topics, sub.Topics) {
			return true
		}
	}

	return false
}

// RemoveSubscriptions removes all subscriptions registered by a dispatcher
func (sm *SubscriptionMapper) RemoveSubscriptions(dispatcherID uuid.UUID) {
	sm.rwMut.Lock()
	defer sm.rwMut.Unlock()

	if _, ok := sm.subscriptions[dispatcherID]; ok {
		delete(sm.subscriptions, dispatcherID)
	}

	log.Debug("unsubscribed dispatcher", "dispatcherID", dispatcherID)
}

// Subscriptions returns a slice reflecting the subscriptions present in the map
func (sm *SubscriptionMapper) Subscriptions() map[string][]data.Subscription {
	sm.rwMut.RLock()
	defer sm.rwMut.RUnlock()

	subscriptions := make(map[string][]data.Subscription)
	for _, sub := range sm.subscriptions {
		for _, s := range sub {
			subscriptions[s.EventType] = append(subscriptions[s.EventType], s)
		}
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

func getEventType(subEntry data.SubscriptionEntry) string {
	if subEntry.EventType == common.FinalizedBlockEvents ||
		subEntry.EventType == common.RevertBlockEvents ||
		subEntry.EventType == common.BlockTxs ||
		subEntry.EventType == common.BlockScrs ||
		subEntry.EventType == common.BlockEvents ||
		subEntry.EventType == common.BlockStateAccesses {
		return subEntry.EventType
	}

	return common.PushLogsAndEvents
}

// IsInterfaceNil returns true if there is no value under the interface
func (sm *SubscriptionMapper) IsInterfaceNil() bool {
	return sm == nil
}
