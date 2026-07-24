package utils_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/zuudevs/saq-inventory-system-backend/internal/utils"
)

func createDummyPNG(t *testing.T) ([]byte, string) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode dummy png: %v", err)
	}

	return buf.Bytes(), "test.png"
}

func createMultipartFile(t *testing.T, content []byte, filename string) (multipart.File, *multipart.FileHeader, func()) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}
	writer.Close()

	req, err := multipart.NewReader(body, writer.Boundary()).ReadForm(10 << 20)
	if err != nil {
		t.Fatalf("failed to read form: %v", err)
	}

	files := req.File["file"]
	if len(files) == 0 {
		t.Fatalf("no files found in multipart form")
	}

	header := files[0]
	file, err := header.Open()
	if err != nil {
		t.Fatalf("failed to open header: %v", err)
	}

	cleanup := func() {
		file.Close()
		req.RemoveAll()
	}

	return file, header, cleanup
}

func TestSaveImageFile_PNG_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upload_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content, filename := createDummyPNG(t)
	file, header, cleanup := createMultipartFile(t, content, filename)
	defer cleanup()

	relPath, err := utils.SaveImageFile(tempDir, file, header)
	if err != nil {
		t.Fatalf("SaveImageFile returned error: %v", err)
	}

	if relPath == "" {
		t.Errorf("expected relative path, got empty string")
	}

	fullPath := filepath.Join(tempDir, filepath.FromSlash(relPath))
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Errorf("saved file does not exist at %s", fullPath)
	}
}

func TestSaveImageFile_InvalidExtension(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upload_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file, header, cleanup := createMultipartFile(t, []byte("fake text"), "document.txt")
	defer cleanup()

	_, err = utils.SaveImageFile(tempDir, file, header)
	if err == nil {
		t.Errorf("expected error for unsupported file extension, got nil")
	}
}

func TestSaveImageFile_ExceedsSizeLimit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upload_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content, filename := createDummyPNG(t)
	file, header, cleanup := createMultipartFile(t, content, filename)
	defer cleanup()

	// Artificially inflate header.Size to test size check
	header.Size = utils.MaxImageUploadSize + 1

	_, err = utils.SaveImageFile(tempDir, file, header)
	if err == nil {
		t.Errorf("expected error for file exceeding size limit, got nil")
	}
}

func TestDeleteImageFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "upload_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	imagesDir := filepath.Join(tempDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		t.Fatalf("failed to create images dir: %v", err)
	}

	dummyPath := filepath.Join(imagesDir, "test_delete.webp")
	if err := os.WriteFile(dummyPath, []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}

	relPath := "images/test_delete.webp"
	if err := utils.DeleteImageFile(tempDir, relPath); err != nil {
		t.Fatalf("DeleteImageFile failed: %v", err)
	}

	if _, err := os.Stat(dummyPath); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted, but it still exists")
	}

	// Test idempotency: deleting non-existent file should not return error
	if err := utils.DeleteImageFile(tempDir, relPath); err != nil {
		t.Errorf("expected nil error on deleting non-existent file, got: %v", err)
	}
}
