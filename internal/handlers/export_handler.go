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

func (h *ExportHandler) ExportItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=items.csv")

	if err := h.ExportService.ExportItemsToCSV(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
