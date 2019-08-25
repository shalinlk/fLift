package file

type FileContent struct {
	Size    int64
	Name    string
	content []byte
}

func NewFileContent(size int64, name string) FileContent {
	return FileContent{size, name, make([]byte, 0)}
}
func (f *FileContent) Append(buffer []byte) {
	f.content = append(f.content, buffer...)
}
func (f *FileContent) getBytes() []byte {
	return f.content
}
