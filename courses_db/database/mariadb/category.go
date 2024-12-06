package mariadb

import (
	"database/sql"

	"errors"

	"github.com/antoniofmoliveira/courses/dto"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	c := &CategoryRepository{db: db}
	c.db.Exec("CREATE TABLE IF NOT EXISTS categories (id CHAR(36) PRIMARY KEY, name TEXT, description TEXT)")
	return c
}

func (c *CategoryRepository) Create(categoryDto dto.CategoryInputDto) (dto.CategoryOutputDto, error) {
	// func (c *Category) Create(name string, description string) (dto.CategoryOutputDto, error) {

	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES (?, ?, ?)",
		id, categoryDto.Name, categoryDto.Description)
	if err != nil {
		return dto.CategoryOutputDto{}, err
	}
	return dto.CategoryOutputDto{ID: id, Name: categoryDto.Name, Description: categoryDto.Description}, nil
}

func (c *CategoryRepository) FindAll() (dto.CategoryListOutputDto, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return dto.CategoryListOutputDto{}, err
	}
	defer rows.Close()
	categories := dto.CategoryListOutputDto{}
	for rows.Next() {
		var id, name, description string
		if err := rows.Scan(&id, &name, &description); err != nil {
			return dto.CategoryListOutputDto{}, err
		}
		categories.Categories = append(categories.Categories, dto.CategoryOutputDto{ID: id, Name: name, Description: description})
	}
	return categories, nil
}

func (c *CategoryRepository) FindByCourseID(courseID string) (dto.CategoryOutputDto, error) {
	var id, name, description string
	err := c.db.QueryRow("SELECT c.id, c.name, c.description FROM categories c JOIN courses co ON c.id = co.category_id WHERE co.id = ?", courseID).
		Scan(&id, &name, &description)
	if err != nil {
		return dto.CategoryOutputDto{}, err
	}
	return dto.CategoryOutputDto{ID: id, Name: name, Description: description}, nil
}

func (c *CategoryRepository) Find(id string) (dto.CategoryOutputDto, error) {
	var name, description string
	err := c.db.QueryRow("SELECT name, description FROM categories WHERE id = ?", id).
		Scan(&name, &description)
	if err != nil {
		return dto.CategoryOutputDto{}, err
	}
	return dto.CategoryOutputDto{ID: id, Name: name, Description: description}, nil
}

func (c *CategoryRepository) Update(category dto.CategoryInputDto) error {
	_, err := c.db.Exec("UPDATE categories SET name = ?, description = ? WHERE id = ?",
		category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepository) Delete(id string) error {
	query, err := c.db.Prepare("select count(*) from courses where category_id = ?")
	if err != nil {
		return err
	}
	var count int
	err = query.QueryRow(id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("category has courses")
	}
	_, err = c.db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
