package test

import (
	"archive/zip"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/linkiog/doodocs/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func TestGetArchiveInformation(t *testing.T) {
	router := gin.Default()
	router.POST("/api/archive/information", handlers.GetArchiveInformation)

	zipFilePath, err := createTestZipFile()
	if err != nil {
		t.Fatalf("Failed to create test ZIP file: %v", err)
	}
	defer os.Remove(zipFilePath)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", zipFilePath)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	file, err := os.Open(zipFilePath)
	if err != nil {
		t.Fatalf("Failed to open test ZIP file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file content: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/api/archive/information", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	expectedFilename := "test-archive-"
	if !bytes.Contains(recorder.Body.Bytes(), []byte(expectedFilename)) {
		t.Errorf("Expected response to contain %s, got %s", expectedFilename, recorder.Body.String())
	}

}
func createTestZipFile() (string, error) {
	tmpFile, err := os.CreateTemp("", "test-archive-*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	w, err := zipWriter.Create("test-file.txt")
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte("This is a test file"))
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
