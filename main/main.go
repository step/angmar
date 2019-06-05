package main

import (
	"net/http"

	"github.com/step/angmar/gh"
	"github.com/step/angmar/tarutils"
)

func main() {
	extractor := tarutils.NewDefaultExtractor("/tmp/angmar")
	api := gh.GithubAPI{http.DefaultClient}

	api.FetchTarball("http://localhost:8003/boo.tar.gz", extractor)
}
