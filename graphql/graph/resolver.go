package graph

import "github.com/antoniofmoliveira/courses/db/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CategoryDB database.CategoryRepositoryInterface
	CourseDB   database.CourseRepositoryInterface
}
