package common

const (
	// WSAPIType defines a webserver api type using WebSockets
	WSAPIType string = "notifier"

	// MessageQueueAPIType defines a webserver api type using a message queueing service
	MessageQueueAPIType string = "rabbit-api"

	// CommonHubType defines the common hub type name
	CommonHubType string = "common"
)

const (
	// RedisInstanceConnType specifies a redis connection to a single instance
	RedisInstanceConnType string = "instance"

	// RedisSentinelConnType specifies a redis connection to a setup with sentinel
	RedisSentinelConnType string = "sentinel"
)

const (
	// PushLogsAndEvents defines the subscription event type for pushing block events
	PushLogsAndEvents string = "all_events"

	// BlockEvents defines the subscription event type for block info with logs and events
	BlockEvents string = "block_events"

	// RevertBlockEvents defines the subscription event type for revert block
	RevertBlockEvents string = "revert_events"

	// FinalizedBlockEvents defines the subscription event type for finalized blocks
	FinalizedBlockEvents string = "finalized_events"

	// BlockTxs defines the subscription event type for block txs
	BlockTxs string = "block_txs"

	// BlockScrs defines the subscription event type for block scrs
	BlockScrs string = "block_scrs"
)
