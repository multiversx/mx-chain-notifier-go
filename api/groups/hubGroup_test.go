package groups_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core/check"
	apiErrors "github.com/ElrondNetwork/notifier-go/api/errors"
	"github.com/ElrondNetwork/notifier-go/api/groups"
	"github.com/ElrondNetwork/notifier-go/mocks"
	"github.com/stretchr/testify/require"
)

const hubPath = "/hub"

func TestNewHubGroup(t *testing.T) {
	t.Parallel()

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		hg, err := groups.NewHubGroup(nil)
		require.True(t, errors.Is(err, apiErrors.ErrNilFacadeHandler))
		require.True(t, check.IfNil(hg))
	})

	t.Run("dispatch all, should work", func(t *testing.T) {
		t.Parallel()

		facade := &mocks.FacadeStub{
			GetDispatchTypeCalled: func() string {
				return "dispatch:*"
			},
		}

		hg, err := groups.NewHubGroup(facade)
		require.NoError(t, err)
		require.NotNil(t, hg)

		require.Equal(t, 0, len(hg.GetAdditionalMiddlewares()))
	})
}
