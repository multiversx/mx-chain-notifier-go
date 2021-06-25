package filters

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
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

type defaultFilter struct{}

func NewDefaultFilter() *defaultFilter {
	return &defaultFilter{}
}

func (f *defaultFilter) MatchEvent(subscription dispatcher.Subscription, event data.Event) bool {
	switch subscription.MatchLevel {
	case MatchAddress:
		return event.Address == subscription.Address
	case MatchIdentifier:
		return event.Identifier == subscription.Identifier
	case MatchTopics:
		return f.matchTopics(subscription, event)
	default:
		return false
	}
}

func (f *defaultFilter) matchTopics(subscription dispatcher.Subscription, event data.Event) bool {
	return false
}
