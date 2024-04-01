package txtarfs_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"sort"
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
		sort.Slice(ar.Files, func(i, j int) bool { return ar.Files[i].Name < ar.Files[j].Name })
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

		err := fs.WalkDir(arfs, ".", func(path string, d fs.DirEntry, err error) error {
			fi, err := d.Info()
			if err != nil {
				return err
			}
			if mode := fi.Mode().Perm(); mode&0444 != 0444 {
				return fmt.Errorf("%s has mode %v", path, mode)
			}
			return nil
		})
		if err != nil {
			t.Errorf("fs.WalkDir(...) = %v", err)
			continue
		}

		ar2, err := txtarfs.From(arfs)
		if err != nil {
			t.Errorf("failed to write fsys for %v: %v", tt, err)
			continue
		}
		sort.Slice(ar2.Files, func(i, j int) bool { return ar2.Files[i].Name < ar2.Files[j].Name })
		in := string(txtar.Format(ar))
		out := string(txtar.Format(ar2))
		if in != out {
			t.Errorf("As/From round trip failed: %v != %v", in, out)
		}
	}
}
