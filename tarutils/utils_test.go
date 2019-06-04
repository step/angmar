package tarutils

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/step/durin/testutils"
)

func TestUntar(t *testing.T) {
	var buffer bytes.Buffer

	var files = []testutils.MockFile{
		{"dir/foo", "hello"},
	}
	var dirs = []string{"dir/"}
	testutils.ZipFiles(files, dirs, &buffer)

	mapFiles := testutils.NewMapFiles()
	Untar(&buffer, &mapFiles)

	expected := testutils.CreateMapFiles(map[string]string{
		"dir/foo": "hello",
	}, []string{"dir/"})

	if !reflect.DeepEqual(mapFiles, expected) {
		t.Errorf("Untar failed: Wanted %s Got %s", expected, mapFiles)
	}
}
