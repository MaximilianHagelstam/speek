package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/maximilianhagelstam/speek/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		http.Error(w, "error fetching posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(posts)
}

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post internal.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusInternalServerError)
		return
	}

	if post.Audio == "" {
		http.Error(w, "audio field is required", http.StatusBadRequest)
		return
	}

	if err := h.repository.CreatePost(post); err != nil {
		http.Error(w, "error creating post", http.StatusInternalServerError)
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
		http.Error(w, "error deleting post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	if postID == "" {
		http.Error(w, "post id is required", http.StatusBadRequest)
		return
	}

	var comment internal.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if comment.Audio == "" {
		http.Error(w, "audio is a required field", http.StatusBadRequest)
		return
	}

	comment.PostID, err = primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	comment.ID = primitive.NewObjectID()
	comment.CreatedAt = time.Now()

	if err := h.repository.CreateComment(comment); err != nil {
		http.Error(w, "error creating comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")
	commentID := chi.URLParam(r, "commentId")

	if postID == "" || commentID == "" {
		http.Error(w, "post id and comment id are required fields", http.StatusBadRequest)
		return
	}

	if err := h.repository.DeleteComment(postID, commentID); err != nil {
		http.Error(w, "error deleting comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	if postID == "" {
		http.Error(w, "post id is required", http.StatusBadRequest)
		return
	}

	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	like := internal.Like{
		ID:        primitive.NewObjectID(),
		PostID:    postObjectID,
		CreatedAt: time.Now(),
	}

	if err := h.repository.CreateLike(like); err != nil {
		http.Error(w, "error creating like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UnLikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	if postID == "" {
		http.Error(w, "post id is required", http.StatusBadRequest)
		return
	}

	// TODO: find likeID with current userID

	if err := h.repository.DeleteLike(postID, ""); err != nil {
		http.Error(w, "error deleting like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
