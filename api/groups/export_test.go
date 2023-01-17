package groups

import nodeData "github.com/multiversx/mx-chain-core-go/data"

// GetHeader -
func GetHeader(marshaledData []byte) (nodeData.HeaderHandler, error) {
	return getHeader(marshaledData)
}
