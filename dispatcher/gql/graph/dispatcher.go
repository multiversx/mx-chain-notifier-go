package graph

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/ElrondNetwork/notifier-go/dispatcher/gql/graph/model"
	"github.com/google/uuid"
)

type graphqlDispatcher struct {
	id           uuid.UUID
	dispatchChan chan<- []*model.Event
}

func NewGraphqlDispatcher() (*graphqlDispatcher, error) {
	return &graphqlDispatcher{
		id: uuid.New(),
	}, nil
}

func (gd *graphqlDispatcher) PushEvents(events []data.Event) {
	var gqlEvents []*model.Event
	for _, event := range events {
		var ptrTopics []*string
		for _, topic := range event.Topics {
			strTopic := string(topic)
			ptrTopics = append(ptrTopics, &strTopic)
		}

		strData := string(event.Data)
		gqlEvents = append(gqlEvents, &model.Event{
			Address:    &event.Address,
			Identifier: &event.Identifier,
			Data:       &strData,
			Topics:     ptrTopics,
		})
	}
	gd.dispatchChan <- gqlEvents
}

func (gd *graphqlDispatcher) GetID() uuid.UUID {
	return gd.id
}

func (gd *graphqlDispatcher) SetDispatchChan(dispatchChan chan<- []*model.Event) {
	gd.dispatchChan = dispatchChan
}
