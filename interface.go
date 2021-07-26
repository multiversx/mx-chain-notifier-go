package notifier

import (
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type Indexer interface {
	SaveBlock(args *indexer.ArgsSaveBlockData)
	RevertIndexedBlock(header coreData.HeaderHandler, body coreData.BodyHandler)
	SaveRoundsInfo(roundsInfos []*indexer.RoundInfo)
	SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32)
	SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo)
	SaveAccounts(blockTimestamp uint64, acc []coreData.UserAccountHandler)
	Close() error
	IsInterfaceNil() bool
	IsNilIndexer() bool
}
