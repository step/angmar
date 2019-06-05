package gh_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/step/angmar/gh"
	"github.com/step/angmar/testutils"
)

func createServer() *httptest.Server {
	var buffer bytes.Buffer

	var files = []testutils.MockFile{
		{Name: "dir/foo", Body: "hello", Mode: 0777},
	}
	var dirs = []string{"dir/"}
	testutils.TarGzFiles(files, dirs, &buffer)

	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		if request.URL.String() == "/404" {
			rw.WriteHeader(404)
			if _, err := rw.Write([]byte("")); err != nil {
				fmt.Println("Unable to write empty byte!!!")
			}
			return
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
}

func TestFetchTarball(t *testing.T) {
	server := createServer()
	defer server.Close()

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
	server := createServer()
	defer server.Close()

	api := gh.GithubAPI{Client: server.Client()}

	mapFiles := testutils.NewMapFiles()
	err := api.FetchTarball(server.URL+"/404", &mapFiles)

	if err == nil {
		t.Error("expecting error but got no error")
	}
}

func TestFetchTarballFailsWhileUntarring(t *testing.T) {
	server := createServer()
	defer server.Close()

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
	server := createServer()
	defer server.Close()

	api := gh.GithubAPI{Client: &ClientMock{}}

	mapFiles := testutils.NewMapFiles()
	err := api.FetchTarball(server.URL+"/badtar", &mapFiles)

	if err == nil {
		t.Error("expecting error but got no error")
	}
}
