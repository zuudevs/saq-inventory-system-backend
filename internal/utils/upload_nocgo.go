//go:build !cgo

package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// SaveImageFile memvalidasi ekstensi & ukuran, lalu menyimpan file asli 
// ke {storageRoot}/images/<uuid>.<ext> (digunakan sebagai fallback saat CGO tidak aktif).
func SaveImageFile(storageRoot string, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > MaxImageUploadSize {
		return "", fmt.Errorf("file size exceeds %d bytes limit", MaxImageUploadSize)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("unsupported file extension: %s", ext)
	}

	imagesDir := filepath.Join(storageRoot, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return "", err
	}

	filename := uuid.NewString() + ext
	fullPath := filepath.Join(imagesDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(fullPath)
		return "", err
	}

	return filepath.ToSlash(filepath.Join("images", filename)), nil
}
