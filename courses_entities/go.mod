module github.com/antoniofmoliveira/courses

go 1.23.4

replace github.com/antoniofmoliveira/courses/db => ../courses_db

replace github.com/antoniofmoliveira/courses => ../courses_entities

require (
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.28.0
)
