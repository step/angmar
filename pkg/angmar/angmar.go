package angmar

import (
	"encoding/json"

	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/queueclient"
	"github.com/step/angmar/pkg/tarutils"
)

type Angmar struct {
	QueueClient queueclient.QueueClient
	Generator   tarutils.ExtractorGenerator
	ApiClient   gh.GithubAPI
}

type AngmarMessage struct {
	Url    string
	SHA    string
	Pusher string
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

		var message AngmarMessage
		err = json.Unmarshal([]byte(val), &message)
		if err != nil {
			continue
		}
		extractor := a.Generator.Generate(message.Pusher, message.SHA)
		err = a.ApiClient.FetchTarball(message.Url, extractor)
		if err != nil {
			continue
		}
		r <- true
	}
}
