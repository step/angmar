package angmar_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/step/angmar/pkg/queueclient"
	"github.com/step/saurontypes"

	a "github.com/step/angmar/pkg/angmar"
	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/tarutils"
	"github.com/step/angmar/pkg/testutils"
)

type DefaultExtractorGenerator struct {
	mapFiles *testutils.MapFiles
}

func (d *DefaultExtractorGenerator) Generate(args ...string) tarutils.Extractor {
	basePath := filepath.Join(args[0], args[1])
	d.mapFiles = testutils.CreateMapFiles(map[string]string{}, []string{}, basePath)
	return d.mapFiles
}
func (d *DefaultExtractorGenerator) String() string {
	return ""
}

func TestAngmar(t *testing.T) {
	queueClient := queueclient.NewDefaultClient()
	generator := DefaultExtractorGenerator{}

	server, archiveServer := testutils.CreateServer()
	defer archiveServer.Close()
	apiClient := gh.GithubAPI{Client: server.Client()}

	logger := a.AngmarLogger{Logger: log.New(ioutil.Discard, "", log.LstdFlags)}
	angmar := a.NewAngmar(queueClient, &generator, apiClient, logger, 1)
	responseCh := make(chan bool)
	stopCh := make(chan bool)

	message := saurontypes.AngmarMessage{Url: server.URL, SHA: "0abcdef1234", Pusher: "me", Tasks: []string{"test", "lint"}}
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
	}, []string{"dir/"}, "me/0abcdef1234")

	if !reflect.DeepEqual(generator.mapFiles, expected) {
		t.Errorf("Untar failed: Wanted %s Got %s", expected, generator.mapFiles)
	}

	for _, q := range []string{"test", "lint"} {
		val, err := queueClient.Dequeue(q)
		if err != nil {
			t.Errorf("Unexpected error while dequeuing from test")
		}

		if val != "me/0abcdef1234" {
			t.Errorf("Expected %s, got %s while testing downstream queue", "me/0abcdef1234", val)
		}
	}
}
