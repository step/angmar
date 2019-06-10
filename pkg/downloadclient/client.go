package downloadclient

import "github.com/step/angmar/pkg/tarutils"

type DownloadClient interface {
	Download(url string, extractor tarutils.Extractor) error
}
