package tarutils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

// Untar expects a reader that provides a gzipped and tarred
// stream and an extractor. It reads the given stream and
// calls the extractor's ExtractFile or ExtractDir based
// on whether the tar Header is of a file or a directory
func Untar(reader io.Reader, extractor Extractor) (rerr error) {
	// first unzip and defer close
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("Unable to create gzip reader\n%s", err.Error())
	}

	// trap the error in the defer
	defer func() {
		err := gzipReader.Close()
		if err != nil {
			rerr = err
		}
	}()

	tarReader := tar.NewReader(gzipReader)

	// For each header in the tar stream call the appropriate Extractor function
	for header, err := tarReader.Next(); err != io.EOF; header, err = tarReader.Next() {
		if err != nil {
			return fmt.Errorf("Error reading the tar header\n%s", err.Error())
		}
		extract := extractor.ExtractFile
		if header.FileInfo().IsDir() {
			extract = extractor.ExtractDir
		}
		if err := extract(*header, tarReader); err != nil {
			return fmt.Errorf("Error while extracting %s\n%s", header.Name, err.Error())
		}
	}

	return nil
}
