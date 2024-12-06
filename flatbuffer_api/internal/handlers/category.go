package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/dto"
	"github.com/antoniofmoliveira/courses/flatbuffersapi/fb"
	flatbuffers "github.com/google/flatbuffers/go"
)

type CategoryHandler struct {
	CategoryRepository database.CategoryRepositoryInterface
}

func NewCategoryHandler(categoryDB database.CategoryRepositoryInterface) *CategoryHandler {
	return &CategoryHandler{
		CategoryRepository: categoryDB,
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	defer func() {
		if err := recover(); err != nil {
			slog.Info("createCategory", "msg", "unexpected payload")
			slog.Error("createCategory", "msg", err)
			sendFlatBufferMessage(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()

	if r.Header.Get("Content-Type") != octetStream {
		slog.Error("createCategory", "msg", "invalid content type")
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	if r.Header.Get("Accept") != octetStream {
		slog.Error("createCategory", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("createCategory", "msg", err)
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

	categoryOutputDto, err := h.CategoryRepository.Create(categoryInputDto)
	if err != nil {
		slog.Error("createCategory", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	buf := categoryAsBytes(&categoryOutputDto)
	w.Write(*buf)

	slog.Info("createCategory", "msg", "category created", "id", categoryOutputDto.ID)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	defer func() {
		if err := recover(); err != nil {
			slog.Info("updateCategory", "msg", "unexpected payload")
			slog.Error("updateCategory", "msg", err)
			sendFlatBufferMessage(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()

	if r.Header.Get("Content-Type") != octetStream {
		slog.Error("updateCategory", "msg", "invalid content type")
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	if r.Header.Get("Accept") != octetStream {
		slog.Error("updateCategory", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("updateCategory", "msg", err)
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

	err = h.CategoryRepository.Update(categoryInputDto)
	if err != nil {
		slog.Error("updateCategory", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendFlatBufferMessage(w, "category updated", http.StatusOK)

	slog.Info("updateCategory", "msg", "category updated", "id", categoryInputDto.ID)
}

func (h *CategoryHandler) FindAllCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("FindAllCategories", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	categories, err := h.CategoryRepository.FindAll()
	if err != nil {
		slog.Error("FindAllCategories", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fbBuilder := flatbuffers.NewBuilder(0)

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

	w.WriteHeader(http.StatusOK)
	w.Write(fbBuilder.FinishedBytes())

	slog.Info("FindAllCategories", "msg", "categories found", "count", len(categories.Categories))
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
	}
	return &elements
}

func (h *CategoryHandler) FindCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("FindCategory", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")

	category, err := h.CategoryRepository.Find(id)
	if err != nil {
		slog.Error("FindCategory", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf := categoryAsBytes(&category)
	w.Write(*buf)

	slog.Info("FindCategory", "msg", "category found", "id", id)
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

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("DeleteCategory", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")

	err := h.CategoryRepository.Delete(id)
	if err != nil {
		slog.Error("DeleteCategory", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendFlatBufferMessage(w, "category deleted", http.StatusOK)

	slog.Info("DeleteCategory", "msg", "category deleted", "id", id)
}
