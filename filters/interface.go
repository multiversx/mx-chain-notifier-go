package filters

import (
	"github.com/multiversx/mx-chain-notifier-go/data"
)

// BloomFilter defines the behaviour of a bloom filter
type BloomFilter interface {
	Set(data []byte) error
	SetMany(data [][]byte) error
	IsInSet(data []byte) bool
}

// EventFilter defines the behaviour of an event filter component
type EventFilter interface {
	MatchEvent(subscription data.Subscription, event data.Event) bool
	IsInterfaceNil() bool
}
