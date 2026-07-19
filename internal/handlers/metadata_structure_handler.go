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

type MetadataStructureHandler struct {
	MetadataStructureService *services.MetadataStructureService
}

func NewMetadataStructureHandler(service *services.MetadataStructureService) *MetadataStructureHandler {
	return &MetadataStructureHandler{
		MetadataStructureService: service,
	}
}

func (h *MetadataStructureHandler) Create(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(
		chi.URLParam(r, "categoryId"),
		10,
		64,
	)

	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid category id"),
		)
		return
	}

	var req dto.CreateMetadataStructureRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid request body"),
		)
		return
	}

	structure, err := h.MetadataStructureService.Create(categoryID, req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		}

		utils.JSON(
			w,
			status,
			dto.Error[any](err.Error()),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusCreated,
		dto.Success(
			"metadata structure created successfully",
			structure,
		),
	)
}

func (h *MetadataStructureHandler) FindByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(
		chi.URLParam(r, "categoryId"),
		10,
		64,
	)

	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid category id"),
		)
		return
	}

	structure, err := h.MetadataStructureService.FindByCategoryID(categoryID)
	if err != nil {
		utils.JSON(
			w,
			http.StatusInternalServerError,
			dto.Error[any](err.Error()),
		)
		return
	}

	if structure == nil {
		utils.JSON(
			w,
			http.StatusNotFound,
			dto.Error[any]("metadata structure not found"),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"metadata structure retrieved successfully",
			structure,
		),
	)
}
