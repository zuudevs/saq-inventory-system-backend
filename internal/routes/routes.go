package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zuudevs/saq-inventory-system-backend/internal/handlers"
)

func HealthRoutes(r chi.Router) {
	r.Get("/health", handlers.HealthHandler)
}

func BrandRoutes(r chi.Router, h *handlers.BrandHandler) {
	r.Route("/brands", func(r chi.Router) {
		r.Get("/", h.FindAll)
		r.Get("/{id}", h.FindByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func CategoryRoutes(r chi.Router, h *handlers.CategoryHandler) {
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", h.FindAll)
		r.Get("/{id}", h.FindByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func LocationRoutes(r chi.Router, h *handlers.LocationHandler) {
	r.Route("/locations", func(r chi.Router) {
		r.Get("/", h.FindAll)
		r.Get("/{id}", h.FindByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func ItemRoutes(r chi.Router, h *handlers.ItemHandler) {
	r.Route("/items", func(r chi.Router) {
		r.Get("/", h.FindAll)
		r.Get("/{id}", h.FindByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

func ImageRoutes(r chi.Router, h *handlers.ImageHandler) {
	r.Route("/images", func(r chi.Router) {
		r.Get("/", h.FindAll)
		r.Get("/{id}", h.FindByID)
		r.Post("/", h.Create)
		r.Post("/upload", h.Upload)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// StorageRoutes menge-serve file statis di dalam storagePath (hasil
// SaveImageFile) lewat prefix "/storage/". image_path yang tersimpan di DB
// (mis. "images/<uuid>.png") jadi bisa diakses via GET /storage/images/<uuid>.png.
// http.Dir + http.FileServer sudah aman terhadap path traversal secara
// bawaan (path dibersihkan sebelum dibuka), jadi tidak perlu validasi manual.
func StorageRoutes(r chi.Router, storagePath string) {
	fileServer := http.FileServer(http.Dir(storagePath))
	r.Handle("/storage/*", http.StripPrefix("/storage/", fileServer))
}

func MetadataStructureRoutes(r chi.Router, h *handlers.MetadataStructureHandler) {
	r.Route("/categories/{categoryId}/metadata-structure", func(r chi.Router) {
		r.Get("/", h.FindByCategoryID)
		r.Post("/", h.Create)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})
}

func ExportRoutes(r chi.Router, h *handlers.ExportHandler) {
	r.Route("/exports", func(r chi.Router) {
		r.Get("/items", h.ExportItems)
		r.Get("/items/xlsx", h.ExportItemsXLSX)
		r.Get("/items.xlsx", h.ExportItemsXLSX)
	})
}