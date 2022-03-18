package factory

import (
	"errors"
	"strings"

	"github.com/ElrondNetwork/notifier-go/common"
	"github.com/ElrondNetwork/notifier-go/dispatcher"
	"github.com/ElrondNetwork/notifier-go/dispatcher/hub"
	"github.com/ElrondNetwork/notifier-go/filters"
)

// TODO: refactor this after dispacher refactoring analysis

const separator = ":"

// ErrMakeCustomHub -
var ErrMakeCustomHub = errors.New("failed to make custom hub")

var availableHubDelegates = map[string]func() dispatcher.Hub{
	"default-custom": func() dispatcher.Hub {
		return nil
	},
}

func CreateCommonHub(hubType common.HubType) (dispatcher.Hub, error) {
	commonHub, err := makeHub(hubType)
	if err != nil {
		return nil, err
	}

	return commonHub, nil
}

func makeHub(hubType common.HubType) (dispatcher.Hub, error) {
	if hubType == common.CommonHubType {
		eventFilter := filters.NewDefaultFilter()
		commonHub := hub.NewCommonHub(eventFilter)
		return commonHub, nil
	}

	hubConfig := strings.Split(string(hubType), separator)
	return tryMakeCustomHubForID(hubConfig[1])
}

func tryMakeCustomHubForID(id string) (dispatcher.Hub, error) {
	if makeHubFunc, ok := availableHubDelegates[id]; ok {
		customHub := makeHubFunc()
		return customHub, nil
	}

	return nil, ErrMakeCustomHub
}
