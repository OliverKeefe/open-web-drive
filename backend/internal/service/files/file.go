package files

import (
	"bytes"
	"io"
)

type File struct {
	Metadata FileMetadata
	Bytes    FileData
}

type FileData struct {
	Filename string
	Reader   io.ReadCloser
}

type FileResponse struct {
	Data bytes.Buffer `json:"data"`
}

//func (f *FileData) ToResponse() FileResponse {
//	return FileResponse{Data: f.Reader}
//}
