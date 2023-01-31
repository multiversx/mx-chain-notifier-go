package gin

import (
	"context"
	"net/http"
	"time"

	apiErrors "github.com/multiversx/mx-chain-notifier-go/api/errors"
)

const contextTimeout = 5 * time.Second

type httpServerWrapper struct {
	server HTTPServerHandler
}

// NewHTTPServerWrapper returns a new instance of httpServer
func NewHTTPServerWrapper(server HTTPServerHandler) (*httpServerWrapper, error) {
	if server == nil {
		return nil, apiErrors.ErrNilHTTPServer
	}

	return &httpServerWrapper{
		server: server,
	}, nil
}

// Start will handle the starting of the gin web server
func (h *httpServerWrapper) Start() {
	err := h.server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Error("could not start webserver",
				"error", err.Error(),
			)
		} else {
			log.Debug("ListenAndServe - webserver closed")
		}
	}
}

// Close will handle the stopping of the gin web server
func (h *httpServerWrapper) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return h.server.Shutdown(ctx)
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *httpServerWrapper) IsInterfaceNil() bool {
	return h == nil
}
