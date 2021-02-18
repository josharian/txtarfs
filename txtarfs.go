// Package txtarfs turns a txtar into an fs.FS.
package txtarfs

import (
	"io/fs"
	"testing/fstest"

	"golang.org/x/tools/txtar"
)

func As(ar *txtar.Archive) fs.FS {
	m := make(fstest.MapFS, len(ar.Files))
	for _, f := range ar.Files {
		m[f.Name] = &fstest.MapFile{
			Data: f.Data,
			Mode: 0666,
			// TODO: maybe ModTime: time.Now(),
			Sys: f,
		}
	}
	return m
}
