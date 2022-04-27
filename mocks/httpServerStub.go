package mocks

import "context"

// HTTPServerStub defines a stub that implements HTTPServerHandler interface
type HTTPServerStub struct {
	ListenAndServeCalled func() error
	ShutdownCalled       func(ctx context.Context) error
}

// ListenAndServe -
func (hss *HTTPServerStub) ListenAndServe() error {
	if hss.ListenAndServeCalled != nil {
		return hss.ListenAndServeCalled()
	}

	return nil
}

// Shutdown -
func (hss *HTTPServerStub) Shutdown(ctx context.Context) error {
	if hss.ShutdownCalled != nil {
		return hss.ShutdownCalled(ctx)
	}

	return nil
}
