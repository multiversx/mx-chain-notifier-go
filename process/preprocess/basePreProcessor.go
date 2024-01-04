package preprocess

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/process"
)

var log = logger.GetOrCreate("preprocess")

// ArgsEventsPreProcessor defines the arguments needed to create a new events data preprocessor
type ArgsEventsPreProcessor struct {
	Marshaller marshal.Marshalizer
	Facade     process.EventsFacadeHandler
}

type baseEventsPreProcessor struct {
	marshaller        marshal.Marshalizer
	emptyBlockCreator EmptyBlockCreatorContainer
	facade            process.EventsFacadeHandler
}

// newBaseEventsPreProcessor will create a new base events data preprocessor instance
func newBaseEventsPreProcessor(args ArgsEventsPreProcessor) (*baseEventsPreProcessor, error) {
	err := checkBaseEventsPreProcessorArgs(args)
	if err != nil {
		return nil, err
	}

	dp := &baseEventsPreProcessor{
		marshaller: args.Marshaller,
		facade:     args.Facade,
	}

	emptyBlockContainer, err := createEmptyBlockCreatorContainer()
	if err != nil {
		return nil, err
	}

	dp.emptyBlockCreator = emptyBlockContainer

	return dp, nil
}

func checkBaseEventsPreProcessorArgs(args ArgsEventsPreProcessor) error {
	if check.IfNil(args.Marshaller) {
		return common.ErrNilMarshaller
	}
	if check.IfNil(args.Facade) {
		return common.ErrNilFacadeHandler
	}

	return nil
}

func (bep *baseEventsPreProcessor) getHeaderFromBytes(headerType core.HeaderType, headerBytes []byte) (header coreData.HeaderHandler, err error) {
	creator, err := bep.emptyBlockCreator.Get(headerType)
	if err != nil {
		return nil, err
	}

	return block.GetHeaderFromBytes(bep.marshaller, creator, headerBytes)
}

func createEmptyBlockCreatorContainer() (EmptyBlockCreatorContainer, error) {
	container := block.NewEmptyBlockCreatorsContainer()

	err := container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())
	if err != nil {
		return nil, err
	}

	err = container.Add(core.ShardHeaderV2, block.NewEmptyHeaderV2Creator())
	if err != nil {
		return nil, err
	}

	err = container.Add(core.MetaHeader, block.NewEmptyMetaBlockCreator())
	if err != nil {
		return nil, err
	}

	return container, nil
}
