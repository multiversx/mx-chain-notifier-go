package groups

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
)

const (
	metricsPath           = "/metrics"
	prometheusMetricsPath = "/prometheus-metrics"
)

type statusGroup struct {
	*baseGroup
	facade                shared.FacadeHandler
	mutFacade             sync.RWMutex
	additionalMiddlewares []gin.HandlerFunc
}

// NewStatusGroup returns a new instance of status group
func NewStatusGroup(facade shared.FacadeHandler) (*statusGroup, error) {
	if check.IfNil(facade) {
		return nil, fmt.Errorf("%w for status group", errors.ErrNilFacadeHandler)
	}

	sg := &statusGroup{
		facade:                facade,
		baseGroup:             &baseGroup{},
		additionalMiddlewares: make([]gin.HandlerFunc, 0),
	}

	endpoints := []*shared.EndpointHandlerData{
		{
			Path:    metricsPath,
			Handler: sg.getMetrics,
			Method:  http.MethodGet,
		},
		{
			Path:    prometheusMetricsPath,
			Handler: sg.getPrometheusMetrics,
			Method:  http.MethodGet,
		},
	}
	sg.endpoints = endpoints

	return sg, nil
}

// getMetrics will expose endpoints statistics in json format
func (sg *statusGroup) getMetrics(c *gin.Context) {
	metricsResults := sg.facade.GetMetrics()

	shared.JSONResponse(c, http.StatusOK, gin.H{"metrics": metricsResults}, "")
}

// getPrometheusMetrics will expose proxy metrics in prometheus format
func (sg *statusGroup) getPrometheusMetrics(c *gin.Context) {
	metricsResults := sg.facade.GetMetricsForPrometheus()

	c.String(http.StatusOK, metricsResults)
}

// GetAdditionalMiddlewares return additional middlewares for this group
func (sg *statusGroup) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return sg.additionalMiddlewares
}

// IsInterfaceNil returns true if there is no value under the interface
func (sg *statusGroup) IsInterfaceNil() bool {
	return sg == nil
}
