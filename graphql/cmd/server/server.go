package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/graphql/graph"
)

const defaultPort = "8081"

func main() {

	dbi := database.GetDBImplementation()
	categoryDb := dbi.CategoryRepository
	courseDb := dbi.CourseRepository

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CategoryDB: categoryDb,
		CourseDB:   courseDb,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServeTLS(":"+port, "./x509/server_cert.pem", "./x509/server_key.pem", nil))
}
