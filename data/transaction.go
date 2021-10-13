package data

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
