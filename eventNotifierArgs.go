package notifier

import (
	"github.com/ElrondNetwork/elrond-go/marshal"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
)

type EventNotifierArgs struct {
	Hub         dispatcher.Hub
	Marshalizer marshal.Marshalizer
}
