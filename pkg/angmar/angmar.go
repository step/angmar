package angmar

import (
	"encoding/json"
	"fmt"

	"github.com/step/angmar/pkg/downloadclient"
	"github.com/step/angmar/pkg/queueclient"
	"github.com/step/angmar/pkg/tarutils"
)

type Angmar struct {
	QueueClient    queueclient.QueueClient
	Generator      tarutils.ExtractorGenerator
	DownloadClient downloadclient.DownloadClient
}

type AngmarMessage struct {
	Url    string
	SHA    string
	Pusher string
}

type Job struct {
	message        AngmarMessage
	generator      tarutils.ExtractorGenerator
	downloadClient downloadclient.DownloadClient
}

func worker(id int, jobs <-chan Job, rChan chan<- bool) {
	// jobs is buffered, so range is a blocking call if there are no jobs
	for job := range jobs {
		message := job.message
		extractor := job.generator.Generate(message.Pusher, message.SHA, message.Url)
		err := job.downloadClient.Download(message.Url, extractor)

		// Assume everything went well
		response := true

		if err != nil {
			fmt.Println(err)
			response = false
		}

		rChan <- response
	}
}

func (a Angmar) Start(qName string, r chan<- bool, stop <-chan bool) {
	// A flag to stop placing jobs on worker threads
	shouldStop := false
	go func() {
		shouldStop = <-stop
	}()

	jobs := make(chan Job, 10)

	// Create workers. The number 10 should eventually come from config
	// and be a field in the Angmar struct
	for index := 0; index < 10; index++ {
		go worker(index, jobs, r)
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
		jobs <- Job{message, a.Generator, a.DownloadClient}
	}
}
