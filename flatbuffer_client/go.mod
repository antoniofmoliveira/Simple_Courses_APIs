module github.com/antoniofmoliveira/courses/flatbufferclient

go 1.23.4

replace github.com/antoniofmoliveira/courses/db => ../courses_db

replace github.com/antoniofmoliveira/courses => ../courses_entities

replace github.com/antoniofmoliveira/courses/flatbuffersapi => ../flatbuffer_api

require (
	github.com/antoniofmoliveira/courses/flatbuffersapi v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20241204233417-43b7b7cde48d
)

require github.com/google/flatbuffers v24.3.25+incompatible // indirect
