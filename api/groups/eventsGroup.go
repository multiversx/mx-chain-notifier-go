package groups

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-communication-go/websocket"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/common"
)

const (
	pushEventsEndpoint      = "/push"
	revertEventsEndpoint    = "/revert"
	finalizedEventsEndpoint = "/finalized"

	payloadVersionHeaderKey = "version"
)

// ArgsEventsGroup defines the arguments needed to create a new events group component
type ArgsEventsGroup struct {
	Facade         EventsFacadeHandler
	PayloadHandler websocket.PayloadHandler
}

type eventsGroup struct {
	*baseGroup
	facade         EventsFacadeHandler
	payloadHandler websocket.PayloadHandler
}

// NewEventsGroup registers handlers for the /events group
func NewEventsGroup(args ArgsEventsGroup) (*eventsGroup, error) {
	err := checkEventsGroupArgs(args)
	if err != nil {
		return nil, err
	}

	h := &eventsGroup{
		baseGroup:      newBaseGroup(),
		facade:         args.Facade,
		payloadHandler: args.PayloadHandler,
	}

	h.createMiddlewares()

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

func checkEventsGroupArgs(args ArgsEventsGroup) error {
	if check.IfNil(args.Facade) {
		return fmt.Errorf("%w for events group", errors.ErrNilFacadeHandler)
	}
	if check.IfNil(args.PayloadHandler) {
		return fmt.Errorf("%w for events group", errors.ErrNilPayloadHandler)
	}

	return nil
}

// TODO: remove this behaviour after deprecating http connector
// If received version not already handled, go to v0, for pre versioning setup
func getPayloadVersion(c *gin.Context) uint32 {
	version, err := strconv.Atoi(c.GetHeader(payloadVersionHeaderKey))
	if err != nil {
		log.Warn("failed to parse version header, used default version")
		return common.PayloadV0
	}

	return uint32(version)
}

func (h *eventsGroup) pushEvents(c *gin.Context) {
	pushEventsRawData, err := c.GetRawData()
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	payloadVersion := getPayloadVersion(c)

	err = h.payloadHandler.ProcessPayload(pushEventsRawData, outport.TopicSaveBlock, payloadVersion)
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shared.JSONResponse(c, http.StatusOK, nil, "")
}

func (h *eventsGroup) revertEvents(c *gin.Context) {
	revertEventsRawData, err := c.GetRawData()
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	payloadVersion := getPayloadVersion(c)

	err = h.payloadHandler.ProcessPayload(revertEventsRawData, outport.TopicRevertIndexedBlock, payloadVersion)
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shared.JSONResponse(c, http.StatusOK, nil, "")
}

func (h *eventsGroup) finalizedEvents(c *gin.Context) {
	finalizedRawData, err := c.GetRawData()
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	payloadVersion := getPayloadVersion(c)

	err = h.payloadHandler.ProcessPayload(finalizedRawData, outport.TopicFinalizedBlock, payloadVersion)
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shared.JSONResponse(c, http.StatusOK, nil, "")
}

func (h *eventsGroup) createMiddlewares() {
	user, pass := h.facade.GetConnectorUserAndPass()

	if user != "" && pass != "" {
		basicAuth := gin.BasicAuth(gin.Accounts{
			user: pass,
		})
		h.authMiddleware = basicAuth
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *eventsGroup) IsInterfaceNil() bool {
	return h == nil
}
