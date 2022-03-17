package errors

import "errors"

// ErrNilHTTPServer signals that a nil http server has been provided
var ErrNilHTTPServer = errors.New("nil http server")

// ErrNilFacadeHandler signals that a nnil dacade handler has been provided
var ErrNilFacadeHandler = errors.New("nil facade handler")
