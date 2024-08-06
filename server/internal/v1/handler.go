package v1

import (
	"encoding/json"
	"net/http"

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
	if err != nil || post.Caption == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repository.CreatePost(&post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
