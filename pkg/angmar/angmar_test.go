package angmar_test

import (
	"fmt"
	"reflect"
	"testing"

	a "github.com/step/angmar/pkg/angmar"
	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/tarutils"
	"github.com/step/angmar/pkg/testutils"
)

type inMemoryQueue struct {
	queues map[string][]string
}

func (q inMemoryQueue) Enqueue(name, value string) error {
	q.queues[name] = append(q.queues[name], value)
	return nil
}

func (q inMemoryQueue) Dequeue(name string) (string, error) {
	queue := q.queues[name]
	if len(queue) == 0 {
		return "", fmt.Errorf("queue length is 0")
	}
	value := queue[len(queue)-1]
	q.queues[name] = queue[:len(queue)-1]
	return value, nil
}

type DefaultExtractorGenerator struct {
	mapFiles *testutils.MapFiles
}

func (d *DefaultExtractorGenerator) Generate(args ...string) tarutils.Extractor {
	d.mapFiles = testutils.CreateMapFiles(map[string]string{}, []string{})
	return d.mapFiles
}

func TestAngmar(t *testing.T) {
	queues := make(map[string][]string)
	queueClient := inMemoryQueue{queues}
	generator := DefaultExtractorGenerator{}

	server, archiveServer := testutils.CreateServer()
	defer archiveServer.Close()
	apiClient := gh.GithubAPI{server.Client()}

	angmar := a.Angmar{queueClient, &generator, apiClient}

	queueClient.Enqueue("queue", server.URL)
	responseCh := make(chan bool)
	stopCh := make(chan bool)
	go func() {
		angmar.Start("queue", responseCh, stopCh)
	}()

	result := <-responseCh
	if result != true {
		t.Errorf("an unexpected error occurred")
	}

	expected := testutils.CreateMapFiles(map[string]string{
		"dir/foo": "hello",
	}, []string{"dir/"})

	if !reflect.DeepEqual(generator.mapFiles, expected) {
		t.Errorf("Untar failed: Wanted %s Got %s", expected, generator.mapFiles)
	}
}
