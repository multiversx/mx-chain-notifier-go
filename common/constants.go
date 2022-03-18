package common

// TODO: comments udpate

type APIType string

const (
	WSAPIType           APIType = "notifier"
	MessageQueueAPIType APIType = "rabbit-api"
)

type HubType string

const (
	CommonHubType HubType = "common"
)
