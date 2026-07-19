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

type BrandHandler struct {
	BrandService *services.BrandService
}

func NewBrandHandler(service *services.BrandService) *BrandHandler {
	return &BrandHandler{
		BrandService: service,
	}
}

func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBrandRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid request body"),
		)
		return
	}

	brand, err := h.BrandService.Create(req)
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
			"brand created successfully",
			brand,
		),
	)
}

func (h *BrandHandler) FindAll(w http.ResponseWriter, _ *http.Request) {
	brands, err := h.BrandService.FindAll()
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
			"brands retrieved successfully",
			brands,
		),
	)
}

func (h *BrandHandler) FindByID(w http.ResponseWriter, r *http.Request) {
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

	brand, err := h.BrandService.FindByID(id)
	if err != nil {
		utils.JSON(
			w,
			http.StatusInternalServerError,
			dto.Error[any](err.Error()),
		)
		return
	}

	if brand == nil {
		utils.JSON(
			w,
			http.StatusNotFound,
			dto.Error[any]("brand not found"),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"brand retrieved successfully",
			brand,
		),
	)
}

func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateBrandRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid request body"),
		)
		return
	}

	brand, err := h.BrandService.Update(id, req)
	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any](err.Error()),
		)
		return
	}

	if brand == nil {
		utils.JSON(
			w,
			http.StatusNotFound,
			dto.Error[any]("brand not found"),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"brand updated successfully",
			brand,
		),
	)
}

func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.BrandService.Delete(id); err != nil {
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
			"brand deleted successfully",
			nil,
		),
	)
}
