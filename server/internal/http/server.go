package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/maximilianhagelstam/speek/internal/database"
	v1 "github.com/maximilianhagelstam/speek/internal/v1"
)

type Server struct {
	handler *v1.Handler
	port    int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db := database.New()
	repository := v1.NewRepository(db)
	handler := v1.NewHandler(repository)

	NewServer := &Server{
		port:    port,
		handler: handler,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/posts", s.handler.GetPostsHandler)
	r.Post("/api/v1/posts", s.handler.CreatePostHandler)
	r.Delete("/api/v1/posts/{id}", s.handler.DeletePostHandler)
	r.Post("/api/v1/posts/{id}/comments", s.handler.CreateCommentHandler)
	r.Delete("/api/v1/posts/{postId}/comments/{commentId}", s.handler.DeleteCommentHandler)

	return r
}
