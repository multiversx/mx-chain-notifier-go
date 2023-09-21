package factory

import (
	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	"github.com/multiversx/mx-chain-notifier-go/api/gin"
	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/multiversx/mx-chain-notifier-go/config"
)

func CreateWebServerHandler(facade shared.FacadeHandler, configs config.Configs) (shared.WebServerHandler, error) {
	marshaller, err := marshalFactory.NewMarshalizer(configs.MainConfig.General.InternalMarshaller.Type)
	if err != nil {
		return nil, err
	}

	payloadHandler, err := createPayloadHandler(marshaller, facade)
	if err != nil {
		return nil, err
	}

	webServerArgs := gin.ArgsWebServerHandler{
		Facade:         facade,
		PayloadHandler: payloadHandler,
		Configs:        configs,
	}

	return gin.NewWebServerHandler(webServerArgs)
}
