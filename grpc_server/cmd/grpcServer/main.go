package main

import (
	"log"
	"net"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/grpcproto/pb"
	"github.com/antoniofmoliveira/courses/grpcserver/internal/service"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	dbi := database.GetDBImplementation()
	
	categoryService := service.NewCategoryService(dbi.CategoryRepository)

	// with authentication
	creds, err := credentials.NewServerTLSFromFile("x509/server_cert.pem", "x509/server_key.pem")
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	// without authentication
	// grpcServer := grpc.NewServer()

	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
