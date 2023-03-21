package data

import (
	"github.com/multiversx/mx-chain-core-go/core"
	nodeData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
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
	TxsWithOrder  map[string]*NotifierTransaction
	Scrs          map[string]*smartContractResult.SmartContractResult
	ScrsWithOrder map[string]*NotifierSmartContractResult
	LogEvents     []Event
}

// ArgsSaveBlockData holds the block data that will be received on push events
type ArgsSaveBlockData struct {
	HeaderHash             []byte
	Body                   nodeData.BodyHandler
	Header                 nodeData.HeaderHandler
	SignersIndexes         []uint64
	NotarizedHeadersHashes []string
	HeaderGasConsumption   *outport.HeaderGasConsumption
	TransactionsPool       *TransactionsPool
	AlteredAccounts        map[string]*alteredAccount.AlteredAccount
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
	Txs      map[string]*NodeTransaction
	Scrs     map[string]*NodeSmartContractResult
	Rewards  map[string]*NodeRewardTx
	Invalid  map[string]*NodeTransaction
	Receipts map[string]*NodeReceipt
	Logs     []*LogData
}

// NodeTransaction defines a wrapper over transaction
type NodeTransaction struct {
	TransactionHandler *transaction.Transaction
	outport.FeeInfo
	ExecutionOrder int
}

// NodeSmartContractResult defines a wrapper over scr
type NodeSmartContractResult struct {
	TransactionHandler *smartContractResult.SmartContractResult
	outport.FeeInfo
	ExecutionOrder int
}

// NodeRewardTx defines a wrapper over rewardTx
type NodeRewardTx struct {
	TransactionHandler *rewardTx.RewardTx
	outport.FeeInfo
	ExecutionOrder int
}

// NodeReceipt defines a wrapper over receipt
type NodeReceipt struct {
	TransactionHandler *receipt.Receipt
	outport.FeeInfo
	ExecutionOrder int
}
