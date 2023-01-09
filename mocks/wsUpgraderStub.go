package mocks

import (
	"net/http"

	"github.com/multiversx/mx-chain-notifier-go/dispatcher"
)

// WSUpgraderStub implements dispatcher.WSUpgrader interface
type WSUpgraderStub struct {
	UpgradeCalled func(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (dispatcher.WSConnection, error)
}

// Upgrade -
func (wus *WSUpgraderStub) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (dispatcher.WSConnection, error) {
	if wus.UpgradeCalled != nil {
		return wus.UpgradeCalled(w, r, responseHeader)
	}

	return nil, nil
}

// IsInterfaceNil -
func (wus *WSUpgraderStub) IsInterfaceNil() bool {
	return wus == nil
}
