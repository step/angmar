package tarutils

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

// Extractor is an interface that can be used whenever one
// wants to extract a file or a directory based on a tar header
// It is left to the caller that calls the interface to decide
// if the header is that of a file or a directory
type Extractor interface {
	ExtractFile(tar.Header, io.Reader) error
	ExtractDir(tar.Header, io.Reader) error
}

// Untar expects a reader that provides a gzipped and tarred
// stream and an extractor. It reads the given stream and
// calls the extractor's ExtractFile or ExtractDir based
// on whether the tar Header is of a file or a directory
func Untar(reader io.Reader, extractor Extractor) {
	// first unzip and defer close
	gzipReader, _ := gzip.NewReader(reader)
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// For each header in the tar stream call the appropriate Extractor function
	for header, err := tarReader.Next(); err != io.EOF; header, err = tarReader.Next() {
		extract := extractor.ExtractFile
		if header.FileInfo().IsDir() {
			extract = extractor.ExtractDir
		}
		extract(*header, tarReader)
	}
}
