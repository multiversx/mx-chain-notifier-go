package errors

import "errors"

// ErrNilHTTPServer signals that a nil http server has been provided
var ErrNilHTTPServer = errors.New("nil http server")

// ErrNilFacadeHandler signals that a nil facade handler has been provided
var ErrNilFacadeHandler = errors.New("nil facade handler")

// ErrNilPayloadHandler signals that a nil payload handler has been provided
var ErrNilPayloadHandler = errors.New("nil payload handler")
