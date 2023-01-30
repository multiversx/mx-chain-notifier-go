package data

import (
	"github.com/multiversx/mx-chain-core-go/core"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/receipt"
	"github.com/multiversx/mx-chain-core-go/data/rewardTx"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

// SaveBlockData holds the filtered block data that will be received on push events
type SaveBlockData struct {
	Hash      string                                              `json:"hash"`
	Txs       map[string]*transaction.Transaction                 `json:"txs"`
	Scrs      map[string]*smartContractResult.SmartContractResult `json:"scrs"`
	LogEvents []Event                                             `json:"events"`
}

// InterceptorBlockData holds the block data needed for processing
type InterceptorBlockData struct {
	Hash          string
	Body          nodeData.BodyHandler
	Header        nodeData.HeaderHandler
	Txs           map[string]*transaction.Transaction
	TxsWithOrder  map[string]TransactionWrapped
	Scrs          map[string]*smartContractResult.SmartContractResult
	ScrsWithOrder map[string]SmartContractResultWrapped
	LogEvents     []Event
}

// ArgsSaveBlockData holds the block data that will be received on push events
type ArgsSaveBlockData struct {
	HeaderHash             []byte
	Body                   nodeData.BodyHandler
	Header                 nodeData.HeaderHandler
	SignersIndexes         []uint64
	NotarizedHeadersHashes []string
	HeaderGasConsumption   outport.HeaderGasConsumption
	TransactionsPool       *TransactionsPool
	AlteredAccounts        map[string]*outport.AlteredAccount
	NumberOfShards         uint32
	IsImportDB             bool
}

// ArgsSaveBlock holds block data with header type
type ArgsSaveBlock struct {
	HeaderType core.HeaderType
	ArgsSaveBlockData
}

// LogData holds the data needed for indexing logs and events
type LogData struct {
	LogHandler *transaction.Log
	TxHash     string
}

// TransactionsPool holds all types of transaction
type TransactionsPool struct {
	Txs      map[string]TransactionWrapped
	Scrs     map[string]SmartContractResultWrapped
	Rewards  map[string]RewardTxWrapped
	Invalid  map[string]TransactionWrapped
	Receipts map[string]ReceiptWrapped
	Logs     []*LogData
}

// TransactionWrapped defines a wrapper over transaction
type TransactionWrapped struct {
	TransactionHandler *transaction.Transaction
	outport.FeeInfo
	ExecutionOrder int
}

// SmartContractResultWrapped defines a wrapper over scr
type SmartContractResultWrapped struct {
	TransactionHandler *smartContractResult.SmartContractResult
	outport.FeeInfo
	ExecutionOrder int
}

// RewardTxWrapped defines a wrapper over rewardTx
type RewardTxWrapped struct {
	TransactionHandler *rewardTx.RewardTx
	outport.FeeInfo
	ExecutionOrder int
}

// ReceiptWrapped defines a wrapper over receipt
type ReceiptWrapped struct {
	TransactionHandler *receipt.Receipt
	outport.FeeInfo
	ExecutionOrder int
}
