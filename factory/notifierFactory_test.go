package factory

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core"
	hashMock "github.com/ElrondNetwork/elrond-go-core/core/mock"
	"github.com/ElrondNetwork/elrond-go-core/data/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateEventNotifier(t *testing.T) {
	tests := []struct {
		args func() *EventNotifierFactoryArgs
		err  error
	}{
		{
			args: func() *EventNotifierFactoryArgs {
				return &EventNotifierFactoryArgs{
					Marshalizer: nil,
					Hasher:      &hashMock.HasherStub{},
				}
			},
			err: core.ErrNilMarshalizer,
		},
		{
			args: func() *EventNotifierFactoryArgs {
				return &EventNotifierFactoryArgs{
					Marshalizer: &mock.MarshalizerStub{},
					Hasher:      nil,
				}
			},
			err: core.ErrNilHasher,
		},
		{
			args: func() *EventNotifierFactoryArgs {
				return &EventNotifierFactoryArgs{
					Marshalizer: &mock.MarshalizerStub{},
					Hasher:      mock.HasherMock{},
				}
			},
			err: nil,
		},
	}

	for _, test := range tests {
		_, err := CreateEventNotifier(test.args())
		require.Equal(t, test.err, err)
	}
}
