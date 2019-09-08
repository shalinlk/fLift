package file

type FileContent struct {
	Size    int
	Name    string
	Content []byte
}

func NewFileContent(size int, name string) FileContent {
	return FileContent{size, name, make([]byte, 0)}
}
func (f *FileContent) Append(buffer []byte) {
	f.Content = append(f.Content, buffer...)
}
func (f *FileContent) getBytes() []byte {
	return f.Content
}
