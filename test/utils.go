package test

import (
	"math/rand"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

func RandStr(length int) string {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randSeed.Intn(len(charset))]
	}
	return string(b)
}

func SaveBlockArgsMock() *indexer.ArgsSaveBlockData {
	return &indexer.ArgsSaveBlockData{
		TransactionsPool: &indexer.Pool{
			Logs: map[string]data.LogHandler{
				"-": &log{},
			},
		},
	}
}

func HeaderHandler() data.HeaderHandler {
	return &header{}
}
