// Package txtarfs turns a txtar into an fs.FS.
package txtarfs

import (
	"io/fs"

	"github.com/josharian/mapfs"
	"golang.org/x/tools/txtar"
)

// As returns an fs.FS containing ar's contents.
// Subsequent changes to ar may or may not
// be reflected in the returned fs.FS.
func As(ar *txtar.Archive) fs.FS {
	m := make(mapfs.MapFS, len(ar.Files))
	for _, f := range ar.Files {
		m[f.Name] = &mapfs.MapFile{
			Data: f.Data,
			// TODO: maybe ModTime: time.Now(),
			Sys: f,
		}
	}
	m.ChmodAll(0666)
	return m
}

// From constructs a txtar.Archive with the contents of fsys and an empty Comment.
// Subsequent changes to fsys are not reflected in the returned archive.
//
// The transformation is lossy.
// For example, because directories are implicit in txtar archives,
// empty directories in fsys will be lost.
// And txtar does not represent file mode, mtime, or other file metadata.
//
// Note also this warning from function txtar.Format:
//   > It is assumed that the Archive data structure is well-formed:
//   > a.Comment and all a.File[i].Data contain no file marker lines,
//   > and all a.File[i].Name is non-empty.
// From does not guarantee that a.File[i].Data contain no file marker lines.
//
// In short, it is unwise to use From/As as a generic filesystem serialization mechanism.
func From(fsys fs.FS) (*txtar.Archive, error) {
	ar := new(txtar.Archive)
	walkfn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Directories in txtar are implicit.
			return nil
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		ar.Files = append(ar.Files, txtar.File{Name: path, Data: data})
		return nil
	}

	err := fs.WalkDir(fsys, ".", walkfn)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
