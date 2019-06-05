package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/step/angmar/testutils"
)

func createServer() *httptest.Server {
	var buffer bytes.Buffer

	var files = []testutils.MockFile{
		{"dir/foo", "hello", 0777},
	}
	var dirs = []string{"dir/"}
	testutils.TarGzFiles(files, dirs, &buffer)

	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		rw.Write(buffer.Bytes())
	}))
}

func TestFetchTarball(t *testing.T) {
	server := createServer()
	defer server.Close()

	api := GithubAPI{server.Client()}

	mapFiles := testutils.NewMapFiles()
	api.FetchTarball(server.URL, &mapFiles)

	expected := testutils.CreateMapFiles(map[string]string{
		"dir/foo": "hello",
	}, []string{"dir/"})

	if !reflect.DeepEqual(mapFiles, expected) {
		t.Errorf("Wanted %s Got %s", expected, mapFiles)
	}
}
