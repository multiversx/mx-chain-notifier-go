package filters

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

type defaultFilter struct{}

// NewDefaultFilter creates a new default filter
func NewDefaultFilter() *defaultFilter {
	return &defaultFilter{}
}

// MatchEvent will try to match subscription data with an event
func (f *defaultFilter) MatchEvent(subscription data.Subscription, event data.Event) bool {
	switch subscription.MatchLevel {
	case dispatcher.MatchAll:
		return true
	case dispatcher.MatchAddress:
		return event.Address == subscription.Address
	case dispatcher.MatchAddressIdentifier:
		return event.Address == subscription.Address && event.Identifier == subscription.Identifier
	case dispatcher.MatchIdentifier:
		return event.Identifier == subscription.Identifier
	case dispatcher.MatchTopics:
		return f.matchTopics(subscription, event)
	default:
		return false
	}
}

func (f *defaultFilter) matchTopics(subscription data.Subscription, event data.Event) bool {
	return false
}

// IsInterfaceNil returns true if there is no value under the interface
func (f *defaultFilter) IsInterfaceNil() bool {
	return f == nil
}
