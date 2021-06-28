package filters

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

type BloomFilter interface {
	Set(data []byte) error
	SetMany(data [][]byte) error
	IsInSet(data []byte) bool
}

type EventFilter interface {
	MatchEvent(subscription dispatcher.Subscription, event data.Event) bool
}
