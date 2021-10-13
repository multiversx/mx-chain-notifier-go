package notifier

import (
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type Indexer interface {
	SaveBlock(args *indexer.ArgsSaveBlockData)
	SaveRoundsInfo(roundsInfos []*indexer.RoundInfo)
	SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32)
	RevertIndexedBlock(header coreData.HeaderHandler, body coreData.BodyHandler)
	SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo)
	SaveAccounts(blockTimestamp uint64, acc []coreData.UserAccountHandler)
	FinalizedBlock(headerHash []byte)
	IsInterfaceNil() bool
	IsNilIndexer() bool
	Close() error
}
