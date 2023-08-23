package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	apiErrors "github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/assert"
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

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		facade := &mocks.FacadeStub{
			ServeCalled: func(w http.ResponseWriter, r *http.Request) {
				wasCalled = true
			},
		}

		hg, err := groups.NewHubGroup(facade)
		require.NoError(t, err)
		require.NotNil(t, hg)

		ws := startWebServer(hg, hubPath, getHubRoutesConfig())

		req, _ := http.NewRequest("GET", "/hub/ws", nil)
		resp := httptest.NewRecorder()

		ws.ServeHTTP(resp, req)

		require.Equal(t, 0, len(hg.GetAdditionalMiddlewares()))
		assert.True(t, wasCalled)
	})
}

func getHubRoutesConfig() config.APIRoutesConfig {
	return config.APIRoutesConfig{
		APIPackages: map[string]config.APIPackageConfig{
			"hub": {
				Routes: []config.RouteConfig{
					{Name: "/ws", Open: true},
				},
			},
		},
	}
}
