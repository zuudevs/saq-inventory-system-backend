package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/services"
	"github.com/zuudevs/saq-inventory-system-backend/internal/utils"
)

type ImageHandler struct {
	ImageService *services.ImageService
	StoragePath  string
}

func NewImageHandler(service *services.ImageService, storagePath string) *ImageHandler {
	return &ImageHandler{
		ImageService: service,
		StoragePath:  storagePath,
	}
}

// Upload menerima multipart/form-data dengan field "file", menyimpannya ke
// disk, dan mengembalikan path relatifnya. Endpoint ini terpisah dari
// Create/Update dengan sengaja: client upload file dulu lewat sini, dapat
// image_path, baru kirim image_path itu ke POST/PUT /images — supaya upload
// yang gagal validasi owner (item/location tidak ada) tidak perlu bongkar
// pasang multipart request lagi, dan supaya ganti gambar pada image yang
// sudah ada bisa upload dulu sebelum PUT.
func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Header == nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("missing file header"),
		)
		return 
	}

	if err := r.ParseMultipartForm(utils.MaxImageUploadSize); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid multipart form or file too large"),
		)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("file is required"),
		)
		return
	}
	defer file.Close()

	relativePath, err := utils.SaveImageFile(h.StoragePath, file, header)
	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any](err.Error()),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusCreated,
		dto.Success(
			"file uploaded successfully",
			dto.UploadImageResponse{ImagePath: relativePath},
		),
	)
}

func (h *ImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateImageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid request body"),
		)
		return
	}

	image, err := h.ImageService.Create(req)
	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any](err.Error()),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusCreated,
		dto.Success(
			"image created successfully",
			image,
		),
	)
}

// FindAll mengembalikan semua image, atau bisa difilter dengan query param
// ?item_id= atau ?location_id= untuk mengambil galeri milik satu item/location
// saja (mis. dipakai halaman detail item untuk menampilkan foto-fotonya).
func (h *ImageHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if raw := query.Get("item_id"); raw != "" {
		itemID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			utils.JSON(
				w,
				http.StatusBadRequest,
				dto.Error[any]("invalid item_id"),
			)
			return
		}

		images, err := h.ImageService.FindByItemID(itemID)
		if err != nil {
			utils.JSON(
				w,
				http.StatusInternalServerError,
				dto.Error[any](err.Error()),
			)
			return
		}

		utils.JSON(
			w,
			http.StatusOK,
			dto.Success(
				"images retrieved successfully",
				images,
			),
		)
		return
	}

	if raw := query.Get("location_id"); raw != "" {
		locationID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			utils.JSON(
				w,
				http.StatusBadRequest,
				dto.Error[any]("invalid location_id"),
			)
			return
		}

		images, err := h.ImageService.FindByLocationID(locationID)
		if err != nil {
			utils.JSON(
				w,
				http.StatusInternalServerError,
				dto.Error[any](err.Error()),
			)
			return
		}

		utils.JSON(
			w,
			http.StatusOK,
			dto.Success(
				"images retrieved successfully",
				images,
			),
		)
		return
	}

	images, err := h.ImageService.FindAll()
	if err != nil {
		utils.JSON(
			w,
			http.StatusInternalServerError,
			dto.Error[any](err.Error()),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"images retrieved successfully",
			images,
		),
	)
}

func (h *ImageHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(
		chi.URLParam(r, "id"),
		10,
		64,
	)

	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid id"),
		)
		return
	}

	image, err := h.ImageService.FindByID(id)
	if err != nil {
		utils.JSON(
			w,
			http.StatusInternalServerError,
			dto.Error[any](err.Error()),
		)
		return
	}

	if image == nil {
		utils.JSON(
			w,
			http.StatusNotFound,
			dto.Error[any]("image not found"),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"image retrieved successfully",
			image,
		),
	)
}

func (h *ImageHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(
		chi.URLParam(r, "id"),
		10,
		64,
	)

	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid id"),
		)
		return
	}

	var req dto.UpdateImageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid request body"),
		)
		return
	}

	image, err := h.ImageService.Update(id, req)
	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any](err.Error()),
		)
		return
	}

	if image == nil {
		utils.JSON(
			w,
			http.StatusNotFound,
			dto.Error[any]("image not found"),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"image updated successfully",
			image,
		),
	)
}

func (h *ImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(
		chi.URLParam(r, "id"),
		10,
		64,
	)

	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid id"),
		)
		return
	}

	if err := h.ImageService.Delete(id); err != nil {
		utils.JSON(
			w,
			http.StatusInternalServerError,
			dto.Error[any](err.Error()),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success[any](
			"image deleted successfully",
			nil,
		),
	)
}
