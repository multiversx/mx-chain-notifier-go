package groups_test

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/multiversx/mx-chain-notifier-go/api/shared"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func startWebServer(group shared.GroupHandler, path string) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	routes := ws.Group(path)
	group.RegisterRoutes(routes)
	return ws
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	logError(err)
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
