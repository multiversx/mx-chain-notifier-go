package preprocess_test

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/mock"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
	"github.com/stretchr/testify/require"
)

func createMockEventsDataPreProcessorArgs() preprocess.ArgsEventsPreProcessor {
	return preprocess.ArgsEventsPreProcessor{
		Marshaller: &mock.MarshalizerMock{},
		Facade:     &mocks.FacadeStub{},
	}
}

func TestNewBaseEventsPreProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()
		args.Marshaller = nil

		dp, err := preprocess.NewBaseEventsPreProcessor(args)
		require.Nil(t, dp)
		require.Equal(t, common.ErrNilMarshaller, err)
	})

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()
		args.Facade = nil

		dp, err := preprocess.NewBaseEventsPreProcessor(args)
		require.Nil(t, dp)
		require.Equal(t, common.ErrNilFacadeHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsDataPreProcessorArgs()
		dp, err := preprocess.NewBaseEventsPreProcessor(args)
		require.Nil(t, err)
		require.NotNil(t, dp)
	})
}

func TestCreateEmptyBlockCreatorContainer(t *testing.T) {
	t.Parallel()

	cont, err := preprocess.CreateEmptyBlockCreatorContainer()
	require.Nil(t, err)

	// shard
	_, err = cont.Get(core.ShardHeaderV1)
	require.Nil(t, err)

	_, err = cont.Get(core.ShardHeaderV2)
	require.Nil(t, err)

	_, err = cont.Get(core.ShardHeaderV3)
	require.Nil(t, err)

	// meta
	_, err = cont.Get(core.MetaHeader)
	require.Nil(t, err)

	_, err = cont.Get(core.MetaHeaderV3)
	require.Nil(t, err)
}
