package config_test

import (
	"strconv"
	"testing"

	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
)

func TestMainConfig(t *testing.T) {
	t.Parallel()

	generalMarshallerType := "json"
	adrConverterType := "bech32"
	adrConverterPrefix := "erd"
	adrConverterLength := 32
	checkDuplicates := true

	connectorAPIHost := "5000"
	connectorAPIUsername := "guest"
	connectorAPIPassword := "guest"

	redisURL := "redis://localhost:6379/0"
	redisMasterName := "mymaster"
	redisSentinelURL := "localhost:26379"
	redisConnectionType := "sentinel"
	redisTTL := 30

	rabbitmqURL := "amqp://guest:guest@localhost:5672"
	rabbitmqEventsExchangeName := "all_events"
	rabbitmqEventsExchangeType := "fanout"
	rabbitmqRevertExchangeName := "revert_events"
	rabbitmqRevertExchangeType := "fanout"

	wsConnURL := "localhost:22111"
	wsConnMode := "server"
	wsConnRetryDurationInSec := 5
	wsConnAckTimeoutInSec := 60
	wsConnMarshallerType := "gogo protobuf"

	expectedConfig := config.MainConfig{
		General: config.GeneralConfig{
			ExternalMarshaller: config.MarshallerConfig{
				Type: generalMarshallerType,
			},
			AddressConverter: config.AddressConverterConfig{
				Type:   adrConverterType,
				Prefix: adrConverterPrefix,
				Length: adrConverterLength,
			},
			CheckDuplicates: checkDuplicates,
		},
		WebSocketConnector: config.WebSocketConfig{
			Enabled:                    true,
			URL:                        wsConnURL,
			Mode:                       wsConnMode,
			RetryDurationInSec:         wsConnRetryDurationInSec,
			AcknowledgeTimeoutInSec:    wsConnAckTimeoutInSec,
			WithAcknowledge:            true,
			BlockingAckOnError:         false,
			DropMessagesIfNoConnection: false,
			DataMarshallerType:         wsConnMarshallerType,
		},
		ConnectorApi: config.ConnectorApiConfig{
			Enabled:  true,
			Host:     connectorAPIHost,
			Username: connectorAPIUsername,
			Password: connectorAPIPassword,
		},
		Redis: config.RedisConfig{
			Url:            redisURL,
			MasterName:     redisMasterName,
			SentinelUrl:    redisSentinelURL,
			ConnectionType: redisConnectionType,
			TTL:            uint32(redisTTL),
		},
		RabbitMQ: config.RabbitMQConfig{
			Url: rabbitmqURL,
			EventsExchange: config.RabbitMQExchangeConfig{
				Name: rabbitmqEventsExchangeName,
				Type: rabbitmqEventsExchangeType,
			},
			RevertEventsExchange: config.RabbitMQExchangeConfig{
				Name: rabbitmqRevertExchangeName,
				Type: rabbitmqRevertExchangeType,
			},
		},
	}

	testString := `
[General]
    # CheckDuplicates signals if the events received from observers have been already pushed to clients
    # Requires a redis instance/cluster and should be used when multiple observers push from the same shard
    CheckDuplicates = true

    # ExternalMarshaller is used for handling incoming/outcoming api requests 
    [General.ExternalMarshaller]
        Type = "` + generalMarshallerType + `"
    # InternalMarshaller is used for handling internal structs
    # This has to be mapped with the internal marshalling used for notifier outport driver
    [General.InternalMarshaller]
        Type = "` + generalMarshallerType + `"

    # Address pubkey converter config options
    [General.AddressConverter]
        Type = "` + adrConverterType + `"
        Prefix = "` + adrConverterPrefix + `"
        Length = ` + strconv.Itoa(adrConverterLength) + `

[WebSocketConnector]
    # Enabled will determine if websocket connector will be enabled or not
    Enabled = true

    # URL for the WebSocket client/server connection
    # This value represents the IP address and port number that the WebSocket client or server will use to establish a connection.
    URL = "` + wsConnURL + `"

    # This flag describes the mode to start the WebSocket connector. Can be "client" or "server"
    Mode = "` + wsConnMode + `"

    # Possible values: json, gogo protobuf. Should be compatible with mx-chain-node outport driver config
    DataMarshallerType = "` + wsConnMarshallerType + `"

    # Retry duration (receive/send ack signal) in seconds
    RetryDurationInSec = ` + strconv.Itoa(wsConnRetryDurationInSec) + ` 

    # Signals if in case of data payload processing error, we should send the ack signal or not
    BlockingAckOnError = false

    # Set to true to drop messages if there is no active WebSocket connection to send to.
    DropMessagesIfNoConnection = false

    # After a message will be sent it will wait for an ack message if this flag is enabled
    WithAcknowledge = true

    # The duration in seconds to wait for an acknowledgment message, after this time passes an error will be returned
    AcknowledgeTimeoutInSec = ` + strconv.Itoa(wsConnAckTimeoutInSec) + `

[ConnectorApi]
    # Enabled will determine if http connector will be enabled or not.
    # It will determine if http connector endpoints will be created.
    # If set to false, the web server will still be created for other endpoints (for metrics, or for WS if needed)
    Enabled = true

    # The address on which the events notifier listens for subscriptions
    # It can be specified as "localhost:5000" or only as "5000"
    Host = "` + connectorAPIHost + `"

    # Username and Password needed to authorize the connector
    # BasicAuth is enabled only for the endpoints with "Auth" flag enabled
    # in api.toml config file 
    Username = "` + connectorAPIUsername + `"
    Password = "` + connectorAPIPassword + `"

[Redis]
    # The url used to connect to a pubsub server
    Url = "` + redisURL + `"

    # The master name for failover client
    MasterName = "` + redisMasterName + `"

    # The sentinel url for failover client
    SentinelUrl = "` + redisSentinelURL + `"

    # The redis connection type. Options: | instance | sentinel |
    # instance - it will try to connect to a single redis instance
    # sentinel - it will try to connect to redis setup with master, slave and sentinel instances
    ConnectionType = "` + redisConnectionType + `"

    # Time to live (in minutes) for redis lock entry
    TTL = ` + strconv.Itoa(redisTTL) + `

[RabbitMQ]
    # The url used to connect to a rabbitMQ server
    # Note: not required for running in the notifier mode
    Url = "` + rabbitmqURL + `"

    # The exchange which holds all logs and events
    [RabbitMQ.EventsExchange]
        Name = "` + rabbitmqEventsExchangeName + `"
        Type = "` + rabbitmqEventsExchangeType + `"

    # The exchange which holds revert events
    [RabbitMQ.RevertEventsExchange]
        Name = "` + rabbitmqRevertExchangeName + `"
        Type = "` + rabbitmqRevertExchangeType + `"
	`

	config := config.MainConfig{}

	err := toml.Unmarshal([]byte(testString), &config)
	require.Nil(t, err)
	require.Equal(t, expectedConfig, config)
}

func TestAPIConfig(t *testing.T) {
	t.Parallel()

	expectedAPIConfig := config.APIRoutesConfig{
		APIPackages: map[string]config.APIPackageConfig{
			"events": {
				Routes: []config.RouteConfig{
					{
						Name: "/push",
						Open: true,
						Auth: false,
					},
					{
						Name: "/revert",
						Open: true,
						Auth: false,
					},
					{
						Name: "/finalized",
						Open: true,
						Auth: false,
					},
				},
			},
			"hub": {
				Routes: []config.RouteConfig{
					{
						Name: "/ws",
						Open: true,
					},
				},
			},
			"status": {
				Routes: []config.RouteConfig{
					{
						Name: "/metrics",
						Open: true,
					},
					{
						Name: "/prometheus-metrics",
						Open: true,
					},
				},
			},
		},
	}

	testString := `
# API routes configuration
[APIPackages]

[APIPackages.events]
    Routes = [
        { Name = "/push", Open = true, Auth = false },
        { Name = "/revert", Open = true, Auth = false },
        { Name = "/finalized", Open = true, Auth = false },
    ]

[APIPackages.hub]
    Routes = [
        { Name = "/ws", Open = true },
    ]

[APIPackages.status]
    Routes = [
        { Name = "/metrics", Open = true },
        { Name = "/prometheus-metrics", Open = true },
    ]
	`

	config := config.APIRoutesConfig{}

	err := toml.Unmarshal([]byte(testString), &config)
	require.Nil(t, err)
	require.Equal(t, expectedAPIConfig, config)
}
