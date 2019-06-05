package tarutils

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/step/angmar/testutils"
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
		{Name: "dir/foo", Body: "hello", Mode: 0777},
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

	if err == nil {
		t.Errorf("Not supposed to give a nil error! Half reader used!")
	}
}

func formattedTime() string {
	return time.Now().Format("02_01_06__03_04_05")
}

func TestDefaultExtractor(t *testing.T) {
	// Verify contents of file
	// Erase directory
	tmp := os.TempDir()
	timeSuffix := formattedTime()
	src := filepath.Join(tmp, "test"+timeSuffix)

	// Check if temp directory present
	_, err := os.Stat(src)
	if os.IsExist(err) {
		t.Errorf("Directory already present: %s", src)
	}

	// Tar the following files and dirs
	var buffer bytes.Buffer

	files := []testutils.MockFile{
		{Name: "dir/foo", Body: "hello", Mode: 0777},
	}
	dirs := []string{"dir/"}

	testutils.TarGzFiles(files, dirs, &buffer)

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

	err = os.RemoveAll(src)
	if err != nil {
		t.Errorf("An unexpected error occurred removing %s\n%s", src, err.Error())
	}
}
