package data

import (
	"encoding/json"

	nodeData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/ElrondNetwork/elrond-go-core/data/receipt"
	"github.com/ElrondNetwork/elrond-go-core/data/rewardTx"
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

// LogEvent defines a log event associated with corresponding tx hash
type LogEvent struct {
	EventHandler nodeData.EventHandler
	TxHash       string
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

// BlockTxs holds the block transactions
type BlockTxs struct {
	Hash string                             `json:"hash"`
	Txs  map[string]transaction.Transaction `json:"txs"`
}

// BlockScrs holds the block smart contract results
type BlockScrs struct {
	Hash string                                             `json:"hash"`
	Scrs map[string]smartContractResult.SmartContractResult `json:"scrs"`
}

// SaveBlockData holds the block data that will be received on push events
type SaveBlockData struct {
	Hash      string                                             `json:"hash"`
	Txs       map[string]transaction.Transaction                 `json:"txs"`
	Scrs      map[string]smartContractResult.SmartContractResult `json:"scrs"`
	LogEvents []Event                                            `json:"events"`
}

// ArgsSaveBlockData will contain all information that are needed to save block data
type ArgsSaveBlockData struct {
	HeaderHash             []byte
	Body                   *block.Body
	Header                 *block.HeaderV2
	SignersIndexes         []uint64
	NotarizedHeadersHashes []string
	HeaderGasConsumption   outport.HeaderGasConsumption
	TransactionsPool       *TransactionsPool
	AlteredAccounts        map[string]*outport.AlteredAccount
	NumberOfShards         uint32
	IsImportDB             bool
}

// LogData holds the data needed for indexing logs and events
type LogData struct {
	LogHandler *transaction.Log
	TxHash     string
}

// TransactionsPool holds all types of transaction
type TransactionsPool struct {
	Txs      map[string]transaction.Transaction
	Scrs     map[string]smartContractResult.SmartContractResult
	Rewards  map[string]rewardTx.RewardTx
	Invalid  map[string]transaction.Transaction
	Receipts map[string]receipt.Receipt
	Logs     []*LogData
}
