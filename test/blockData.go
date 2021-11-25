package test

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go-core/data"
)

type header struct {
	data.HeaderHandler
}

func (h *header) GetNonce() uint64 {
	return 1111
}

func (h *header) GetRound() uint64 {
	return 2222
}

func (h *header) GetEpoch() uint32 {
	return 123
}

type log struct{}

func (l *log) GetAddress() []byte {
	return []byte{}
}

func (l *log) GetLogEvents() []data.EventHandler {
	return []data.EventHandler{
		&event{},
	}
}

func (l *log) IsInterfaceNil() bool {
	return false
}

type event struct{}

func (e *event) GetAddress() []byte {
	b, _ := hex.DecodeString("4f95d47324b6a6a8fbcd7727ae8d1b32fd0d32192e905aef4a92a38f6cf56111")
	return b
}

func (e *event) GetIdentifier() []byte {
	return []byte("identifier")
}

func (e *event) GetTopics() [][]byte {
	return [][]byte{
		[]byte("topic1"),
		[]byte("topic1"),
	}
}

func (e *event) GetData() []byte {
	return []byte("data")
}

func (e *event) IsInterfaceNil() bool {
	return false
}
