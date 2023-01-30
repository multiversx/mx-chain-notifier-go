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
// TODO: set transaction with order here also
type BlockTxs struct {
	Hash   string                             `json:"hash"`
	Txs    map[string]transaction.Transaction `json:"txs"`
	Events []Event                            `json:"events"`
}

// BlockEventsWithOrder holds the block transactions with order
type BlockEventsWithOrder struct {
	Hash      string                                  `json:"hash"`
	ShardID   uint32                                  `json:"shardID"`
	TimeStamp uint64                                  `json:"timestamp"`
	Txs       map[string]TransactionWithOrder         `json:"txs"`
	Scrs      map[string]SmartContractResultWithOrder `json:"scrs"`
	Events    []Event                                 `json:"events"`
}

// BlockScrs holds the block smart contract results
// TODO: set scr with order here also
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

// InterceptorBlockData holds the block data needed for processing
type InterceptorBlockData struct {
	Hash          string
	Body          *block.Body
	Header        *block.HeaderV2
	Txs           map[string]transaction.Transaction
	TxsWithOrder  map[string]TransactionWithOrder
	Scrs          map[string]smartContractResult.SmartContractResult
	ScrsWithOrder map[string]SmartContractResultWithOrder
	LogEvents     []Event
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
	Txs      map[string]TransactionWithOrder
	Scrs     map[string]SmartContractResultWithOrder
	Rewards  map[string]RewardTxWithOrder
	Invalid  map[string]TransactionWithOrder
	Receipts map[string]ReceiptWithOrder
	Logs     []*LogData
}

// TransactionWithOrder defines a wrapper over transaction
type TransactionWithOrder struct {
	transaction.Transaction
	ExecutionOrder int
}

// SmartContractResultWithOrder defines a wrapper over scr
type SmartContractResultWithOrder struct {
	smartContractResult.SmartContractResult
	ExecutionOrder int
}

// RewardTxWithOrder defines a wrapper over rewardTx
type RewardTxWithOrder struct {
	rewardTx.RewardTx
	ExecutionOrder int
}

// ReceiptWithOrder defines a wrapper over receipt
type ReceiptWithOrder struct {
	receipt.Receipt
	ExecutionOrder int
}
