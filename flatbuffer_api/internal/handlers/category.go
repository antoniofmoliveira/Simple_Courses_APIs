package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/dto"
	"github.com/antoniofmoliveira/courses/flatbuffersapi/fb"
	flatbuffers "github.com/google/flatbuffers/go"
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
	if r.Header.Get("Content-Type") != "application/octet-stream" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/octet-stream" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbCategory := fb.GetRootAsCategory(body, 0)
	categoryInputDto := dto.CategoryInputDto{
		ID:   string(fbCategory.Id()),
		Name: string(fbCategory.Name()),
	}
	categoryOutputDto, err := h.CategoryDB.Create(categoryInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusCreated)
	bb := flatbuffers.NewBuilder(0)
	fb.CategoryStart(bb)
	fb.CategoryAddId(bb, bb.CreateString(categoryOutputDto.ID))
	fb.CategoryAddName(bb, bb.CreateString(categoryOutputDto.Name))
	fbCategoryOutput := fb.CategoryEnd(bb)
	bb.Finish(fbCategoryOutput)
	w.Write(bb.FinishedBytes())
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/octet-stream" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/octet-stream" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbCategory := fb.GetRootAsCategory(body, 0)
	categoryInputDto := dto.CategoryInputDto{
		ID:   string(fbCategory.Id()),
		Name: string(fbCategory.Name()),
	}
	err = h.CategoryDB.Update(categoryInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)

}

func (h *CategoryHandler) FIndAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/octet-stream" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	categories, err := h.CategoryDB.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	bb := flatbuffers.NewBuilder(1024)
	var elements []flatbuffers.UOffsetT
	for _, category := range categories.Categories {
		id := bb.CreateString(category.ID)
		name := bb.CreateString(category.Name)
		fb.CategoryStart(bb)
		fb.CategoryAddId(bb, id)
		fb.CategoryAddName(bb, name)
		fbCategory := fb.CategoryEnd(bb)
		elements = append(elements, fbCategory)
		elements = append(elements, bb.Offset())
	}
	fb.CategoriesStartElementsVector(bb, len(elements))
	for i := len(elements) - 1; i >= 0; i-- {
		bb.PrependUOffsetT(elements[i])
	}
	vec := bb.EndVector(len(elements))

	fb.CategoriesStart(bb)
	fb.CategoriesAddElements(bb, vec)
	fbCategoriesOutput := fb.CategoriesEnd(bb)
	bb.Finish(fbCategoriesOutput)
	w.Write(bb.FinishedBytes())
}

func (h *CategoryHandler) FindCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/octet-stream" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}
	id := r.PathValue("id")
	log.Println(id)
	category, err := h.CategoryDB.Find(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	bb := flatbuffers.NewBuilder(1024)
	idr := bb.CreateString(category.ID)
	name := bb.CreateString(category.Name)
	fb.CategoryStart(bb)
	fb.CategoryAddId(bb, idr)
	fb.CategoryAddName(bb, name)
	fbCategoryOutput := fb.CategoryEnd(bb)
	bb.Finish(fbCategoryOutput)
	w.Write(bb.FinishedBytes())
}

// curl -H "Accept: application/octet-stream" http://localhost:8088/categories/d7f11a84-bb8a-44fd-ad1e-69409c0d8b4d --output file1.bin
