package groups_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	apiErrors "github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/groups"
	"github.com/multiversx/mx-chain-notifier-go/config"
	"github.com/multiversx/mx-chain-notifier-go/data"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const statusPath = "/status"

type statusMetricsResponse struct {
	Data struct {
		Metrics map[string]*data.EndpointMetricsResponse `json:"metrics"`
	}
	Error string `json:"error"`
}

func TestNewStatusGroup(t *testing.T) {
	t.Parallel()

	t.Run("nil facade should error", func(t *testing.T) {
		t.Parallel()

		sg, err := groups.NewStatusGroup(nil)

		require.True(t, errors.Is(err, apiErrors.ErrNilFacadeHandler))
		require.True(t, check.IfNil(sg))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		sg, err := groups.NewStatusGroup(&mocks.FacadeStub{})

		assert.NotNil(t, sg)
		assert.Nil(t, err)
	})
}

func TestGetMetrics_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedMetrics := map[string]*data.EndpointMetricsResponse{
		"/guardian/config": {
			NumRequests:       5,
			TotalResponseTime: 100,
		},
	}
	facade := &mocks.FacadeStub{
		GetMetricsCalled: func() map[string]*data.EndpointMetricsResponse {
			return expectedMetrics
		},
	}

	statusGroup, err := groups.NewStatusGroup(facade)
	require.Nil(t, err)

	ws := startWebServer(statusGroup, statusPath, getStatusRoutesConfig())

	req, _ := http.NewRequest("GET", "/status/metrics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	var apiResp statusMetricsResponse
	loadResponse(resp.Body, &apiResp)
	require.Equal(t, http.StatusOK, resp.Code)

	require.Equal(t, expectedMetrics, apiResp.Data.Metrics)
}

func TestGetPrometheusMetrics_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedMetrics := `num_requests{endpoint="/guardian/config"} 37`
	facade := &mocks.FacadeStub{
		GetMetricsForPrometheusCalled: func() string {
			return expectedMetrics
		},
	}

	statusGroup, err := groups.NewStatusGroup(facade)
	require.NoError(t, err)

	ws := startWebServer(statusGroup, statusPath, getStatusRoutesConfig())

	req, _ := http.NewRequest("GET", "/status/prometheus-metrics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, expectedMetrics, string(bodyBytes))
}

func TestStatusGroup_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	sg, _ := groups.NewStatusGroup(nil)
	assert.True(t, sg.IsInterfaceNil())

	sg, _ = groups.NewStatusGroup(&mocks.FacadeStub{})
	assert.False(t, sg.IsInterfaceNil())
}

func getStatusRoutesConfig() config.APIRoutesConfig {
	return config.APIRoutesConfig{
		APIPackages: map[string]config.APIPackageConfig{
			"status": {
				Routes: []config.RouteConfig{
					{Name: "/metrics", Open: true},
					{Name: "/prometheus-metrics", Open: true},
				},
			},
		},
	}
}
