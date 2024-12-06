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

type CourseHandler struct {
	CourseRepository database.CourseRepositoryInterface
}

func NewCourseHandler(courseRepository database.CourseRepositoryInterface) *CourseHandler {
	return &CourseHandler{
		CourseRepository: courseRepository,
	}
}

func (c *CourseHandler) FindCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("FindCourse", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")

	course, err := c.CourseRepository.Find(id)
	if err != nil {
		slog.Error("FindCourse", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf := courseAsBytes(&course)
	w.Write(*buf)

	slog.Info("FindCourse", "msg", "course found", "id", id)
}

func courseAsBytes(course *dto.CourseOutputDto) *[]byte {
	bb := flatbuffers.NewBuilder(0)
	id := bb.CreateString(course.ID)
	name := bb.CreateString(course.Name)
	category_id := bb.CreateString(course.CategoryID)
	description := bb.CreateString(course.Description)
	fb.CourseStart(bb)
	fb.CourseAddId(bb, id)
	fb.CourseAddName(bb, name)
	fb.CourseAddDescription(bb, description)
	fb.CourseAddCategoryId(bb, category_id)
	fbCourseOutput := fb.CourseEnd(bb)
	bb.Finish(fbCourseOutput)
	buf := bb.FinishedBytes()
	return &buf
}

func (c *CourseHandler) FindAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("FindAllCourses", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	courses, err := c.CourseRepository.FindAll()
	if err != nil {
		slog.Error("FindAllCourses", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fbuilder := flatbuffers.NewBuilder(0)

	elements := coursesAsFlatBufferVector(fbuilder, &courses.Courses)

	fb.CoursesStartElementsVector(fbuilder, len(*elements))
	for i := len(*elements) - 1; i >= 0; i-- {
		fbuilder.PrependUOffsetT((*elements)[i])
	}
	vec := fbuilder.EndVector(len(*elements))

	fb.CoursesStart(fbuilder)
	fb.CoursesAddElements(fbuilder, vec)
	fbCoursesOutput := fb.CoursesEnd(fbuilder)
	fbuilder.Finish(fbCoursesOutput)

	w.WriteHeader(http.StatusOK)
	w.Write(fbuilder.FinishedBytes())

	slog.Info("FindAllCourses", "msg", "courses found", "count", len(courses.Courses))
}

func coursesAsFlatBufferVector(fbuilder *flatbuffers.Builder, courses *[]dto.CourseOutputDto) *[]flatbuffers.UOffsetT {
	elements := make([]flatbuffers.UOffsetT, 0)
	for _, course := range *courses {
		id := fbuilder.CreateString(course.ID)
		name := fbuilder.CreateString(course.Name)
		category_id := fbuilder.CreateString(course.CategoryID)
		description := fbuilder.CreateString(course.Description)
		fb.CourseStart(fbuilder)
		fb.CourseAddId(fbuilder, id)
		fb.CourseAddName(fbuilder, name)
		fb.CourseAddDescription(fbuilder, description)
		fb.CourseAddCategoryId(fbuilder, category_id)
		fbCourseOutput := fb.CourseEnd(fbuilder)
		elements = append(elements, fbCourseOutput)
	}
	return &elements
}

func (c *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	defer func() {
		if err := recover(); err != nil {
			slog.Info("createCourse", "msg", "unexpected payload")
			slog.Error("createCourse", "msg", err)
			sendFlatBufferMessage(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()

	if r.Header.Get("Content-Type") != octetStream {
		slog.Error("createCourse", "msg", "invalid content type")
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	if r.Header.Get("Accept") != octetStream {
		slog.Error("createCourse", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("createCourse", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbCourse := fb.GetRootAsCourse(body, 0)
	courseInputDto := dto.CourseInputDto{
		ID:          string(fbCourse.Id()),
		Name:        string(fbCourse.Name()),
		Description: string(fbCourse.Description()),
		CategoryID:  string(fbCourse.CategoryId()),
	}

	course, err := c.CourseRepository.Create(courseInputDto)
	if err != nil {
		slog.Error("createCourse", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf := courseAsBytes(course)
	w.Write(*buf)

	slog.Info("createCourse", "msg", "course created", "id", course.ID)
}

func (c *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	defer func() {
		if err := recover(); err != nil {
			slog.Info("updateCourse", "msg", "unexpected payload")
			slog.Error("updateCourse", "msg", err)
			sendFlatBufferMessage(w, err.(error).Error(), http.StatusInternalServerError)
		}
	}()

	if r.Header.Get("Content-Type") != octetStream {
		slog.Error("updateCourse", "msg", "invalid content type")
		sendFlatBufferMessage(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	if r.Header.Get("Accept") != octetStream {
		slog.Error("updateCourse", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("updateCourse", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fbCourse := fb.GetRootAsCourse(body, 0)
	courseInputDto := dto.CourseInputDto{
		ID:          string(fbCourse.Id()),
		Name:        string(fbCourse.Name()),
		Description: string(fbCourse.Description()),
		CategoryID:  string(fbCourse.CategoryId()),
	}

	err = c.CourseRepository.Update(courseInputDto)
	if err != nil {
		slog.Error("updateCourse", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendFlatBufferMessage(w, "course updated", http.StatusOK)

	slog.Info("updateCourse", "msg", "course updated", "id", courseInputDto.ID)
}

func (c *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", octetStream)

	if r.Header.Get("Accept") != octetStream {
		slog.Error("DeleteCourse", "msg", "invalid accept header")
		sendFlatBufferMessage(w, "invalid accept header", http.StatusUnsupportedMediaType)
		return
	}

	id := r.PathValue("id")	

	err := c.CourseRepository.Delete(id)
	if err != nil {
		slog.Error("DeleteCourse", "msg", err)
		sendFlatBufferMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	sendFlatBufferMessage(w, "course deleted", http.StatusOK)

	slog.Info("DeleteCourse", "msg", "course deleted", "id", id)
}