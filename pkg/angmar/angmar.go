package angmar

import (
	"fmt"

	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/queueclient"
	"github.com/step/angmar/pkg/tarutils"
)

type Angmar struct {
	QueueClient queueclient.QueueClient
	Generator   tarutils.ExtractorGenerator
	ApiClient   gh.GithubAPI
}

func (a Angmar) Start(qName string, r chan<- bool, stop <-chan bool) {
	shouldStop := false
	go func() {
		shouldStop = <-stop
	}()

	for {
		if shouldStop {
			break
		}
		val, err := a.QueueClient.Dequeue(qName)
		if err != nil {
			continue
		}
		extractor := a.Generator.Generate("")
		err = a.ApiClient.FetchTarball(val, extractor)
		fmt.Println(err)
		r <- true
	}
}
