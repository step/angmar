package angmar

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/step/angmar/pkg/downloadclient"
	"github.com/step/angmar/pkg/queueclient"
	"github.com/step/angmar/pkg/tarutils"
	"github.com/step/saurontypes"
)

type angmar struct {
	QueueClient      queueclient.QueueClient
	Generator        tarutils.ExtractorGenerator
	DownloadClient   downloadclient.DownloadClient
	Logger           AngmarLogger
	NumOfWorkers     int
	SourceMountPoint string
}

func (a angmar) String() string {
	var builder strings.Builder
	builder.WriteString(a.QueueClient.String() + "\n")
	builder.WriteString(a.Generator.String())
	return builder.String()
}

func worker(id int, a angmar, messages <-chan saurontypes.AngmarMessage, rChan chan<- bool) {
	// messages is buffered, so range is a blocking call if there are no messages
	for message := range messages {
		fmt.Println(id, message)
		a.Logger.ReceivedMessage(id, message)
		extractor := a.Generator.Generate(message.Project, message.Pusher, message.SHA)
		err := a.DownloadClient.Download(message.Url, extractor)

		if err != nil {
			a.Logger.LogError(id, err, message)
			rChan <- false
			continue
		}

		for _, q := range message.Tasks {
			extractorBasePath := extractor.GetBasePath()
			repoLocation := strings.Replace(extractorBasePath, a.SourceMountPoint+"/", "", 1)
			urukMessage := saurontypes.ConvertAngmarToUrukMessage(message, repoLocation)
			urukMessageAsJson, err := json.Marshal(urukMessage)
			if err != nil {
				return
			}
			err = a.QueueClient.Enqueue(q, string(urukMessageAsJson))
			if err != nil {
				a.Logger.LogError(id, err, message)
				continue
			}
			a.Logger.TaskPlacedOnQueue(id, message, q)
		}
		rChan <- true
	}
}

func (a angmar) Start(qName string, r chan<- bool, stop <-chan bool) {
	a.Logger.StartAngmar(a, qName)
	// A flag to stop placing jobs on worker threads
	shouldStop := false
	go func() {
		shouldStop = <-stop
	}()

	jobs := make(chan saurontypes.AngmarMessage, 10)

	// Create workers. The number 10 should eventually come from config
	// and be a field in the Angmar struct
	for index := 0; index < a.NumOfWorkers; index++ {
		go worker(index, a, jobs, r)
	}

	var val string
	var err error

	for {
		// Keep running till asked to stop
		if shouldStop {
			break
		}

		// Take the first job off the queue
		val, err = a.QueueClient.Dequeue(qName)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Read the JSON
		message := new(saurontypes.AngmarMessage)
		err = json.Unmarshal([]byte(val), message)
		if err != nil {
			continue
		}

		// Schedule the job to be run by a worker.
		jobs <- *message
	}
}

func NewAngmar(
	qClient queueclient.QueueClient,
	generator tarutils.ExtractorGenerator,
	dClient downloadclient.DownloadClient,
	logger AngmarLogger,
	numOfWorkers int,
	sourceMountPoint string) angmar {

	if numOfWorkers < 1 {
		numOfWorkers = 1
	}

	return angmar{qClient, generator, dClient, logger, numOfWorkers, sourceMountPoint}
}
