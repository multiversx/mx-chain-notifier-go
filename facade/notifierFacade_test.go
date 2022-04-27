package facade_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/facade"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockFacadeArgs() facade.ArgsNotifierFacade {
	return facade.ArgsNotifierFacade{
		EventsHandler: &mocks.EventsHandlerStub{},
		APIConfig:     config.ConnectorApiConfig{},
		WSHandler:     &mocks.WSHandlerStub{},
	}
}

func TestNewNotifierFacade(t *testing.T) {
	t.Parallel()

	t.Run("nil events handler", func(t *testing.T) {
		t.Parallel()

		args := createMockFacadeArgs()
		args.EventsHandler = nil

		f, err := facade.NewNotifierFacade(args)
		require.True(t, check.IfNil(f))
		require.Equal(t, facade.ErrNilEventsHandler, err)
	})

	t.Run("nil ws handler", func(t *testing.T) {
		t.Parallel()

		args := createMockFacadeArgs()
		args.WSHandler = nil

		f, err := facade.NewNotifierFacade(args)
		require.True(t, check.IfNil(f))
		require.Equal(t, facade.ErrNilWSHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockFacadeArgs()
		facade, err := facade.NewNotifierFacade(args)
		require.Nil(t, err)
		require.NotNil(t, facade)
	})
}

func TestGetters(t *testing.T) {
	t.Parallel()

	expuser := "user1"
	exppass := "pass1"

	args := createMockFacadeArgs()
	args.APIConfig.Username = expuser
	args.APIConfig.Password = exppass

	f, err := facade.NewNotifierFacade(args)
	require.Nil(t, err)

	user, pass := f.GetConnectorUserAndPass()
	assert.Equal(t, expuser, user)
	assert.Equal(t, exppass, pass)
}
