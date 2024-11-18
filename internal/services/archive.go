package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"strings"
)

func ProcessArchiveInformation(file multipart.File, header *multipart.FileHeader) (map[string]interface{}, error) {
	if !strings.HasSuffix(header.Filename, ".zip") {
		return nil, errors.New("Only ZIP files are allowed")
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return nil, errors.New("Invalid ZIP archive")
	}

	var totalSize float64
	var files []map[string]interface{}
	var firstFileName string

	rootFolderName := strings.TrimSuffix(header.Filename, ".zip")

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}

		filePath := f.Name
		if strings.HasPrefix(filePath, rootFolderName+"/") {
			filePath = strings.TrimPrefix(filePath, rootFolderName+"/")
		}

		if firstFileName == "" {
			firstFileName = filePath
		}

		totalSize += float64(f.UncompressedSize64)
		files = append(files, map[string]interface{}{
			"file_path": filePath,
			"size":      float64(f.UncompressedSize64),
			"mimetype":  detectMimeType(filePath),
		})
	}

	if len(files) == 0 {
		return nil, errors.New("Archive is empty")
	}

	return map[string]interface{}{
		"filename":     firstFileName,
		"archive_size": header.Size,
		"total_size":   totalSize,
		"total_files":  len(files),
		"files":        files,
	}, nil
}

func detectMimeType(filename string) string {
	if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(filename, ".png") {
		return "image/png"
	} else if strings.HasSuffix(filename, ".docx") {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	} else if strings.HasSuffix(filename, ".xml") {
		return "application/xml"
	}
	return "application/octet-stream"
}
func CreateArchiveFromFiles(files []*multipart.FileHeader) ([]byte, error) {
	if len(files) == 0 {
		return nil, errors.New("No files provided for archiving")
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(zipFile, file)
		if err != nil {
			return nil, err
		}
	}
	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
