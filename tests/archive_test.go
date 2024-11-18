package test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linkiog/doodocs/internal/api/handlers"
)

func TestCreateArchive(t *testing.T) {
	router := gin.Default()
	router.POST("/api/archive/files", handlers.CreateArchive)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file1, _ := os.CreateTemp("", "file1-*.txt")
	defer os.Remove(file1.Name())
	file1.WriteString("Test file 1")

	file2, _ := os.CreateTemp("", "file2-*.txt")
	defer os.Remove(file2.Name())
	file2.WriteString("Test file 2")

	part, _ := writer.CreateFormFile("files[]", file1.Name())
	part.Write([]byte("dummy content"))
	part, _ = writer.CreateFormFile("files[]", file2.Name())
	part.Write([]byte("dummy content"))

	writer.Close()

	req := httptest.NewRequest("POST", "/api/archive/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("Content-Type") != "application/zip" {
		t.Errorf("Expected Content-Type application/zip, got %s", recorder.Header().Get("Content-Type"))
	}
}
