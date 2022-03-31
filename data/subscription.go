package data

import "github.com/google/uuid"

// SubscribeEvent defines a subscription event
type SubscribeEvent struct {
	DispatcherID        uuid.UUID
	SubscriptionEntries []SubscriptionEntry `json:"subscriptionEntries"`
}

// SubscriptionEntry holds the subscription entry data
type SubscriptionEntry struct {
	EventType  string   `json:"eventType"`
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     []string `json:"topics"`
}

// Subscription holds subscription data
type Subscription struct {
	Address      string
	Identifier   string
	Topics       []string
	MatchLevel   string
	EventType    string
	DispatcherID uuid.UUID
}
