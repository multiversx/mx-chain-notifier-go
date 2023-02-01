package data

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

// WebSocketEvent defines a websocket event
type WebSocketEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Event holds event data
type Event struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     [][]byte `json:"topics"`
	Data       []byte   `json:"data"`
	TxHash     string   `json:"txHash"`
}

// BlockEvents holds events data for a block
type BlockEvents struct {
	Hash      string  `json:"hash"`
	ShardID   uint32  `json:"shardId"`
	TimeStamp uint64  `json:"timestamp"`
	Events    []Event `json:"events"`
}

// RevertBlock holds revert event data
type RevertBlock struct {
	Hash  string `json:"hash"`
	Nonce uint64 `json:"nonce"`
	Round uint64 `json:"round"`
	Epoch uint32 `json:"epoch"`
}

// FinalizedBlock holds finalized block data
type FinalizedBlock struct {
	Hash string `json:"hash"`
}

// BlockTxs holds the block transactions
type BlockTxs struct {
	Hash string                              `json:"hash"`
	Txs  map[string]*transaction.Transaction `json:"txs"`
}

// BlockScrs holds the block smart contract results
type BlockScrs struct {
	Hash string                                              `json:"hash"`
	Scrs map[string]*smartContractResult.SmartContractResult `json:"scrs"`
}

// BlockEventsWithOrder holds the block transactions with order
type BlockEventsWithOrder struct {
	Hash      string                                 `json:"hash"`
	ShardID   uint32                                 `json:"shardID"`
	TimeStamp uint64                                 `json:"timestamp"`
	Txs       map[string]*TransactionWrapped         `json:"txs"`
	Scrs      map[string]*SmartContractResultWrapped `json:"scrs"`
	Events    []Event                                `json:"events"`
}
