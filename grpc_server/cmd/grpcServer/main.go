package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/grpcproto/pb"
	"github.com/antoniofmoliveira/courses/grpcserver/internal/service"
	"github.com/antoniofmoliveira/courses/grpcserver/internal/service/configs"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName) 
	db, err := sql.Open(cfg.DBDriver,conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDb)

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
