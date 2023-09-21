package preprocess_test

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/mock"
	"github.com/multiversx/mx-chain-notifier-go/mocks"
	"github.com/multiversx/mx-chain-notifier-go/process/preprocess"
)

func createMockEventsDataPreProcessorArgs() preprocess.ArgsEventsPreProcessor {
	return preprocess.ArgsEventsPreProcessor{
		Marshaller: &mock.MarshalizerMock{},
		Facade:     &mocks.FacadeStub{},
	}
}

func TestNewEventsDataPreProcessor(t *testing.T) {
	t.Parallel()

}
