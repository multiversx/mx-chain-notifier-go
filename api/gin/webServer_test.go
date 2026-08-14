package gin_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-communication-go/testscommon"
	"github.com/multiversx/mx-chain-core-go/core/check"
	apiErrors "github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/gin"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/common"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/require"
)

func createMockArgsWebServerHandler() gin.ArgsWebServerHandler {
	return gin.ArgsWebServerHandler{
		Facade:         &mocks.FacadeStub{},
		PayloadHandler: &testscommon.PayloadHandlerStub{},
		Configs: config.Configs{
			MainConfig: config.MainConfig{
				ConnectorApi: config.ConnectorApiConfig{
					Host: "8080",
				},
			},
			Flags: config.FlagsConfig{
				PublisherType: "notifier",
			},
		},
	}
}

func TestNewWebServerHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Facade = nil

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, apiErrors.ErrNilFacadeHandler, err)
	})

	t.Run("nil payload handler", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.PayloadHandler = nil

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, apiErrors.ErrNilPayloadHandler, err)
	})

	t.Run("invalid api type", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Configs.Flags.PublisherType = ""

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, common.ErrInvalidAPIType, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()

		ws, err := gin.NewWebServerHandler(args)
		require.Nil(t, err)
		require.NotNil(t, ws)

		err = ws.Run()
		require.Nil(t, err)

		err = ws.Close()
		require.Nil(t, err)
	})

	t.Run("route requires auth but no credentials configured, should fail to start", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Configs.MainConfig.ConnectorApi.Enabled = true
		args.Configs.ApiRoutesConfig = config.APIRoutesConfig{
			APIPackages: map[string]config.APIPackageConfig{
				"events": {
					Routes: []config.RouteConfig{
						{Name: "/push", Open: true, Auth: true},
					},
				},
			},
		}

		ws, err := gin.NewWebServerHandler(args)
		require.Nil(t, err)
		require.NotNil(t, ws)

		err = ws.Run()
		require.True(t, errors.Is(err, groups.ErrAuthEnabledWithoutMiddleware))
	})

	t.Run("route requires auth and credentials are configured, should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Configs.MainConfig.ConnectorApi.Enabled = true
		args.Facade = &mocks.FacadeStub{
			GetConnectorUserAndPassCalled: func() (string, string) {
				return "user", "pass"
			},
		}
		args.Configs.ApiRoutesConfig = config.APIRoutesConfig{
			APIPackages: map[string]config.APIPackageConfig{
				"events": {
					Routes: []config.RouteConfig{
						{Name: "/push", Open: true, Auth: true},
					},
				},
			},
		}

		ws, err := gin.NewWebServerHandler(args)
		require.Nil(t, err)
		require.NotNil(t, ws)

		err = ws.Run()
		require.Nil(t, err)

		err = ws.Close()
		require.Nil(t, err)
	})
}
