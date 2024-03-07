package ws

import "errors"

// ErrNilDispatcher signals that a nil dispatcher has been provided
var ErrNilDispatcher = errors.New("nil dispatcher")

// ErrNilWSUpgrader signals that a nil websocket upgrader has been provided
var ErrNilWSUpgrader = errors.New("nil websocket upgrader")

// ErrNilWSConn signals that a nil websocket connection has been provided
var ErrNilWSConn = errors.New("nil ws connection")
