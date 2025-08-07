package utils

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
)

func DetectContentTypeByFileName(fileName string) string {
	ext := filepath.Ext(fileName)

	if contentType := mime.TypeByExtension(ext); contentType != "" {
		return contentType
	}

	return ""
}

func WriteFile(folderPath, fileName string, data []byte) error {

	folderPath = filepath.Clean(folderPath)

	// Create folder if not exists
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Build full file path
	fullPath := filepath.Join(folderPath, fileName)

	// Write to file using os
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func RemoveFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil
}
