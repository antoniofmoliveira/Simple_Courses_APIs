package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/dto"
)

type CourseHandler struct {
	CourseDB database.CourseRepositoryInterface
}

func NewCourseHandler(courseDB database.CourseRepositoryInterface) *CourseHandler {
	return &CourseHandler{
		CourseDB: courseDB,
	}
}

func (c *CourseHandler) FindAllCourses(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	courses, err := c.CourseDB.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}

func (c *CourseHandler) FindCourse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	course, err := c.CourseDB.Find(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

func (c *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	var courseInputDto dto.CourseInputDto
	err := json.NewDecoder(r.Body).Decode(&courseInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	courseOutputDto, err := c.CourseDB.Create(courseInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(courseOutputDto)
}

func (c *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	var courseInputDto dto.CourseInputDto
	err := json.NewDecoder(r.Body).Decode(&courseInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.CourseDB.Update(courseInputDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")
	err := c.CourseDB.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
