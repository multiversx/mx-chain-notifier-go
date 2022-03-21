package common

// APIType defines the webserver api type as string
type APIType string

const (
	// WSAPIType defines a webserver api type using WebSockets
	WSAPIType APIType = "notifier"

	// MessageQueueAPIType defines a webserver api type using a message queueing service
	MessageQueueAPIType APIType = "rabbit-api"
)

// HubType defines the hub type as string
type HubType string

const (
	// CommonHubType defines the common hub type name
	CommonHubType HubType = "common"
)
