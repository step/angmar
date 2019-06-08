package angmar_test

import (
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
