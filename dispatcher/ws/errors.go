package ws

import "errors"

// ErrNilHubHandler signals that a nil hub handler has been provided
var ErrNilHubHandler = errors.New("nil hub handler")

// ErrNilWSUpgrader signals that a nil websocket upgrader has been provided
var ErrNilWSUpgrader = errors.New("nil websocket upgrader")
