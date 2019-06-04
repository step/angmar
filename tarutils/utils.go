package tarutils

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

type Extractor interface {
	ExtractFile(tar.Header, io.Reader) error
	ExtractDir(tar.Header, io.Reader) error
}

func Untar(reader io.Reader, extractor Extractor) {
	gzipReader, _ := gzip.NewReader(reader)
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)

	for header, err := tarReader.Next(); err != io.EOF; header, err = tarReader.Next() {
		if header.FileInfo().IsDir() {
			extractor.ExtractDir(*header, tarReader)
		} else {
			extractor.ExtractFile(*header, tarReader)
		}
	}
}
