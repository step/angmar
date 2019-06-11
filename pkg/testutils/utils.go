package testutils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	files    map[string]string
	dirs     []string
	basePath string
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

func (mapFiles *MapFiles) GetBasePath() string {
	return mapFiles.basePath
}

func (mapFiles *MapFiles) String() string {
	return ""
}

func NewMapFiles() MapFiles {
	return MapFiles{map[string]string{}, []string{}, ""}
}

func CreateMapFiles(filesAndContents map[string]string, dirs []string, basePath string) *MapFiles {
	m := &MapFiles{filesAndContents, dirs, basePath}
	return m
}

func CreateServer() (*httptest.Server, *httptest.Server) {
	var buffer bytes.Buffer

	var files = []MockFile{
		{Name: "dir/foo", Body: "hello", Mode: 0777},
	}
	var dirs = []string{"dir/"}
	TarGzFiles(files, dirs, &buffer)

	archiveServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		rw.Write(buffer.Bytes())
	}))

	ghProxy := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
		if request.URL.String() == "/404" {
			rw.WriteHeader(404)
			if _, err := rw.Write([]byte("")); err != nil {
				fmt.Println("Unable to write empty byte!!!")
			}
			return
		}

		if request.URL.String() == "/archive" {
			rw.Header().Set("Location", archiveServer.URL)
			rw.WriteHeader(302)
		}

		if request.URL.String() == "/badtar" {
			if _, err := rw.Write([]byte("")); err != nil {
				fmt.Println("Unable to write empty byte!!!")
			}
			return
		}

		numOfBytesWritten, err := rw.Write(buffer.Bytes())
		if err != nil {
			fmt.Printf("Something went wrong while responding!\nWritten %d bytes\n%s", numOfBytesWritten, err.Error())
		}
	}))
	return ghProxy, archiveServer
}
