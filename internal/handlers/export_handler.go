package handlers

import (
	"net/http"

	"github.com/zuudevs/saq-inventory-system-backend/internal/services"
)

type ExportHandler struct {
	ExportService *services.ExportService
}

func NewExportHandler(service *services.ExportService) *ExportHandler {
	return &ExportHandler{
		ExportService: service,
	}
}

func (h *ExportHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=exports.zip")

	if err := h.ExportService.ExportCSV(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ExportHandler) ExportXLSX(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=exports.xlsx")

	if err := h.ExportService.ExportXLSX(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
