package filters

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

type defaultFilter struct{}

func NewDefaultFilter() *defaultFilter {
	return &defaultFilter{}
}

func (f *defaultFilter) MatchEvent(subscription dispatcher.Subscription, event data.Event) bool {
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

func (f *defaultFilter) matchTopics(subscription dispatcher.Subscription, event data.Event) bool {
	return false
}
