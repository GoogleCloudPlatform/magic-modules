package google

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
)

// internal interface supporting ReadDirFS and ReadFileFS
//
// Most sensible FS implementations support these, as an example both `embed.FS`
// and `os.DirFS(...)` implement these.
type ReadDirReadFileFS interface {
	fs.ReadDirFS
	fs.ReadFileFS
}

// overlayFS is an fs.FS implementation which supports the concept of overlays.
//
// Given two fs.FS, `overlay` and `base`, NewOverlayFS(overlay, base) builds an FS
// which prioritizes files in `overlay` over files in `base`.
//
// As an example, given:
//   - overlay = {
//     "a/b": "foo",
//     "c": "c",
//     }
//   - base = {
//     "a/b": "bar",
//     "d": "d",
//     }
//
// Then:
//
//	ofs.ReadFile("a/b") -> "foo"
//	ofs.ReadFile("d") -> "d"
//	ofs.ReadDir(".") -> {"a/b":"foo", "c":"c", "d":"d"}
type overlayFS struct {
	overlay, base ReadDirReadFileFS
}

func dirAsReadDirReadFileFS(f fs.FS) (ReadDirReadFileFS, error) {
	fsys, ok := f.(ReadDirReadFileFS)
	if !ok {
		return nil, fmt.Errorf("Golang documentations claim that DirFS implements ReadDirFS and ReadFileFS")
	}
	return fsys, nil
}

// NewOverlayFS create an overlay FS from two FS.
func NewOverlayFS(override, base fs.FS) (ReadDirReadFileFS, error) {
	b, err := dirAsReadDirReadFileFS(base)
	if err != nil {
		return nil, err
	}
	if override == nil {
		return b, nil
	}
	o, err := dirAsReadDirReadFileFS(override)
	if err != nil {
		return nil, err
	}
	return overlayFS{overlay: o, base: b}, nil
}

// Open implements the main FS interface.
func (o overlayFS) Open(name string) (fs.File, error) {
	f, err := o.overlay.Open(name)
	f2, err2 := o.base.Open(name)
	if err != nil {
		return f2, err2
	}
	if err2 != nil {
		return f, err
	}
	overlay, ok := f.(fs.ReadDirFile)
	if !ok {
		// not a directory, we can return the unmerged overlay file
		return f, err
	}
	base, ok := f2.(fs.ReadDirFile)
	if !ok {
		// inconsistency between the two FSes, surprising.
		return nil, fmt.Errorf("Open(%q)'s base did not return ReadDirFile values", name)
	}
	// As a note, here we could have taken shortcuts and not implemented this:
	// implementations will realize that OverlayFS implements ReadDirFS and
	// call overylayfs.ReadDir(dir) instead of overlayfs.Open(dir)+d.ReadDir(int).
	//
	// That being said we do so to be a compliant FS and pass the fstest.TestFS
	// test battery.
	return &overlayDirFile{overlay: overlay, base: base}, nil
}

// ReadFile implements the ReadFileFS interface.
func (o overlayFS) ReadFile(name string) ([]byte, error) {
	b, err := o.overlay.ReadFile(name)
	if err == nil {
		return b, nil
	}
	return o.base.ReadFile(name)
}

// ReadDir implements the ReadDirFS interface.
func (o overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	a, err1 := o.overlay.ReadDir(name)
	b, err2 := o.base.ReadDir(name)
	return mergeReadDirs(a, b, err1, err2)
}

func mergeReadDirs(overlay, base []fs.DirEntry, errOverlay, errBase error) ([]fs.DirEntry, error) {
	if errOverlay != nil {
		// No need to merge (and handle both fs errors case).
		return base, errBase
	}
	var merged []fs.DirEntry
	seen := make(map[string]bool)
	for _, e := range overlay {
		seen[e.Name()] = true
		merged = append(merged, e)
	}
	for _, e := range base {
		if _, ok := seen[e.Name()]; !ok {
			merged = append(merged, e)
		}
	}
	return merged, nil
}

// ReadDirFile implementation when both overlay and base have an existing such
// directory.
type overlayDirFile struct {
	overlay, base fs.ReadDirFile
	initialized   bool
	entries       []fs.DirEntry
	offset        int
}

func (f *overlayDirFile) Stat() (fs.FileInfo, error) {
	return f.overlay.Stat()
}

func (f *overlayDirFile) Read(b []byte) (int, error) {
	// Will be an error: one can't read directories.
	return f.overlay.Read(b)
}

func (f *overlayDirFile) Close() error {
	err := f.overlay.Close()
	err2 := f.base.Close()
	return errors.Join(err, err2)
}

func (f *overlayDirFile) ReadDir(count int) ([]fs.DirEntry, error) {
	if !f.initialized {
		a, err1 := f.overlay.ReadDir(-1)
		b, err2 := f.base.ReadDir(-1)
		if err1 != nil || err2 != nil {
			panic("unexpected error")
		}
		var err error
		f.entries, err = mergeReadDirs(a, b, err1, err2)
		if err != nil {
			panic("unexpected error")
		}
		f.initialized = true
	}
	n := len(f.entries) - f.offset
	if n == 0 {
		if count <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}
	if count > 0 && n > count {
		n = count
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = f.entries[f.offset+i]
	}
	f.offset += n
	return list, nil
}

// Verifying interface implementations
var _ ReadDirReadFileFS = (*overlayFS)(nil)
var _ fs.ReadDirFile = (*overlayDirFile)(nil)
