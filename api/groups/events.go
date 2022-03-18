package groups

import (
	"context"
	"net/http"
	"time"

	"github.com/ElrondNetwork/notifier-go/api/shared"
	"github.com/ElrondNetwork/notifier-go/config"
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/gin-gonic/gin"
)

const (
	baseEventsEndpoint      = "/events"
	pushEventsEndpoint      = "/push"
	revertEventsEndpoint    = "/revert"
	finalizedEventsEndpoint = "/finalized"

	setnxRetryMs = 500

	revertKeyPrefix    = "revert_"
	finalizedKeyPrefix = "finalized_"
)

type eventsHandler struct {
	*baseGroup
	notifierHub           dispatcher.Hub
	config                config.ConnectorApiConfig
	additionalMiddlewares []gin.HandlerFunc
	lockService           LockService
}

// NewEventsHandler registers handlers for the /events group
func NewEventsHandler(
	notifierHub dispatcher.Hub,
	config config.ConnectorApiConfig,
	lockService LockService,
) (*eventsHandler, error) {
	h := &eventsHandler{
		baseGroup:             &baseGroup{},
		notifierHub:           notifierHub,
		config:                config,
		additionalMiddlewares: make([]gin.HandlerFunc, 0),
		lockService:           lockService,
	}

	endpoints := []*shared.EndpointHandlerData{
		{
			Method:  http.MethodPost,
			Path:    pushEventsEndpoint,
			Handler: h.pushEvents,
		},
		{
			Method:  http.MethodPost,
			Path:    revertEventsEndpoint,
			Handler: h.revertEvents,
		},
		{
			Method:  http.MethodPost,
			Path:    finalizedEventsEndpoint,
			Handler: h.finalizedEvents,
		},
	}

	h.endpoints = endpoints

	return h, nil
}

// GetAdditionalMiddlewares return additional middlewares for this group
func (h *eventsHandler) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return h.additionalMiddlewares
}

func (h *eventsHandler) pushEvents(c *gin.Context) {
	var blockEvents data.BlockEvents

	err := c.Bind(&blockEvents)
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shouldProcessEvents := true
	if h.config.CheckDuplicates {
		shouldProcessEvents = h.tryCheckProcessedOrRetry(blockEvents.Hash)
	}

	if blockEvents.Events != nil && shouldProcessEvents {
		log.Info("received events for block",
			"block hash", blockEvents.Hash,
			"shouldProcess", shouldProcessEvents,
		)
		h.notifierHub.BroadcastChan() <- blockEvents
	} else {
		log.Info("received duplicated events for block",
			"block hash", blockEvents.Hash,
			"ignoring", true,
		)
	}

	shared.JSONResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) revertEvents(c *gin.Context) {
	var revertBlock data.RevertBlock

	err := c.Bind(&revertBlock)
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shouldProcessRevert := true
	if h.config.CheckDuplicates {
		revertKey := revertKeyPrefix + revertBlock.Hash
		shouldProcessRevert = h.tryCheckProcessedOrRetry(revertKey)
	}

	if shouldProcessRevert {
		log.Info("received revert event for block",
			"block hash", revertBlock.Hash,
			"shouldProcess", shouldProcessRevert,
		)
		h.notifierHub.BroadcastRevertChan() <- revertBlock
	} else {
		log.Info("received duplicated revert event for block",
			"block hash", revertBlock.Hash,
			"ignoring", true,
		)
	}

	shared.JSONResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) finalizedEvents(c *gin.Context) {
	var finalizedBlock data.FinalizedBlock

	err := c.Bind(&finalizedBlock)
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shouldProcessFinalized := true
	if h.config.CheckDuplicates {
		finalizedKey := finalizedKeyPrefix + finalizedBlock.Hash
		shouldProcessFinalized = h.tryCheckProcessedOrRetry(finalizedKey)
	}

	if shouldProcessFinalized {
		log.Info("received finalized events for block",
			"block hash", finalizedBlock.Hash,
			"shouldProcess", shouldProcessFinalized,
		)
		h.notifierHub.BroadcastFinalizedChan() <- finalizedBlock
	} else {
		log.Info("received duplicated finalized event for block",
			"block hash", finalizedBlock.Hash,
			"ignoring", true,
		)
	}

	shared.JSONResponse(c, http.StatusOK, nil, "")
}

func (h *eventsHandler) createMiddlewares() []gin.HandlerFunc {
	var middleware []gin.HandlerFunc

	if h.config.Username != "" && h.config.Password != "" {
		basicAuth := gin.BasicAuth(gin.Accounts{
			h.config.Username: h.config.Password,
		})
		middleware = append(middleware, basicAuth)
	}

	return middleware
}

func (h *eventsHandler) tryCheckProcessedOrRetry(blockHash string) bool {
	var err error
	var setSuccessful bool

	for {
		setSuccessful, err = h.lockService.IsBlockProcessed(context.Background(), blockHash)

		if err != nil {
			if !h.lockService.HasConnection(context.Background()) {
				log.Error("failure connecting to redis")
				time.Sleep(time.Millisecond * setnxRetryMs)
				continue
			}
		}

		break
	}

	return err == nil && setSuccessful
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *eventsHandler) IsInterfaceNil() bool {
	return h == nil
}
