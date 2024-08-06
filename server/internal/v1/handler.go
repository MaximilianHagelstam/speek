package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/maximilianhagelstam/speek/internal"
)

type Handler struct {
	repository internal.Repository
}

func NewHandler(repository internal.Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.repository.GetPosts()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post internal.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil || post.Audio == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repository.CreatePost(post.Audio); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	if postID == "" {
		http.Error(w, "post id is required", http.StatusBadRequest)
		return
	}

	if err := h.repository.DeletePost(postID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
