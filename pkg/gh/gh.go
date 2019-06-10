package gh

import (
	"net/http"

	"github.com/step/angmar/pkg/tarutils"
)

type HttpClient interface {
	Get(string) (*http.Response, error)
}

// GithubAPI is a struct that will have several methods associated with
// Github API calls. For ex: FetchTarball etc
// It takes an http.Client as its only field
type GithubAPI struct {
	Client HttpClient
}

// FetchTarball fetches an archive that is supposedly at the url provided
// It assumes that the tarball is gzipped and hands the response contents
// to tarutils.Untar and is extracted by the tarutils.Extractor provided
func (api *GithubAPI) FetchTarball(url string, extractor tarutils.Extractor) error {
	location := "FetchTarball"
	resp, err := api.Client.Get(url)

	if err != nil {
		return ClientFetchError{url, "GET", err, location}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return StatusCodeError{resp.StatusCode, url, "GET", location}
	}

	if err = tarutils.Untar(resp.Body, extractor); err != nil {
		return FetchUntarError{url, err, location}
	}

	return nil
}
