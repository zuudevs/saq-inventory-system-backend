package handlers

import (
	"net/http"
	"time"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/utils"
	"github.com/zuudevs/saq-inventory-system-backend/internal/version"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		Status:    "ok",
		Service:   version.Service,
		Version:   version.Version,
		Timestamp: time.Now().UTC(),
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success("", response),
	)
}