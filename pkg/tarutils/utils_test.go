package tarutils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/step/angmar/pkg/testutils"
)

func testUntarOfFiles(files []testutils.MockFile, dirs []string, expected testutils.MapFiles) func(t *testing.T) {
	return func(t *testing.T) {
		var buffer bytes.Buffer

		testutils.TarGzFiles(files, dirs, &buffer)

		mapFiles := testutils.NewMapFiles()
		err := Untar(&buffer, &mapFiles)

		if err != nil {
			t.Errorf("Unexpected error: %s\n", err.Error())
		}

		if !reflect.DeepEqual(mapFiles, expected) {
			t.Errorf("Untar failed: Wanted %s Got %s", expected, mapFiles)
		}
	}
}

func TestUntar(t *testing.T) {
	files := []testutils.MockFile{
		{Name: "dir/foo", Body: "hello", Mode: 0755},
	}
	dirs := []string{"dir/"}

	expected := testutils.CreateMapFiles(map[string]string{
		"dir/foo": "hello",
	}, []string{"dir/"})

	t.Run("Single file in single directory", testUntarOfFiles(files, dirs, expected))
}

func TestBadTar(t *testing.T) {
	var buffer bytes.Buffer

	mapFiles := testutils.NewMapFiles()
	err := Untar(&buffer, &mapFiles)

	if actualErr, ok := err.(GzipReaderCreateError); !ok {
		t.Errorf("Expected GzipReaderCreateError but got\n%s", actualErr)
	}
}

func formattedTime() string {
	return time.Now().Format("02_01_06__03_04_05")
}

func tmpSrcDir(prefix string) string {
	tmp := os.TempDir()
	timeSuffix := formattedTime()
	return filepath.Join(tmp, prefix+timeSuffix)
}

func createDefaultTarGz(writer io.Writer) {
	files := []testutils.MockFile{
		{Name: "dir/foo", Body: "hello", Mode: 0755},
		{Name: "pax_global_header", Body: "", Mode: 0000},
	}
	dirs := []string{"dir/"}

	testutils.TarGzFiles(files, dirs, writer)
}

func TestDefaultExtractor(t *testing.T) {
	// name of temporary timestamp based src dir
	src := tmpSrcDir("test")

	// Check if temp directory present
	_, err := os.Stat(src)
	if err == nil {
		t.Errorf("Directory already present: %s", src)
	}

	// Tar the following files and dirs
	var buffer bytes.Buffer
	createDefaultTarGz(&buffer)

	// Untar the files that are tarred into buffer
	extractor := NewDefaultExtractor(src)

	err = Untar(&buffer, extractor)
	if err != nil {
		t.Errorf("Unexpected error while Untarring: %s", err.Error())
	}

	// Verify contents of files written to disk
	untarredFileToTest := filepath.Join(src, "dir", "foo")
	file, err := os.Open(untarredFileToTest)
	if err != nil {
		t.Errorf("Unable to open file %s\n%s", untarredFileToTest, err.Error())
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Errorf("Unable to read file %s\n%s", file.Name(), err.Error())
	}

	if string(contents) != "hello" {
		t.Errorf("contents of file incorrect. Expected hello, got %s", contents)
	}

	// Verify no pax_global_header
	paxGlobalHeader := filepath.Join(src, "pax_global_header")
	if _, err := os.Stat(paxGlobalHeader); err == nil {
		t.Errorf("pax_global_header created!")
	}

	err = os.RemoveAll(src)
	if err != nil {
		t.Errorf("An unexpected error occurred removing %s\n%s", src, err.Error())
	}
}
