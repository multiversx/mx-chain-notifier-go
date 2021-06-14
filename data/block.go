package data

import "time"

type Block struct {
	Nonce                 uint64        `json:"nonce"`
	Round                 uint64        `json:"round"`
	Epoch                 uint32        `json:"epoch"`
	Hash                  string        `json:"-"`
	MiniBlocksHashes      []string      `json:"miniBlocksHashes"`
	NotarizedBlocksHashes []string      `json:"notarizedBlocksHashes"`
	Proposer              uint64        `json:"proposer"`
	Validators            []uint64      `json:"validators"`
	PubKeyBitmap          string        `json:"pubKeyBitmap"`
	Size                  int64         `json:"size"`
	SizeTxs               int64         `json:"sizeTxs"`
	Timestamp             time.Duration `json:"timestamp"`
	StateRootHash         string        `json:"stateRootHash"`
	PrevHash              string        `json:"prevHash"`
	ShardID               uint32        `json:"shardId"`
	TxCount               uint32        `json:"txCount"`
	AccumulatedFees       string        `json:"accumulatedFees"`
	DeveloperFees         string        `json:"developerFees"`
	EpochStartBlock       bool          `json:"epochStartBlock"`
	SearchOrder           uint64        `json:"searchOrder"`
}
