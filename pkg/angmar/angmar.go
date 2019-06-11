package angmar

import (
	"encoding/json"
	"strings"

	"github.com/step/angmar/pkg/downloadclient"
	"github.com/step/angmar/pkg/queueclient"
	"github.com/step/angmar/pkg/tarutils"
)

type Angmar struct {
	QueueClient    queueclient.QueueClient
	Generator      tarutils.ExtractorGenerator
	DownloadClient downloadclient.DownloadClient
	Logger         AngmarLogger
}

func (a Angmar) String() string {
	var builder strings.Builder
	builder.WriteString(a.QueueClient.String() + "\n")
	builder.WriteString(a.Generator.String())
	return builder.String()
}

func worker(id int, angmar Angmar, messages <-chan AngmarMessage, rChan chan<- bool) {
	// jobs is buffered, so range is a blocking call if there are no jobs
	for message := range messages {
		angmar.Logger.ReceivedMessage(id, message)
		extractor := angmar.Generator.Generate(message.Pusher, message.SHA, message.Url)
		err := angmar.DownloadClient.Download(message.Url, extractor)

		if err != nil {
			angmar.Logger.LogError(id, err, message)
			rChan <- false
			continue
		}

		for _, q := range message.Tasks {
			err := angmar.QueueClient.Enqueue(q, extractor.GetBasePath())
			if err != nil {
				angmar.Logger.LogError(id, err, message)
				continue
			}
			angmar.Logger.TaskPlacedOnQueue(id, message, q)
		}
		rChan <- true
	}
}

func (a Angmar) Start(qName string, r chan<- bool, stop <-chan bool) {
	a.Logger.StartAngmar(a, qName)
	// A flag to stop placing jobs on worker threads
	shouldStop := false
	go func() {
		shouldStop = <-stop
	}()

	jobs := make(chan AngmarMessage, 10)

	// Create workers. The number 10 should eventually come from config
	// and be a field in the Angmar struct
	for index := 0; index < 10; index++ {
		go worker(index, a, jobs, r)
	}

	for {
		// Keep running till asked to stop
		if shouldStop {
			break
		}

		// Take the first job off the queue
		val, err := a.QueueClient.Dequeue(qName)
		if err != nil {
			continue
		}

		// Read the JSON
		var message AngmarMessage
		err = json.Unmarshal([]byte(val), &message)
		if err != nil {
			continue
		}

		// Schedule the job to be run by a worker.
		jobs <- message
	}
}
