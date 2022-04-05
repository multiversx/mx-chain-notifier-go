package mocks

import "net/http"

// WSHandlerStub implements WSHandler interface
type WSHandlerStub struct {
	ServeHTTPCalled func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP -
func (whs *WSHandlerStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if whs.ServeHTTPCalled != nil {
		whs.ServeHTTPCalled(w, r)
	}
}

// IsInterfaceNil -
func (whs *WSHandlerStub) IsInterfaceNil() bool {
	return whs == nil
}
