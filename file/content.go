package file

type FileContent struct {
	Size    int
	Name    string
	Content []byte
	Index   int64
	Path    string
}

func NewFileContent(size int, name string, index int64, path string, content []byte) FileContent {
	return FileContent{
		Size:    size,
		Name:    name,
		Content: content,
		Index:   index,
		Path:    path,
	}
}
func (f *FileContent) Append(buffer []byte) {
	f.Content = append(f.Content, buffer...)
}
func (f *FileContent) getBytes() []byte {
	return f.Content
}
