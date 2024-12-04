package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/antoniofmoliveira/courses/grpcproto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// with authentication
	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile("x509/ca_cert.pem", "grpccourses")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(creds))

	// without authentication
	// conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCategoryServiceClient(conn)

	// list categories
	fmt.Println("### list categories")
	listCategories(c)

	// create category
	fmt.Println("### create category")
	category := createCategory(c)

	// get category
	fmt.Println("### get category")
	getCategory(c, category.Id)

	// create category with stream
	fmt.Println("### create category with stream")
	createCategoriesWithStream(c)

	// list categories
	fmt.Println("### list categories")
	listCategories(c)

	// create category with bidirectional streams
	fmt.Println("### create category with bidirectional streams")
	createCategoriesWithBidirecionalStream(c)

	// list categories
	fmt.Println("### list categories")
	listCategories(c)

}

func createCategoriesWithBidirecionalStream(c pb.CategoryServiceClient) {
	newCategories := []*pb.CreateCategoryRequest{
		{Name: "Teste 5", Description: "Teste 5"},
		{Name: "Teste 6", Description: "Teste 6"},
		{Name: "Teste 7", Description: "Teste 7"},
	}
	stream2, err := c.CreateCategoryStreamBidirectional(context.Background())
	if err != nil {
		log.Fatalf("could not create category stream: %v", err)
	}

	for _, category := range newCategories {
		err = stream2.Send(category)
		if err == io.EOF {
			stream2.CloseSend()
		}
		if err != nil {
			log.Fatalf("could not send category: %v", err)
		}

		category2, err := stream2.Recv()
		if err != nil {
			log.Fatalf("could not receive category: %v", err)
		}
		log.Printf("category: %v", category2)

	}
}

func createCategoriesWithStream(c pb.CategoryServiceClient) []*pb.CreateCategoryRequest {

	newCategories := []*pb.CreateCategoryRequest{
		{Name: "Teste 1", Description: "Teste 1"},
		{Name: "Teste 2", Description: "Teste 2"},
		{Name: "Teste 3", Description: "Teste 3"},
	}
	stream, err := c.CreateCategoryStream(context.Background())
	if err != nil {
		log.Fatalf("could not create category stream: %v", err)
	}

	for _, category := range newCategories {
		err = stream.Send(category)
		if err != nil {
			log.Fatalf("could not send category: %v", err)
		}
	}
	categories, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive categories: %v", err)
	}
	log.Printf("categories: %v", categories)
	return newCategories
}

func getCategory(c pb.CategoryServiceClient, id string) {
	category, err := c.GetCategory(context.Background(), &pb.CategoryGetRequest{Id: id})
	if err != nil {
		log.Fatalf("could not get category: %v", err)
	}
	log.Printf("category: %v", category)
}

func createCategory(c pb.CategoryServiceClient) *pb.Category {
	category, err := c.CreateCategory(context.Background(), &pb.CreateCategoryRequest{Name: "Teste CLI", Description: "Teste CLI Desc"})
	if err != nil {
		log.Fatalf("could not create category: %v", err)
	}
	log.Printf("category: %v", category)
	return category
}

func listCategories(c pb.CategoryServiceClient) *pb.CategoryList {
	categories, err := c.ListCategories(context.Background(), &pb.Blank{})
	if err != nil {
		log.Fatalf("could not list categories: %v", err)
	}
	for _, category := range categories.Categories {
		log.Printf("category: %v", category)
	}
	return categories
}
