package file

import "os"

type Meta struct {
	os.FileInfo
	Index    int64
	Path     string
	FullPath string
}
