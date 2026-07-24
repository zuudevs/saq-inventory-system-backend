package utils

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/HugoSmits86/nativewebp"
	"github.com/google/uuid"
)

var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

// MaxImageUploadSize adalah batas ukuran file gambar yang diupload (5 MB),
// dipakai baik untuk validasi per-file maupun sebagai limit ParseMultipartForm.
const MaxImageUploadSize = 5 << 20 // 5 MB

// SaveImageFile memvalidasi ekstensi & ukuran, lalu mengonversi gambar ke WebP
// dan menyimpannya ke {storageRoot}/images/<uuid>.webp (menggunakan pure Go nativewebp).
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

	if err := nativewebp.Encode(dst, img, nil); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("failed to encode to webp: %w", err)
	}

	return filepath.ToSlash(filepath.Join("images", filename)), nil
}

// DeleteImageFile menghapus file fisik di storageRoot/{relativePath}.
// Dipanggil best-effort (error diabaikan pemanggil) dari ImageService saat
// record image dihapus atau image_path-nya diganti — supaya tidak ada file
// yatim (orphan) yang menumpuk di disk. Sengaja idempotent: kalau filenya
// memang sudah tidak ada, itu bukan error.
func DeleteImageFile(storageRoot string, relativePath string) error {
	if relativePath == "" {
		return nil
	}

	fullPath := filepath.Join(storageRoot, filepath.FromSlash(relativePath))

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
