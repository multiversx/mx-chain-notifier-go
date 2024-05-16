# MultiversX events notifier

The notifier service is a component that receives block events synchronously
from multiversx observer nodes, and it forwards them to a subscribing component
(via message queuing service or websockets)
  
## Prerequisites

In order to be able to run the notifier proxy and receive data, one 
has to setup one or multiple observers. For running an observing squad,
these [docs](https://docs.multiversx.com/integrators/observing-squad/) 
cover the whole process. 

The observer node can be connected using WebSocket integration. Please check observer
node config for setting up the connector. Enable `HostDriversConfig` in order to use
WebSocket integration.

The required configs for launching an observer/s with a driver attached,
can be found [here](https://github.com/multiversx/mx-chain-go/blob/master/cmd/node/config/external.toml).

The HTTP integration is still available for backwards compatibility, but it will be
deprecated in the future.

## How to run

Using the `cmd/notifier` package as root, execute the following commands:

- install go dependencies: `go install`
- build executable: `go build -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)" -o event-notifier`
- run `./event-notifier`

Or use the build script:
```bash
bash scripts/build.sh
```

---

This can also be done using a single command from `Makefile`:
```bash
# by default, rabbitmq type
# `run` command will also trigger make `build` command
make run

# specify notifier running mode (eq: rabbitmq, ws)
make run publisher_type=rabbitmq
```


CLI: run `--help` to get the command line parameters
```bash
./cmd/notifier/event-notifier --help
```

### Prerequisites

Before launching the proxy (notifier) service, it has to be configured so that it runs with the
correct config variables.

The main config file can be found [here](https://github.com/multiversx/mx-chain-notifier-go/blob/main/cmd/notifier/config/config.toml).

The supported config variables are:
- `Host`: the host and port on which the http server listens on. Should be the same port
  as the one specified in the `ProxyUrl` described above.
- `Username`: the username used to authorize an observer. Can be left empty for `UseAuthorization = false` on observer connector.
- `Password`: the password used to authorize an observer. Can be left empty for `UseAuthorization = false` on observer connector.

If observer connector is set to use BasicAuth with `UseAuthorization = true`, `Username` and `Password` has to be
set here on events notifier, and `Auth` flag has to be enabled in
[`api.toml`](https://github.com/multiversx/mx-chain-notifier-go/blob/main/cmd/notifier/config/api.toml) config file for events path.
For example:
```toml
[APIPackages.events]
    Routes = [
        { Name = "/push", Open = true, Auth = true },
        { Name = "/revert", Open = true, Auth = false },
```

After the configuration file is set up, the notifier instance can be
launched.

There are multiple publishing options when starting notifier service:
* `ws`: it will launch a websocket handler (check [WebSockets](#websockets) section)
* `rabbitmq`: it will set up a rabbitMQ client based on the RabbitMQ section from main config file (check [RabbitMQ](#rabbitmq) section)

## Development setup

There is a development setup using docker containers (with
docker compose) that can be used for this.

If it is important that no event is lost, the setup with rabbitmq (as messaging
system) and redis (as locker service) is recommended (to make sure no duplicated
events are being processed).

> If you want to use a similar setup in production systems, make sure to check
> `docker-compose.yaml` file and set up proper infrastructure and security
> configurations

* `notifier` mode can be launched as following (check `Makefile` for details): 
```bash
# Starts setup with one notifier instance
make docker-new publisher_type=ws

# Stop notifier instance
make docker-stop

# Start notifier instance
make docker-start
```

* `rabbit-api` mode can be launched as following (check `Makefile` for details): 
Manage RabbitMQ and Redis
```bash
# Starts setup with one notifier instance, redis sentinel setup and rabbitmq
make compose-new

# Stop all containers
make compose-stop

# Start start all containers
make compose-start

# Shutdown entire setup
make compose-rm
```

Start Notifier instance
```bash
make docker-new publisher_type=rabbitmq
```

### API Endpoints

Notifier service will expose several events routes, the observer nodes will
push events to these routes:
- `/events/push` (POST) -> it will handle all events for each round
- `/events/revert` (POST) -> if there is a reverted block, the event will be
  pushed on this route
- `/events/finalized` (POST) -> when the block has been finalized, the events
  will be pushed on this route

If the service will be in "notifier" mode, it will expose a additional route:
- `/hub/ws` (GET) - this route can be used to manage the websocket connection (check [websocket subscribing](#websockets) section for more details on this)

## Redis

In this setup, `Redis` is used as a locker service. If `CheckDuplicates` config
parameter is set to `true` notifier instance will check for duplicated events
in locker service.

Check `Redis` section from config in order to set up the available options.

## RabbitMQ

If `--api-type` command line parameter is set to `rabbit-api`, the notifier instance
will try to publish events to `RabbitMQ`. Check `RabbitMQ` section for config in order to
set up the url and exchanges properly.

There are multiple exchanges that can be used, they can be found in the main config file
in the `RabbitMQ` section. The data structures corresponding to these exchanges are defined
in code in `data/outport.go` file.

## Subscribing

Once the proxy is launched together with the observer/s, the driver's methods
will be called. 

### RabbitMQ

When using a setup with `RabbitMQ` you have to subscribe to each exchange
separately.

### WebSockets

In order for a consumer to subscribe, it needs to select the correct
communication protocol and send a payload signalling the intention of
subscribing. This will generate a subscription for that session.

There are two types of events:
- Protocol based events, such as `ESDTTransfer` or `NFTCreate`
- Smart contract based events. These are defined inside a smart contract. 
  The event will automatically be assigned the smart contract address, 
  and the identifier will be the function by which it was triggered.
  
Example:

```json
{
  "address": "erd111",
  "identifier": "swapTokens",
  "topics": ["RUdMRA==", "RVRI"],
}
```

Note: 
- The address field is `bech32` encoded with the tag `erd`.
- Topics are base64 encoded and require custom filters for decoding/filtering.

The subscribe message should be sent in `json` format and has the following form:

```json
{
  "subscriptionEntries": [
    {
      "address": "erd123",
      "identifier": "swapExact"
    },
    {
      "address": "erdqqq",
      "identifier": "setValue"
    }
  ]
}
```

Each subscription upon creation is assigned a `MatchLevel`:
- Match all `*`. All events are broadcast.
- Match by `address`. Events are filtered by address.
- Match by `address && identifier`. Events are filtered by (address, identifier).
- Match by `topics`. Filtering is done by topics, it currently requires custom filter implementation.

The `MatchLevel` is assigned using the input payload sent while subscribing. Examples:

- Match `*`:
```json
{
  "subscriptionEntries": []
}
```

- Match `address`:
```json
{
  "subscriptionEntries": [
    {
      "address": "erdFirst"
    },
    {
      "address": "erdSecond"
    }
  ]
}
```

- Match `address && identifier`:
```json
{
  "subscriptionEntries": [
    {
      "address": "erdFirst",
      "identifier": "ESDTTransfer"
    },
    {
      "address": "erdSecond",
      "identifier": "setValue"
    }
  ]
}
```

The subscription entry has also a field for specifying event type, which can be
one of the followings: `all_events`, `revert_events`, `finalized_events`.  By
default, it is set to `all_events`, for backwards compatibility reasons.

All other fields, like `address`, `identifier`, `topics` can be used for
`all_events`.  The other events type (`revert_events` and `finalized_events`)
do not have these fields associated with them.

A subscription example with `eventType` will be like this:
```json
{
  "subscriptionEntries": [
    {
      "eventType": "all_events",
      "address": "erdFirst",
      "identifier": "ESDTTransfer"
    },
    {
      "eventType": "finalized_events",
    }
  ]
}
```

Note: if eventType type is not specified, it will be set to `all_events` by
default.

The payload data will consist of a marshalled object containing the event type and the
inner marshalled data like:
```json
{
  "type": "all_events",
  "data": "<< marshalled object here >>",
}
```

There are multiple event types available, they can be found as constants in common package,
[constants](https://github.com/multiversx/mx-chain-notifier-go/blob/main/common/constants.go). Below there is the event type together with the associated marshalled data type.
- `all_events`
```json
{
  "hash": "blockHash1",
  "events": [
    {
      "address": "addr1",
      "identifier": "identifier",
      ...
    }
  ]
}
```

- `revert_events`
```json
{
    "hash": "blockHash1",
    "nonce": 11,
    "round": 2,
    "epoch": 1,
}
```

- `finalized_events`
```json
{
    "hash": "blockHash"
}
```

- `block_txs`:
```json
{
  "hash": "blockHash1",
  "txs": {
    "txHash1": {      
        "Nonce": 123,
        "Round": 1,
        "Epoch": 1,
        ...
    }
  }
}
```

- `block_scrs`
```json
{
  "hash": "blockHash1",
  "scrs": {
    "scrHash1": {      
        "Nonce": 123,
        ...
    }
  }
}
```
