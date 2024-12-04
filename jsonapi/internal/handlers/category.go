package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/dto"
)

type CategoryHandler struct {
	CategoryDB database.CategoryRepositoryInterface
}

func NewCategoryHandler(categoryDB database.CategoryRepositoryInterface) *CategoryHandler {
	return &CategoryHandler{
		CategoryDB: categoryDB,
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	var categoryInputDto dto.CategoryInputDto
	err := json.NewDecoder(r.Body).Decode(&categoryInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoryOutputDto, err := h.CategoryDB.Create(categoryInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(categoryOutputDto)
}

func (h *CategoryHandler) FindCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	category, err := h.CategoryDB.Find(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) FindAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	categories, err := h.CategoryDB.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	var categoryInputDto dto.CategoryInputDto
	err := json.NewDecoder(r.Body).Decode(&categoryInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.CategoryDB.Update(categoryInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")
	err := h.CategoryDB.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
