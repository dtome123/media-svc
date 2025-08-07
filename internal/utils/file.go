package utils

import (
	"mime"
	"path/filepath"
)

func DetectContentTypeByFileName(fileName string) string {
	ext := filepath.Ext(fileName)

	if contentType := mime.TypeByExtension(ext); contentType != "" {
		return contentType
	}

	return ""
}
