package gh_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/testutils"
)

func createServer() (*httptest.Server, *httptest.Server) {
	return testutils.CreateServer()
}

func TestFetchTarball(t *testing.T) {
	server, archiveServer := createServer()
	defer server.Close()
	defer archiveServer.Close()

	api := gh.GithubAPI{Client: server.Client()}

	mapFiles := testutils.NewMapFiles()
	if err := api.FetchTarball(server.URL, &mapFiles); err != nil {
		t.Errorf("Unexpected error while fetching %s\n%s", server.URL, err.Error())
	}
	expected := testutils.CreateMapFiles(map[string]string{
		"dir/foo": "hello",
	}, []string{"dir/"}, "")

	if !reflect.DeepEqual(&mapFiles, expected) {
		t.Errorf("Wanted %s Got %s", expected, &mapFiles)
	}
}

func TestFetchTarballFailsWhileFetching(t *testing.T) {
	server, archiveServer := createServer()
	defer server.Close()
	defer archiveServer.Close()

	api := gh.GithubAPI{Client: server.Client()}

	mapFiles := testutils.NewMapFiles()
	err := api.FetchTarball(server.URL+"/404", &mapFiles)

	if err == nil {
		t.Error("expecting error but got no error")
	}
}

func TestFetchTarballFailsWhileUntarring(t *testing.T) {
	server, archiveServer := createServer()
	defer server.Close()
	defer archiveServer.Close()

	api := gh.GithubAPI{Client: server.Client()}

	mapFiles := testutils.NewMapFiles()
	err := api.FetchTarball(server.URL+"/badtar", &mapFiles)

	if err == nil {
		t.Error("expecting error but got no error")
	}

}

type ClientMock struct{}

func (c *ClientMock) Get(url string) (*http.Response, error) {
	return nil, fmt.Errorf("bad client")
}
func TestFetchTarballFailsWithBadClient(t *testing.T) {
	server, archiveServer := createServer()
	defer server.Close()
	defer archiveServer.Close()

	api := gh.GithubAPI{Client: &ClientMock{}}

	mapFiles := testutils.NewMapFiles()
	err := api.FetchTarball(server.URL+"/badtar", &mapFiles)

	if err == nil {
		t.Error("expecting error but got no error")
	}
}

func TestFetchTarballWithRedirect(t *testing.T) {
	server, archiveServer := createServer()
	defer server.Close()
	defer archiveServer.Close()

	client := server.Client()
	api := gh.GithubAPI{Client: client}

	mapFiles := testutils.NewMapFiles()
	if err := api.FetchTarball(server.URL+"/archive", &mapFiles); err != nil {
		t.Errorf("Unexpected error while fetching %s\n%s", server.URL, err.Error())
	}

	expected := testutils.CreateMapFiles(map[string]string{
		"dir/foo": "hello",
	}, []string{"dir/"}, "")

	if !reflect.DeepEqual(&mapFiles, expected) {
		t.Errorf("Wanted %s Got %s", expected, &mapFiles)
	}
}
