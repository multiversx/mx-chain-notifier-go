package data

import "time"

type Transaction struct {
	Hash                 string        `json:"-"`
	MBHash               string        `json:"miniBlockHash"`
	BlockHash            string        `json:"-"`
	Nonce                uint64        `json:"nonce"`
	Round                uint64        `json:"round"`
	Value                string        `json:"value"`
	Receiver             string        `json:"receiver"`
	Sender               string        `json:"sender"`
	ReceiverShard        uint32        `json:"receiverShard"`
	SenderShard          uint32        `json:"senderShard"`
	GasPrice             uint64        `json:"gasPrice"`
	GasLimit             uint64        `json:"gasLimit"`
	GasUsed              uint64        `json:"gasUsed"`
	Fee                  string        `json:"fee"`
	Data                 []byte        `json:"data"`
	Signature            string        `json:"signature"`
	Timestamp            time.Duration `json:"timestamp"`
	Status               string        `json:"status"`
	SearchOrder          uint32        `json:"searchOrder"`
	SmartContractResults []ScResult    `json:"scResults,omitempty"`
	SenderUserName       []byte        `json:"senderUsername,omitempty"`
	ReceiverUserName     []byte        `json:"receiverUsername,omitempty"`
	Log                  TxLog         `json:"-"`
	ReceiverAddressBytes []byte        `json:"-"`
}

type ScResult struct {
	Hash           string `json:"hash"`
	Nonce          uint64 `json:"nonce"`
	GasLimit       uint64 `json:"gasLimit"`
	GasPrice       uint64 `json:"gasPrice"`
	Value          string `json:"value"`
	Sender         string `json:"sender"`
	Receiver       string `json:"receiver"`
	RelayerAddr    string `json:"relayerAddr,omitempty"`
	RelayedValue   string `json:"relayedValue,omitempty"`
	Code           string `json:"code,omitempty"`
	Data           []byte `json:"data,omitempty"`
	PreTxHash      string `json:"prevTxHash"`
	OriginalTxHash string `json:"originalTxHash"`
	CallType       string `json:"callType"`
	CodeMetadata   []byte `json:"codeMetaData,omitempty"`
	ReturnMessage  string `json:"returnMessage,omitempty"`
}

type TxLog struct {
	Address string  `json:"scAddress"`
	Events  []Event `json:"events"`
}

type Event struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     []string `json:"topics"`
	Data       string   `json:"data"`
}
