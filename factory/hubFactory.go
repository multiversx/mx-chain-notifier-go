package factory

import (
	"errors"
	"strings"

	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/disabled"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/filters"
)

// TODO: refactor this after dispacher refactoring analysis

const separator = ":"

// ErrMakeCustomHub signals that creation of a new custom hub failed
var ErrMakeCustomHub = errors.New("failed to make custom hub")

var availableHubDelegates = map[string]func() dispatcher.Hub{
	"default-custom": func() dispatcher.Hub {
		return nil
	},
}

// CreateHub creates a common hub component
func CreateHub(apiType string, hubType string) (dispatcher.Hub, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return &disabled.Hub{}, nil
	case common.WSAPIType:
		return makeHub(hubType)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func makeHub(hubType string) (dispatcher.Hub, error) {
	if hubType == common.CommonHubType {
		eventFilter := filters.NewDefaultFilter()
		commonHub := hub.NewCommonHub(eventFilter)
		return commonHub, nil
	}

	hubConfig := strings.Split(hubType, separator)
	return tryMakeCustomHubForID(hubConfig[1])
}

func tryMakeCustomHubForID(id string) (dispatcher.Hub, error) {
	if makeHubFunc, ok := availableHubDelegates[id]; ok {
		customHub := makeHubFunc()
		return customHub, nil
	}

	return nil, ErrMakeCustomHub
}
