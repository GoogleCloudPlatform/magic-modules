package google

import (
	"testing"
	"testing/fstest"
)

func TestSimple(t *testing.T) {
	o := overlayFS{
		overlay: fstest.MapFS{
			"a/b": {Data: []byte("foo")},
			"c":   {Data: []byte("c")},
		},
		base: fstest.MapFS{
			"a/b": {Data: []byte("bar")},
			"d":   {Data: []byte("d")},
			"a/d": {Data: []byte("ad")},
		}}
	// TestFS checks for existence of files, but on top of it
	// runs fairly extensive validation tests on FS implementation.
	if err := fstest.TestFS(o, "a", "a/b", "c", "d", "a/d"); err != nil {
		t.Error(err)
	}
	contents, err := o.ReadFile("a/b")
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "foo" {
		t.Errorf("Unexpected contents for a/b, wanted 'foo', got %q", contents)
	}
}
