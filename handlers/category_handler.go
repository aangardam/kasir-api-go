package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleCategorys - GET /api/category
func (h *CategoryHandler) HandleCategorys(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.ResponseSuccess(w, categories, http.StatusOK, "Categories retrieved successfully")
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.service.Create(&category)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseSuccess(w, category, http.StatusCreated, "Category created successfully")
}

// HandleCategoryByID - GET/PUT/DELETE /api/category/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetByID - GET /api/category/{id}
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		utils.ResponseError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.ResponseSuccess(w, category, http.StatusOK, "Category retrieved successfully")
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category.ID = id
	err = h.service.Update(&category)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseSuccess(w, category, http.StatusOK, "Category updated successfully")
}

// Delete - DELETE /api/category/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseSuccess(w, map[string]string{}, http.StatusOK, "Category deleted successfully")
}
