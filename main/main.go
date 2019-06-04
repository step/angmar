package main

import (
	"archive/tar"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/step/durin/tarutils"
)

type GithubAPI struct {
	Client *http.Client
}

func (api *GithubAPI) FetchTarball(url string, extractor tarutils.Extractor) {
	resp, _ := api.Client.Get(url)
	tarutils.Untar(resp.Body, extractor)
}

type Foo struct{}

func (_ Foo) ExtractFile(header tar.Header, reader io.Reader) error {
	file, ferr := os.OpenFile("/tmp/target/"+header.Name, os.O_CREATE|os.O_RDWR, header.FileInfo().Mode())
	fmt.Println("file error", ferr)
	_, err := io.Copy(file, reader)
	fmt.Println(file, err)
	file.Close()
	return nil
}

func (_ Foo) ExtractDir(header tar.Header, reader io.Reader) error {
	os.MkdirAll("/tmp/target/"+header.Name, header.FileInfo().Mode())
	return nil
}

func main() {
	extractor := Foo{}
	api := GithubAPI{http.DefaultClient}

	api.FetchTarball("http://localhost:8003/foo.tar.gz", extractor)
}
