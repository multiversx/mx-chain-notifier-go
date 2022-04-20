package data

import (
	"encoding/json"

	"github.com/ElrondNetwork/elrond-go-core/data/smartContractResult"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
)

// WSEvent defines a websocket event
type WSEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// TxLog holds log data
type TxLog struct {
	Address string  `json:"scAddress"`
	Events  []Event `json:"events"`
}

// Event holds event data
type Event struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     [][]byte `json:"topics"`
	Data       []byte   `json:"data"`
}

// BlockEvents holds events data for a block
type BlockEvents struct {
	Hash   string  `json:"hash"`
	Events []Event `json:"events"`
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

type BlockTxs struct {
	Hash string                             `json:"hash"`
	Txs  map[string]transaction.Transaction `json:"txs"`
}

type BlockScrs struct {
	Hash string                                             `json:"hash"`
	Scrs map[string]smartContractResult.SmartContractResult `json:"scrs"`
}

type SaveBlockData struct {
	Hash      string                                             `json:"hash"`
	Txs       map[string]transaction.Transaction                 `json:"txs"`
	Scrs      map[string]smartContractResult.SmartContractResult `json:"scrs"`
	LogEvents []Event                                            `json:"events"`
}
