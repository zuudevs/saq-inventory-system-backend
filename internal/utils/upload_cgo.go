//go:build cgo

package utils

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// SaveImageFile memvalidasi ekstensi & ukuran, lalu mengonversi gambar ke WebP 
// dan menyimpannya ke {storageRoot}/images/<uuid>.webp (menggunakan go-webp dengan CGO).
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

	// Jika file asal sudah WebP, simpan langsung tanpa re-encode
	if ext == ".webp" {
		filename := uuid.NewString() + ".webp"
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

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	filename := uuid.NewString() + ".webp"
	fullPath := filepath.Join(imagesDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 85)
	if err != nil {
		os.Remove(fullPath)
		return "", err
	}

	if err := webp.Encode(dst, img, options); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("failed to encode to webp: %w", err)
	}

	return filepath.ToSlash(filepath.Join("images", filename)), nil
}
