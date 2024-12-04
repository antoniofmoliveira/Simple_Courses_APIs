package database

import "github.com/antoniofmoliveira/courses/dto"

type CourseRepositoryInterface interface {
	Create(dto dto.CourseInputDto) (*dto.CourseOutputDto, error)
	FindAll() (dto.CourseListOutputDto, error)
	FindByCategoryID(categoryID string) (dto.CourseListOutputDto, error)
	Find(id string) (dto.CourseOutputDto, error)
	Update(course dto.CourseInputDto) error
	Delete(id string) error
}

type CategoryRepositoryInterface interface {
	Create(dto dto.CategoryInputDto) (dto.CategoryOutputDto, error)
	FindAll() (dto.CategoryListOutputDto, error)
	FindByCourseID(courseID string) (dto.CategoryOutputDto, error)
	Find(id string) (dto.CategoryOutputDto, error)
	Update(category dto.CategoryInputDto) error
	Delete(id string) error
}

type UserRepositoryInterface interface {
	Create(user *dto.CreateUserInput) error
	FindByEmail(email string) (*dto.GetJWTInput, error)
}
