package main

import (
	"fmt"

	"github.com/multiversx/mx-chain-notifier-go/testdata"
)

func main() {
	args := HTTPClientWrapperArgs{
		UseAuthorization:  false,
		BaseUrl:           "http://localhost:5000",
		RequestTimeoutSec: 10,
	}
	httpClient, err := NewHTTPWrapperClient(args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = httpClient.Post("/events/push", testdata.OutportBlockV0())
	if err != nil {
		fmt.Println(fmt.Errorf("%w in eventNotifier.SaveBlock while posting block data", err))
		return
	}

	err = httpClient.Post("/events/revert", testdata.RevertBlockV0())
	if err != nil {
		fmt.Println(fmt.Errorf("%w in eventNotifier.SaveBlock while posting block data", err))
		return
	}

	err = httpClient.Post("/events/finalized", testdata.FinalizedBlockV0())
	if err != nil {
		fmt.Println(fmt.Errorf("%w in eventNotifier.SaveBlock while posting block data", err))
		return
	}
}
