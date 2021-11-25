package notifier

import (
	"errors"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/mock"
	"github.com/ElrondNetwork/elrond-go-core/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go-core/data/indexer"
	"github.com/ElrondNetwork/notifier-go/test"
	"github.com/ElrondNetwork/notifier-go/test/mocks"
	"github.com/stretchr/testify/require"
)

var pubkeyConv = func() core.PubkeyConverter {
	conv, _ := pubkeyConverter.NewBech32PubkeyConverter(32, log)
	return conv
}

func TestEventNotifier_SaveBlockNilTransactionPool(t *testing.T) {
	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient:      &mocks.HttpClientStub{},
		PubKeyConverter: pubkeyConv(),
	})

	err := en.SaveBlock(&indexer.ArgsSaveBlockData{})
	require.Equal(t, ErrNilTransactionsPool, err)
}

func TestEventNotifier_SaveBlockPostShouldWork(t *testing.T) {
	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient: &mocks.HttpClientStub{
			PostCalled: func(route string, payload interface{}, resp interface{}) error {
				require.True(t, route == pushEventEndpoint)
				return nil
			},
		},
		PubKeyConverter: pubkeyConv(),
	})

	err := en.SaveBlock(test.SaveBlockArgsMock())
	require.Nil(t, err)
}

func TestEventNotifier_SaveBlockErrWhilePosting(t *testing.T) {
	expectedErr := "push failed"

	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient: &mocks.HttpClientStub{
			PostCalled: func(route string, payload interface{}, resp interface{}) error {
				require.True(t, route == pushEventEndpoint)
				return errors.New(expectedErr)
			},
		},
		PubKeyConverter: pubkeyConv(),
	})

	err := en.SaveBlock(test.SaveBlockArgsMock())
	require.True(t, strings.Contains(err.Error(), expectedErr))
}

func TestEventNotifier_RevertIndexedBlockNilMarshalizer(t *testing.T) {
	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient:      &mocks.HttpClientStub{},
		PubKeyConverter: pubkeyConv(),
		Marshalizer:     nil,
		Hasher:          &mock.HasherMock{},
	})

	err := en.RevertIndexedBlock(test.HeaderHandler(), nil)
	require.True(t, strings.Contains(err.Error(), core.ErrNilMarshalizer.Error()))
}

func TestEventNotifier_RevertIndexedBlockNilHasher(t *testing.T) {
	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient:      &mocks.HttpClientStub{},
		PubKeyConverter: pubkeyConv(),
		Marshalizer:     &mock.MarshalizerMock{},
		Hasher:          nil,
	})

	err := en.RevertIndexedBlock(test.HeaderHandler(), nil)
	require.True(t, strings.Contains(err.Error(), core.ErrNilHasher.Error()))
}

func TestEventNotifier_RevertIndexedBlockShouldWork(t *testing.T) {
	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient: &mocks.HttpClientStub{
			PostCalled: func(route string, payload interface{}, resp interface{}) error {
				require.True(t, route == revertEventsEndpoint)
				return nil
			},
		},
		PubKeyConverter: pubkeyConv(),
		Marshalizer:     &mock.MarshalizerMock{},
		Hasher:          &mock.HasherMock{},
	})

	err := en.RevertIndexedBlock(test.HeaderHandler(), nil)
	require.Nil(t, err)
}

func TestEventNotifier_RevertIndexedBlockErrWhilePosting(t *testing.T) {
	expectedErr := "revert fail"

	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient: &mocks.HttpClientStub{
			PostCalled: func(route string, payload interface{}, resp interface{}) error {
				require.True(t, route == revertEventsEndpoint)
				return errors.New(expectedErr)
			},
		},
		PubKeyConverter: pubkeyConv(),
		Marshalizer:     &mock.MarshalizerMock{},
		Hasher:          &mock.HasherMock{},
	})

	err := en.RevertIndexedBlock(test.HeaderHandler(), nil)
	require.True(t, strings.Contains(err.Error(), expectedErr))
}

func TestEventNotifier_FinalizedBlockShouldWork(t *testing.T) {
	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient: &mocks.HttpClientStub{
			PostCalled: func(route string, payload interface{}, resp interface{}) error {
				require.True(t, route == finalizedEventsEndpoint)
				return nil
			},
		},
		PubKeyConverter: pubkeyConv(),
		Marshalizer:     &mock.MarshalizerMock{},
		Hasher:          &mock.HasherMock{},
	})

	err := en.FinalizedBlock([]byte(test.RandStr(32)))
	require.Nil(t, err)
}

func TestEventNotifier_FinalizedBlockErrWhilePosting(t *testing.T) {
	expectedErr := "finalize fail"

	en, _ := NewEventNotifier(EventNotifierArgs{
		HttpClient: &mocks.HttpClientStub{
			PostCalled: func(route string, payload interface{}, resp interface{}) error {
				require.True(t, route == finalizedEventsEndpoint)
				return errors.New(expectedErr)
			},
		},
		PubKeyConverter: pubkeyConv(),
		Marshalizer:     &mock.MarshalizerMock{},
		Hasher:          &mock.HasherMock{},
	})

	err := en.FinalizedBlock([]byte(test.RandStr(32)))
	require.True(t, strings.Contains(err.Error(), expectedErr))
}
