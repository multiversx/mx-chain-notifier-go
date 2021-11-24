# Elrond events notifier

The notifier service is a component that implements the Driver interface
declared in the [elrond-go](https://github.com/ElrondNetwork/elrond-go)
repository. This allows for one or multiple observers to push block data 
after each round. 

The notifier has two main modules:
- The factory module, which is used by an observer to initialize 
  an `eventNotifier` instance, that implements the Driver interface
- The proxy module, which exposes a REST API, is used by the 
`eventNotifier` to push events to the Hub, which then broadcasts on 
  the opened channels. 
  
## Prerequisites

In order to be able to run the notifier proxy and receive data, one 
has to setup one or multiple observers. For running an observing squad,
these [docs](https://docs.elrond.com/integrators/observing-squad/) 
cover the whole process. 

The required configs for launching an observer/s with a driver attached,
can be found [here](https://github.com/ElrondNetwork/elrond-go/blob/master/cmd/node/config/external.toml).

The supported config variables are as follows:

- `Enabled`: signals whether a driver should be attached when launching the node.
- `UseAuthorization`: signal whether to use authorization. For testing purposes it can be set to `false`.
- `ProxyUrl`: host and port on which the `eventNotifier` will push the parsed event data.
- `Username`: the username used for authorization.
- `Password`: the password used for authorization.

The corresponding config section for enabling the driver:

```toml
[EventNotifierConnector]
    # Enabled will turn on or off the event notifier connector
    Enabled = true

    # UseAuthorization signals the proxy to use authorization
    # Never run a production setup without authorization
    UseAuthorization = false

    # ProxyUrl is used to communicate with the subscriptions hub
    # The indexer instance will broadcast data using ProxyUrl
    ProxyUrl = "http://localhost:5000"

    # Username and Password need to be specified if UseAuthorization is set to true
    Username = ""

    # Password is used to authorize an observer to push event data
    Password = ""
```

## Install

Using the `cmd` package as root, execute the following commands:

- install go dependencies: `go install`
- build executable: `go build -o event-notifier`
- run `./event-notifier`

## Launching the proxy

Before launching the proxy service, it has to be configured so that it runs with the
correct environment variables.

The supported config variables are:
- `Port`: the port on which the http server listens on. Should be the same 
  as the port in the `ProxyUrl` described above.
- `HubType`: defaults to `common`. If one wants to use a custom hub implementation, 
  it can be added to the hub [factory](https://github.com/ElrondNetwork/notifier-go/blob/main/proxy/handlers/hub.go#L30-L34). 
  This can be registered with the custom name in the `config.toml` file.
channels the hub supports and uses for events broadcasting. The current implementation 
  supports `websocket` and `graphql subscriptions` for broadcasting.
- `Username`: the username used to authorize an observer. Can be left empty for `UseAuthorization = false`.
- `Password`: the password used to authorize an observer. Can be left empty for `UseAuthorization = false`.
  
The [config](https://github.com/ElrondNetwork/notifier-go/blob/main/config/config.toml) file:

```toml
[ConnectorApi]
    # The port on which the Hub listens for subscriptions
    Port = "5000"

    # The type of the hub. Options: | common | custom:<your_hub_id> |
    # Used for custom implementations for subscriptions/events filtering
    # Defaults to: common
    HubType = "common"

    # Username is the username needed to authorize an observer to push data
    Username = ""
    
    # Password is the password needed to authorize an observer to push event data
    Password = ""
```

## Subscribing

Once the proxy is launched together with the observer/s, the driver's `SaveBlock`
will be called. 

Note: For empty logs in any given block the driver won't push data, 
so subscribers won't be notified.

In order for a consumer to subscribe, it needs to select the correct
communication protocol and send a payload signalling the intention of
subscribing. This will generate a subscription for that session.

There are two types of events:
- Protocol based events, such as `ESDTTrasnfer` or `NFTCreate`
- Smart contract based events. These are defined inside a smart contract. 
  The event will automatically be assigned the smart contract address, 
  and the identifier will be the function by which it was triggered.
  
Example:

```json
{
  "address": "erd111",
  "identifier": "swapTokens",
  "topics": ["RUdMRA==", "RVRI"],
  "data": ""
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




