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

type ItemHandler struct {
	ItemService *services.ItemService
}

func NewItemHandler(service *services.ItemService) *ItemHandler {
	return &ItemHandler{
		ItemService: service,
	}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid request body"),
		)
		return
	}

	item, err := h.ItemService.Create(req)
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
			"item created successfully",
			item,
		),
	)
}

func (h *ItemHandler) FindAll(w http.ResponseWriter, _ *http.Request) {
	items, err := h.ItemService.FindAll()
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
			"items retrieved successfully",
			items,
		),
	)
}

func (h *ItemHandler) FindByID(w http.ResponseWriter, r *http.Request) {
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

	item, err := h.ItemService.FindByID(id)
	if err != nil {
		utils.JSON(
			w,
			http.StatusInternalServerError,
			dto.Error[any](err.Error()),
		)
		return
	}

	if item == nil {
		utils.JSON(
			w,
			http.StatusNotFound,
			dto.Error[any]("item not found"),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"item retrieved successfully",
			item,
		),
	)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any]("invalid request body"),
		)
		return
	}

	item, err := h.ItemService.Update(id, req)
	if err != nil {
		utils.JSON(
			w,
			http.StatusBadRequest,
			dto.Error[any](err.Error()),
		)
		return
	}

	if item == nil {
		utils.JSON(
			w,
			http.StatusNotFound,
			dto.Error[any]("item not found"),
		)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		dto.Success(
			"item updated successfully",
			item,
		),
	)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.ItemService.Delete(id); err != nil {
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
			"item deleted successfully",
			nil,
		),
	)
}
