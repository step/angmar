package main

import (
	"fmt"
	"net/http"

	"github.com/step/angmar/tarutils"
)

type GithubAPI struct {
	Client *http.Client
}

func (api *GithubAPI) FetchTarball(url string, extractor tarutils.Extractor) {
	resp, err := api.Client.Get(url)
	if err != nil {
		fmt.Println("Unable to fetch", url)
		return
	}
	err = tarutils.Untar(resp.Body, extractor)
	if err != nil {
		fmt.Println("Untar unsuccesful\n", err)
	}
}

func main() {
	extractor := tarutils.NewDefaultExtractor("/tmp/angmar")
	api := GithubAPI{http.DefaultClient}

	api.FetchTarball("http://localhost:8003/boo.tar.gz", extractor)
}
