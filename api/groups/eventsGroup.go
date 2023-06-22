package groups

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-communication-go/websocket"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-notifier-go/api/errors"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
)

const (
	pushEventsEndpoint      = "/push"
	revertEventsEndpoint    = "/revert"
	finalizedEventsEndpoint = "/finalized"
)

// ArgsEventsGroup defines the arguments needed to create a new events group component
type ArgsEventsGroup struct {
	Facade            EventsFacadeHandler
	PayloadHandler    websocket.PayloadHandler
	EventsDataHandler EventsDataHandler
}

type eventsGroup struct {
	*baseGroup
	facade                EventsFacadeHandler
	payloadHandler        websocket.PayloadHandler
	eventsDataHandler     EventsDataHandler
	additionalMiddlewares []gin.HandlerFunc
}

// NewEventsGroup registers handlers for the /events group
func NewEventsGroup(args ArgsEventsGroup) (*eventsGroup, error) {
	err := checkEventsGroupArgs(args)
	if err != nil {
		return nil, err
	}

	h := &eventsGroup{
		baseGroup:             &baseGroup{},
		facade:                args.Facade,
		payloadHandler:        args.PayloadHandler,
		eventsDataHandler:     args.EventsDataHandler,
		additionalMiddlewares: make([]gin.HandlerFunc, 0),
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
	if check.IfNil(args.EventsDataHandler) {
		return fmt.Errorf("%w for events group", ErrNilEventsDataHandler)
	}

	return nil
}

// GetAdditionalMiddlewares return additional middlewares for this group
func (h *eventsGroup) GetAdditionalMiddlewares() []gin.HandlerFunc {
	return h.additionalMiddlewares
}

func (h *eventsGroup) pushEvents(c *gin.Context) {
	pushEventsRawData, err := c.GetRawData()
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err = h.payloadHandler.ProcessPayload(pushEventsRawData, outport.TopicSaveBlock)
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

	err = h.payloadHandler.ProcessPayload(revertEventsRawData, outport.TopicRevertIndexedBlock)
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

	err = h.payloadHandler.ProcessPayload(finalizedRawData, outport.TopicFinalizedBlock)
	if err != nil {
		shared.JSONResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	shared.JSONResponse(c, http.StatusOK, nil, "")
}

func (h *eventsGroup) createMiddlewares() {
	var middleware []gin.HandlerFunc

	user, pass := h.facade.GetConnectorUserAndPass()

	if user != "" && pass != "" {
		basicAuth := gin.BasicAuth(gin.Accounts{
			user: pass,
		})
		middleware = append(middleware, basicAuth)
	}

	h.additionalMiddlewares = middleware
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *eventsGroup) IsInterfaceNil() bool {
	return h == nil
}
