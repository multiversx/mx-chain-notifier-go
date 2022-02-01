package notifier

import (
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
)

type Indexer interface {
	SaveBlock(args *indexer.ArgsSaveBlockData) error
	SaveRoundsInfo(roundsInfos []*indexer.RoundInfo) error
	SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) error
	RevertIndexedBlock(header coreData.HeaderHandler, body coreData.BodyHandler) error
	SaveValidatorsRating(indexID string, infoRating []*indexer.ValidatorRatingInfo) error
	SaveAccounts(blockTimestamp uint64, acc map[string]*indexer.AlteredAccount) error
	FinalizedBlock(headerHash []byte) error
	Close() error
	IsInterfaceNil() bool
}
