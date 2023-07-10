package disabled

import "net/http"

// WSHandler defines a disabled wsHandler component
type WSHandler struct {
}

// ServeHTTP does nothing
func (wh *WSHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
}

// Close returns nil
func (wh *WSHandler) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (wh *WSHandler) IsInterfaceNil() bool {
	return wh == nil
}
