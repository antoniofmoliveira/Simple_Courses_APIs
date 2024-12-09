module github.com/antoniofmoliveira/courses/jsonapiclient

go 1.23.4

replace github.com/antoniofmoliveira/courses/db => ../courses_db

replace github.com/antoniofmoliveira/courses => ../courses_entities

require github.com/antoniofmoliveira/courses v0.0.0-00010101000000-000000000000
