package utils

import (
	"os"
	"path/filepath"
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
