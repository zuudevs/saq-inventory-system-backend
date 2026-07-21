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
}

func NewImageHandler(service *services.ImageService) *ImageHandler {
	return &ImageHandler{
		ImageService: service,
	}
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
