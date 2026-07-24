package handlers

import (
	"net/http"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/services"
	"github.com/zuudevs/saq-inventory-system-backend/internal/utils"
)

type ImportHandler struct {
	ImportService *services.ImportService
}

func NewImportHandler(service *services.ImportService) *ImportHandler {
	return &ImportHandler{
		ImportService: service,
	}
}

func (h *ImportHandler) ImportXLSX(w http.ResponseWriter, r *http.Request) {
	if r.Header == nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("missing header"),
		)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid multipart form or file too large"),
		)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("file field is required in multipart form"),
		)
		return
	}
	defer file.Close()

	summary, err := h.ImportService.ImportXLSX(file)
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
		http.StatusOK,
		dto.Success("import completed successfully", summary),
	)
}
