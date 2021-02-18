package txtarfs_test

import (
	"bytes"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

func TestBasics(t *testing.T) {
	tests := []map[string]string{
		nil,
		{"x.txt": "hi"},
		{"a/x.txt": "hi"},
		{"a/x.txt, b/y.txt": "hello"},
		{"a/b/c/x.txt": ""},
	}

	for _, tt := range tests {
		ar := new(txtar.Archive)
		var names []string
		for name, data := range tt {
			ar.Files = append(ar.Files, txtar.File{Name: name, Data: []byte(data)})
			names = append(names, name)
		}
		arfs := txtarfs.As(ar)
		if err := fstest.TestFS(arfs, names...); err != nil {
			t.Fatal(err)
		}
		for name, data := range tt {
			out, err := fs.ReadFile(arfs, name)
			if err != nil {
				t.Errorf("fs.ReadFile(%s) = _, %v", name, err)
				continue
			}
			if !bytes.Equal([]byte(data), out) {
				t.Errorf("fs.ReadFile(%s) = %s want %s", name, out, data)
			}
		}
	}
}
