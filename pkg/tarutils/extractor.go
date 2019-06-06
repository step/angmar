package tarutils

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// Extractor is an interface that can be used whenever one
// wants to extract a file or a directory based on a tar header
// It is left to the caller that calls the interface to decide
// if the header is that of a file or a directory
type Extractor interface {
	ExtractFile(tar.Header, io.Reader) error
	ExtractDir(tar.Header, io.Reader) error
}

// DefaultExtractor is a struct implementing the Extractor interface.
// It has a src field that indicates the base path where one wants to
// extract files and directories to.
type DefaultExtractor struct {
	src string
}

// DefaultExtractor.ExtractFile extracts the given file under src specified
// in DefaultExtractor
func (extractor DefaultExtractor) ExtractFile(header tar.Header, reader io.Reader) (rerr error) {
	// Open file and defer file.Close()
	location := "DefaultExtractor.ExtractFile"
	fileName := filepath.Join(extractor.src, header.Name)
	file, ferr := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, header.FileInfo().Mode())

	// Handle defer in an anonymous func
	defer func() {
		err := file.Close()
		if err != nil {
			rerr = FileCloseError{fileName, err, rerr, location}
		}
	}()

	// Unable to open file
	if ferr != nil {
		return FileOpenError{fileName, header.Name, header.Mode, ferr, location}
	}

	// Copy file
	numBytesCopied, err := io.Copy(file, reader)

	// Unable to copy file
	if err != nil {
		return FileCopyError{header.Name, fileName, numBytesCopied, err, location}
	}

	return nil
}

// DefaultExtractor.ExtractDir extracts the given dir under src specified
// in DefaultExtractor
func (extractor DefaultExtractor) ExtractDir(header tar.Header, reader io.Reader) (rerr error) {
	// Create directory in src
	dirName := filepath.Join(extractor.src, header.Name)
	derr := os.MkdirAll(dirName, header.FileInfo().Mode())
	if derr != nil {
		return MakeDirError{dirName, derr, "DefaultExtractor.ExtractDir"}
	}

	return nil
}

// NewDefaultExtractor creates an instance of DefaultExtractor with the specified src
func NewDefaultExtractor(src string) DefaultExtractor {
	return DefaultExtractor{src}
}
