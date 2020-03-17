package web

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
)

func MultipartFileToOsFile(file multipart.File) *bytes.Buffer {
	bufFile := bytes.NewBuffer(nil)
	_, err := io.Copy(bufFile, file)
	if err != nil {
		log.Println(err)
	}

	return bufFile
}
