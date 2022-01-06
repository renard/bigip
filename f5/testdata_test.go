package f5

import (
	"io/fs"
	"path/filepath"
)

func getFiles(fsys fs.FS) (paths []string) {
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".conf" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return
}
