package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ElrondNetwork/elrond-go/marshal"
)

const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json"
	httpPost         = "POST"
)

type HttpClient interface {
	Post(route string, payload interface{}, response interface{}) error
}

type httpClient struct {
	useAuthorization bool
	baseUrl          string
	marshalizer      marshal.Marshalizer
}

type HttpClientArgs struct {
	UseAuthorization bool
	Username         string
	Password         string
	BaseUrl          string
	Marshalizer      marshal.Marshalizer
}

// NewHttpClient creates an instance of httpClient which is a wrapper for http.Client
func NewHttpClient(args HttpClientArgs) *httpClient {
	return &httpClient{
		useAuthorization: args.UseAuthorization,
		baseUrl:          args.BaseUrl,
	}
}

// Post can be used to send POST requests. It handles marshalling to/from json
func (h *httpClient) Post(
	route string,
	payload interface{},
	response interface{},
) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s%s", h.baseUrl, route)
	req, err := http.NewRequest(httpPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set(contentTypeKey, contentTypeValue)

	if h.useAuthorization {
		h.setAuthorization(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return json.Unmarshal(resBody, &response)
}

func (h *httpClient) setAuthorization(req *http.Request) {

}
