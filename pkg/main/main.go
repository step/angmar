package main

import (
	"fmt"
	"net/http"

	"github.com/step/angmar/pkg/gh"
	"github.com/step/angmar/pkg/tarutils"
)

func main() {
	extractor := tarutils.NewDefaultExtractor("/tmp/angmar")
	api := gh.GithubAPI{Client: http.DefaultClient}

	err := api.FetchTarball("http://localhost:8003/bad.tar.gz", extractor)
	fmt.Println(err)
}
