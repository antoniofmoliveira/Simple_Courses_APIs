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
	"github.com/antoniofmoliveira/courses/jsonapi/internal/configs"
	"github.com/antoniofmoliveira/courses/jsonapi/internal/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dbi := database.GetDBImplementation()
	categoryDb := dbi.CategoryRepository
	courseDb := dbi.CourseRepository
	userDB := dbi.UserRepository

	// public middlewares
	public := func(next http.Handler) http.Handler {
		return middleware.Logger(
			middleware.Recoverer(
				middleware.WithValue("jwt", cfg.TokenAuth)(
					middleware.WithValue("jwtExpiresIn", cfg.JWTExpiresIn)(
						next))))
	}
	// public middlewares plus verification
	private := func(next http.Handler) http.Handler {
		return public(
			jwtauth.Verifier(cfg.TokenAuth)(
				jwtauth.Authenticator(
					next)))
	}
	r := http.NewServeMux()
	categoryHandler := handlers.NewCategoryHandler(categoryDb)
	courseHandler := handlers.NewCourseHandler(courseDb)
	userHandler := handlers.NewUserHandler(userDB)

	r.Handle("GET /categories", private(http.HandlerFunc(categoryHandler.FindAllCategories)))
	r.Handle("GET /categories/{id}", private(http.HandlerFunc(categoryHandler.FindCategory)))
	r.Handle("POST /categories", private(http.HandlerFunc(categoryHandler.CreateCategory)))
	r.Handle("PUT /categories/{id}", private(http.HandlerFunc(categoryHandler.UpdateCategory)))
	r.Handle("DELETE /categories/{id}", private(http.HandlerFunc(categoryHandler.DeleteCategory)))

	r.Handle("GET /courses", private(http.HandlerFunc(courseHandler.FindAllCourses)))
	r.Handle("GET /courses/{id}", private(http.HandlerFunc(courseHandler.FindCourse)))
	r.Handle("POST /courses", private(http.HandlerFunc(courseHandler.CreateCourse)))
	r.Handle("PUT /courses/{id}", private(http.HandlerFunc(courseHandler.UpdateCourse)))
	r.Handle("DELETE /courses/{id}", private(http.HandlerFunc(courseHandler.DeleteCourse)))

	r.Handle("POST /users", private(http.HandlerFunc(userHandler.CreateUser)))
	r.Handle("GET /users", private(http.HandlerFunc(userHandler.FindByEmail)))

	r.Handle("POST /users/generate_token", public(http.HandlerFunc(userHandler.GetJwt)))

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
