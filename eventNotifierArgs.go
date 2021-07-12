package notifier

import (
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/notifier-go/proxy/client"
)

type EventNotifierArgs struct {
	HttpClient  client.HttpClient
	Marshalizer marshal.Marshalizer
}
