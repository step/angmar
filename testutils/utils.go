package testutils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
)

type MockFile struct {
	Name, Body string
	Mode       int64
}

func TarGzFiles(files []MockFile, dirs []string, writer io.Writer) {
	gzipWriter := gzip.NewWriter(writer)
	defer gzipWriter.Close()

	tw := tar.NewWriter(gzipWriter)
	defer tw.Close()

	for _, dir := range dirs {
		tw.WriteHeader(&tar.Header{
			Name: dir,
			Mode: 0777,
		})
	}

	for _, file := range files {
		tw.WriteHeader(&tar.Header{
			Name: file.Name,
			Mode: file.Mode,
			Size: int64(len(file.Body)),
		})
		tw.Write([]byte(file.Body))
	}
}

type MapFiles struct {
	files map[string]string
	dirs  []string
}

func (mapFiles *MapFiles) ExtractFile(header tar.Header, reader io.Reader) error {
	content, _ := ioutil.ReadAll(reader)
	mapFiles.files[header.Name] = string(content)
	return nil
}

func (mapFiles *MapFiles) ExtractDir(header tar.Header, reader io.Reader) error {
	mapFiles.dirs = append(mapFiles.dirs, header.Name)
	return nil
}

func NewMapFiles() MapFiles {
	return MapFiles{map[string]string{}, []string{}}
}

func CreateMapFiles(filesAndContents map[string]string, dirs []string) MapFiles {
	return MapFiles{filesAndContents, dirs}
}
