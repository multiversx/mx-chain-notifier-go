package gql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph/generated"
	"github.com/gorilla/websocket"
)

func NewGraphQLServer(hub dispatcher.Hub) *handler.Server {
	resolver := graph.NewResolver(hub)

	srv := handler.New(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: resolver,
		}),
	)
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.Use(extension.Introspection{})

	return srv
}
