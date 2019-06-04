package tarutils

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/step/durin/testutils"
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
		{Name: "dir/foo", Body: "hello"},
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
