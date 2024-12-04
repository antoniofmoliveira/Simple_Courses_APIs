package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/flatbuffersapi/internal/configs"
	"github.com/antoniofmoliveira/courses/flatbuffersapi/internal/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	// courseDb := database.NewCourse(db)
	// userDB := database.NewUserRepository(db)

	// // public middlewares
	// public := func(next http.Handler) http.Handler {
	// 	return middleware.Logger(
	// 		middleware.Recoverer(
	// 			middleware.WithValue("jwt", cfg.TokenAuth)(
	// 				middleware.WithValue("jwtExpiresIn", cfg.JWTExpiresIn)(
	// 					next))))
	// }
	// // public middlewares plus verification
	// private := func(next http.Handler) http.Handler {
	// 	return public(
	// 		jwtauth.Verifier(cfg.TokenAuth)(
	// 			jwtauth.Authenticator(
	// 				next)))
	// }
	r := http.NewServeMux()

	categoryHandler := handlers.NewCategoryHandler(categoryDb)
	// courseHandler := handlers.NewCourseHandler(courseDb)
	// userHandler := handlers.NewUserHandler(userDB)

	r.HandleFunc("GET /categories", categoryHandler.FIndAllCategories)
	r.HandleFunc("GET /categories/{id}", categoryHandler.FindCategory)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.WebServerPort),
		Handler: r,
	}

	go func() {
		url := fmt.Sprintf("http://%s:%s", cfg.WebServerHost, cfg.WebServerPort)
		slog.Info("Server", "Server is running at ", url)
		if err := server.ListenAndServe(); err != nil && http.ErrServerClosed != err {
			slog.Error("Could not listen on %s: %v\n", server.Addr, err)
			os.Exit(1)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-termChan
	slog.Info("server: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server", "Could not shutdown the server: %v\n", err.Error())
		os.Exit(1)
	}
	slog.Info("Server stopped")
	os.Exit(0)
}

// curl -H "Accept: application/octet-stream" http://localhost:8080/categories
