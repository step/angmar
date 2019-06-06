package gh

import (
	"fmt"
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
	resp, err := api.Client.Get(url)

	if err != nil {
		return fmt.Errorf("Unable to fetch %s\n%s", url, err.Error())
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Non 2xx response for %s", url)
	}

	if err = tarutils.Untar(resp.Body, extractor); err != nil {
		return fmt.Errorf("Unable to untar while fetching %s\n%s", url, err.Error())
	}

	return nil
}
