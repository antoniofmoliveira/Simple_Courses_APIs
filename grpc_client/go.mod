module github.com/antoniofmoliveira/courses/grpcclient

go 1.23.4

replace github.com/antoniofmoliveira/courses/db => ../courses_db

replace github.com/antoniofmoliveira/courses => ../courses_entities

replace github.com/antoniofmoliveira/courses/grpcproto => ../proto

require (
	github.com/antoniofmoliveira/courses/grpcproto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.68.0
)

require (
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)
