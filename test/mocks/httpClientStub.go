package mocks

type HttpClientStub struct {
	PostCalled func(route string, payload interface{}, resp interface{}) error
}

func (hc *HttpClientStub) Post(route string, payload interface{}, resp interface{}) error {
	if hc.PostCalled != nil {
		return hc.PostCalled(route, payload, resp)
	}

	return nil
}
