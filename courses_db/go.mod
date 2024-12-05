module github.com/antoniofmoliveira/courses/db

go 1.23.4

replace github.com/antoniofmoliveira/courses/db => ../courses_db

replace github.com/antoniofmoliveira/courses => ../courses_entities

require (
	github.com/antoniofmoliveira/courses v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.8.1
	github.com/google/uuid v1.6.0
)

require filippo.io/edwards25519 v1.1.0 // indirect
