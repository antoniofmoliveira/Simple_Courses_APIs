package main

import (
	"context"
	"fmt"
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

	dbi := database.GetDBImplementation()
	categoryRepository := dbi.CategoryRepository
	courseRepository := dbi.CourseRepository
	userRepository := dbi.UserRepository

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

	categoryHandler := handlers.NewCategoryHandler(categoryRepository)
	courseHandler := handlers.NewCourseHandler(courseRepository)
	userHandler := handlers.NewUserHandler(userRepository)

	r.HandleFunc("GET /categories", categoryHandler.FindAllCategories)
	r.HandleFunc("GET /categories/{id}", categoryHandler.FindCategory)
	r.HandleFunc("POST /categories", categoryHandler.CreateCategory)
	r.HandleFunc("PUT /categories/{id}", categoryHandler.UpdateCategory)
	r.HandleFunc("DELETE /categories/{id}", categoryHandler.DeleteCategory)

	r.HandleFunc("GET /courses", courseHandler.FindAllCourses)
	r.HandleFunc("GET /courses/{id}", courseHandler.FindCourse)
	r.HandleFunc("POST /courses", courseHandler.CreateCourse)
	r.HandleFunc("PUT /courses/{id}", courseHandler.UpdateCourse)
	r.HandleFunc("DELETE /courses/{id}", courseHandler.DeleteCourse)

	r.HandleFunc("GET /users", userHandler.FindAllUsers)
	r.HandleFunc("GET /users/{id}", userHandler.FindUser)
	r.HandleFunc("POST /users", userHandler.CreateUser)
	r.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	r.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	r.HandleFunc("GET /jwt", userHandler.GetJWT)

	// TODO! only for test - REMOVE! in production
	r.HandleFunc("GET /categorieserror", categoryHandler.CategoriesError)

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
