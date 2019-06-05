package tarutils

import "fmt"

type GzipReaderCreateError struct {
	actualErr error
	location  string
}

func (g GzipReaderCreateError) Error() string {
	return fmt.Sprintf("Unable to create gzip reader at %s\n%s", g.location, g.actualErr.Error())
}

type GzipReaderCloseError struct {
	actualErr error
	location  string
}

func (g GzipReaderCloseError) Error() string {
	return fmt.Sprintf("Unable to close gzip reader at %s\n%s", g.location, g.actualErr.Error())
}

type TarHeaderError struct {
	actualErr error
	location  string
}

func (t TarHeaderError) Error() string {
	return fmt.Sprintf("Unable to read tar header %s\n%s", t.location, t.actualErr.Error())
}

type ExtractionError struct {
	name      string
	mode      int64
	actualErr error
	location  string
}

func (e ExtractionError) Error() string {
	return fmt.Sprintf("Unable to extract %s [%o] at %s\n%s", e.name, e.mode, e.location, e.actualErr.Error())
}

type FileOpenError struct {
	fileName  string
	name      string
	mode      int64
	actualErr error
	location  string
}

func (f FileOpenError) Error() string {
	return fmt.Sprintf("Unable to open %s(%s) [%o] at %s\n%s",
		f.fileName, f.name, f.mode, f.location, f.actualErr.Error())
}

type FileCopyError struct {
	src         string
	dest        string
	bytesCopied int64
	actualErr   error
	location    string
}

func (f FileCopyError) Error() string {
	return fmt.Sprintf("Unable to copy from %s to %s\nCopied %d bytes at %s\n%s",
		f.src, f.dest, f.bytesCopied, f.location, f.actualErr.Error())
}

type MakeDirError struct {
	dirName   string
	actualErr error
	location  string
}

func (d MakeDirError) Error() string {
	return fmt.Sprintf("Unable to make directory %s at %s\n%s", d.dirName, d.location, d.actualErr.Error())
}
