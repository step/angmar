package angmar_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/step/angmar/pkg/queueclient"

	a "github.com/step/angmar/pkg/angmar"
	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/tarutils"
	"github.com/step/angmar/pkg/testutils"
)

type DefaultExtractorGenerator struct {
	mapFiles *testutils.MapFiles
}

func (d *DefaultExtractorGenerator) Generate(args ...string) tarutils.Extractor {
	d.mapFiles = testutils.CreateMapFiles(map[string]string{}, []string{})
	return d.mapFiles
}

func TestAngmar(t *testing.T) {
	queueClient := queueclient.NewDefaultClient()
	generator := DefaultExtractorGenerator{}

	server, archiveServer := testutils.CreateServer()
	defer archiveServer.Close()
	apiClient := gh.GithubAPI{Client: server.Client()}

	angmar := a.Angmar{QueueClient: queueClient, Generator: &generator, DownloadClient: apiClient}
	responseCh := make(chan bool)
	stopCh := make(chan bool)

	message := a.AngmarMessage{Url: server.URL, SHA: "0abcdef1234", Pusher: "me"}
	jsonMessage, _ := json.Marshal(message)

	if err := queueClient.Enqueue("queue", string(jsonMessage)); err != nil {
		t.Errorf("Unexpected error while queuing %s in memory", jsonMessage)
	}

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
