package gql

import (
	"github.com/ElrondNetwork/notifier-go/data"
	"github.com/google/uuid"
)

type graphqlDispatcher struct {
	id uuid.UUID
}

func NewGraphqlDispatcher() (*graphqlDispatcher, error) {
	return &graphqlDispatcher{}, nil
}

func (gd *graphqlDispatcher) PushEvents(events []data.Event) {

}

func (gd *graphqlDispatcher) GetID() uuid.UUID {
	return gd.id
}
