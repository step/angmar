package gh_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/testutils"
)

func createServer() (*httptest.Server, *httptest.Server) {
	var buffer bytes.Buffer

	var files = []testutils.MockFile{
		{Name: "dir/foo", Body: "hello", Mode: 0777},
	}
	var dirs = []string{"dir/"}
	testutils.TarGzFiles(files, dirs, &buffer)

	archiveServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		rw.Write(buffer.Bytes())
	}))

	ghProxy := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		if request.URL.String() == "/404" {
			rw.WriteHeader(404)
			if _, err := rw.Write([]byte("")); err != nil {
				fmt.Println("Unable to write empty byte!!!")
			}
			return
		}

		if request.URL.String() == "/archive" {
			rw.Header().Set("Location", archiveServer.URL)
			rw.WriteHeader(302)
		}

		if request.URL.String() == "/badtar" {
			if _, err := rw.Write([]byte("")); err != nil {
				fmt.Println("Unable to write empty byte!!!")
			}
			return
		}

		numOfBytesWritten, err := rw.Write(buffer.Bytes())
		if err != nil {
			fmt.Printf("Something went wrong while responding!\nWritten %d bytes\n%s", numOfBytesWritten, err.Error())
		}
	}))
	return ghProxy, archiveServer
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
	}, []string{"dir/"})

	if !reflect.DeepEqual(mapFiles, expected) {
		t.Errorf("Wanted %s Got %s", expected, mapFiles)
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
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		fmt.Printf("Redirecting.......%s to %s", via[0].URL, req.URL)
		return nil
	}
	api := gh.GithubAPI{Client: client}

	mapFiles := testutils.NewMapFiles()
	if err := api.FetchTarball(server.URL+"/archive", &mapFiles); err != nil {
		t.Errorf("Unexpected error while fetching %s\n%s", server.URL, err.Error())
	}

	expected := testutils.CreateMapFiles(map[string]string{
		"dir/foo": "hello",
	}, []string{"dir/"})

	if !reflect.DeepEqual(mapFiles, expected) {
		t.Errorf("Wanted %s Got %s", expected, mapFiles)
	}
}
