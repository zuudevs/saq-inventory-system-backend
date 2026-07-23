package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zuudevs/saq-inventory-system-backend/internal/dto"
	"github.com/zuudevs/saq-inventory-system-backend/internal/version"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		Status:    "ok",
		Service:   version.Service,
		Version:   version.Version,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}