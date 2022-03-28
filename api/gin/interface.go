package gin

import "context"

// HTTPServerHandler defines the behaviour of a http server
type HTTPServerHandler interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
