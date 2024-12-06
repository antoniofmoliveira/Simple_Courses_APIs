package mariadb

import (
	"database/sql"

	"github.com/antoniofmoliveira/courses/dto"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Course struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *Course {
	c := &Course{db: db}
	c.db.Exec("CREATE TABLE IF NOT EXISTS courses (id CHAR(36) PRIMARY KEY, name TEXT, description TEXT, category_id TEXT)")
	return c
}

func (c *Course) Create(course dto.CourseInputDto) (*dto.CourseOutputDto, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO courses (id, name, description, category_id) VALUES (?, ?, ?, ?)",
		id, course.Name, course.Description, course.CategoryID)
	if err != nil {
		return nil, err
	}
	return &dto.CourseOutputDto{
		ID:          id,
		Name:        course.Name,
		Description: course.Description,
		CategoryID:  course.CategoryID,
	}, nil
}

func (c *Course) FindAll() (dto.CourseListOutputDto, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses")
	if err != nil {
		return dto.CourseListOutputDto{}, err
	}
	defer rows.Close()
	courses := dto.CourseListOutputDto{}
	for rows.Next() {
		var id, name, description, categoryID string
		if err := rows.Scan(&id, &name, &description, &categoryID); err != nil {
			return dto.CourseListOutputDto{}, err
		}
		courses.Courses = append(courses.Courses, dto.CourseOutputDto{ID: id, Name: name, Description: description, CategoryID: categoryID})
	}
	return courses, nil
}

func (c *Course) FindByCategoryID(categoryID string) (dto.CourseListOutputDto, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses WHERE category_id = ?", categoryID)
	if err != nil {
		return dto.CourseListOutputDto{}, err
	}
	defer rows.Close()
	courses := dto.CourseListOutputDto{}
	for rows.Next() {
		var id, name, description, categoryID string
		if err := rows.Scan(&id, &name, &description, &categoryID); err != nil {
			return dto.CourseListOutputDto{}, err
		}
		courses.Courses = append(courses.Courses, dto.CourseOutputDto{ID: id, Name: name, Description: description, CategoryID: categoryID})
	}
	return courses, nil
}

func (c *Course) Find(id string) (dto.CourseOutputDto, error) {
	var name, description, categoryID string
	err := c.db.QueryRow("SELECT name, description, category_id FROM courses WHERE id = ?", id).
		Scan(&name, &description, &categoryID)
	if err != nil {
		return dto.CourseOutputDto{}, err
	}
	return dto.CourseOutputDto{ID: id, Name: name, Description: description, CategoryID: categoryID}, nil
}

func (c *Course) Update(course dto.CourseInputDto) error {
	_, err := c.db.Exec("UPDATE courses SET name = ?, description = ?, category_id = ? WHERE id = ?",
		course.Name, course.Description, course.CategoryID, course.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Course) Delete(id string) error {
	_, err := c.db.Exec("DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
