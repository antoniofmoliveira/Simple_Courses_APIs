package handlers

import (
	"io"
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
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/octet-stream" {
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbCategory := fb.GetRootAsCategory(body, 0)
	categoryInputDto := dto.CategoryInputDto{
		ID:          string(fbCategory.Id()),
		Name:        string(fbCategory.Name()),
		Description: string(fbCategory.Description()),
	}
	categoryOutputDto, err := h.CategoryDB.Create(categoryInputDto)
	if err != nil {
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusCreated)
	buf := categoryAsBytes(&categoryOutputDto)
	w.Write(*buf)
	// fbBuilder := flatbuffers.NewBuilder(0)
	// id := fbBuilder.CreateString(categoryOutputDto.ID)
	// name := fbBuilder.CreateString(categoryOutputDto.Name)
	// description := fbBuilder.CreateString(categoryOutputDto.Description)
	// fb.CategoryStart(fbBuilder)
	// fb.CategoryAddId(fbBuilder, id)
	// fb.CategoryAddName(fbBuilder, name)
	// fb.CategoryAddDescription(fbBuilder, description)
	// fbCategoryOutput := fb.CategoryEnd(fbBuilder)
	// fbBuilder.Finish(fbCategoryOutput)
	// w.Write(fbBuilder.FinishedBytes())
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/octet-stream" {
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/octet-stream" {
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbCategory := fb.GetRootAsCategory(body, 0)
	categoryInputDto := dto.CategoryInputDto{
		ID:          string(fbCategory.Id()),
		Name:        string(fbCategory.Name()),
		Description: string(fbCategory.Description()),
	}
	err = h.CategoryDB.Update(categoryInputDto)
	if err != nil {
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendFlatBufferMessage(w, "category updated", http.StatusOK)
}

func (h *CategoryHandler) FIndAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/octet-stream" {
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	categories, err := h.CategoryDB.FindAll()
	if err != nil {
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	fbBuilder := flatbuffers.NewBuilder(512)

	elements := categoriesAsFlatBufferVector(fbBuilder, &categories.Categories)

	fb.CategoriesStartElementsVector(fbBuilder, len(*elements))
	for i := len(*elements) - 1; i >= 0; i-- {
		fbBuilder.PrependUOffsetT((*elements)[i])
	}
	vec := fbBuilder.EndVector(len(*elements))

	fb.CategoriesStart(fbBuilder)
	fb.CategoriesAddElements(fbBuilder, vec)
	fbCategoriesOutput := fb.CategoriesEnd(fbBuilder)
	fbBuilder.Finish(fbCategoriesOutput)
	w.Write(fbBuilder.FinishedBytes())
}

func categoriesAsFlatBufferVector(fbBuilder *flatbuffers.Builder, categories *[]dto.CategoryOutputDto) *[]flatbuffers.UOffsetT {
	var elements []flatbuffers.UOffsetT
	for _, category := range *categories {
		id := fbBuilder.CreateString(category.ID)
		name := fbBuilder.CreateString(category.Name)
		description := fbBuilder.CreateString(category.Description)
		fb.CategoryStart(fbBuilder)
		fb.CategoryAddId(fbBuilder, id)
		fb.CategoryAddName(fbBuilder, name)
		fb.CategoryAddDescription(fbBuilder, description)
		fbCategory := fb.CategoryEnd(fbBuilder)
		elements = append(elements, fbCategory)
		elements = append(elements, fbBuilder.Offset())
	}
	return &elements
}

func (h *CategoryHandler) FindCategory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/octet-stream" {
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}
	id := r.PathValue("id")
	// log.Println(id)
	category, err := h.CategoryDB.Find(id)
	if err != nil {
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	buf := categoryAsBytes(&category)
	// bb := flatbuffers.NewBuilder(1024)
	// idr := bb.CreateString(category.ID)
	// name := bb.CreateString(category.Name)
	// fb.CategoryStart(bb)
	// fb.CategoryAddId(bb, idr)
	// fb.CategoryAddName(bb, name)
	// fbCategoryOutput := fb.CategoryEnd(bb)
	// bb.Finish(fbCategoryOutput)
	// w.Write(bb.FinishedBytes())
	w.Write(*buf)
}

func categoryAsBytes(category *dto.CategoryOutputDto) *[]byte {
	bb := flatbuffers.NewBuilder(0)
	id := bb.CreateString(category.ID)
	name := bb.CreateString(category.Name)
	description := bb.CreateString(category.Description)
	fb.CategoryStart(bb)
	fb.CategoryAddId(bb, id)
	fb.CategoryAddName(bb, name)
	fb.CategoryAddDescription(bb, description)
	fbCategoryOutput := fb.CategoryEnd(bb)
	bb.Finish(fbCategoryOutput)
	buf := bb.FinishedBytes()
	return &buf
}
