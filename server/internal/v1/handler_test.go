package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/maximilianhagelstam/speek/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepository struct {
	internal.Repository
	getPosts      func() ([]internal.Post, error)
	createPost    func(internal.Post) error
	deletePost    func(primitive.ObjectID) error
	createComment func(internal.Comment) error
	deleteComment func(primitive.ObjectID, primitive.ObjectID) error
	createLike    func(internal.Like) error
	deleteLike    func(primitive.ObjectID, primitive.ObjectID) error
}

func (m mockRepository) GetPosts() ([]internal.Post, error)     { return m.getPosts() }
func (m mockRepository) CreatePost(p internal.Post) error       { return m.createPost(p) }
func (m mockRepository) DeletePost(id primitive.ObjectID) error { return m.deletePost(id) }
func (m mockRepository) CreateComment(c internal.Comment) error { return m.createComment(c) }
func (m mockRepository) DeleteComment(postID, commentID primitive.ObjectID) error {
	return m.deleteComment(postID, commentID)
}
func (m mockRepository) CreateLike(l internal.Like) error { return m.createLike(l) }
func (m mockRepository) DeleteLike(postID, likeID primitive.ObjectID) error {
	return m.deleteLike(postID, likeID)
}

func TestGetPostsHandler(t *testing.T) {
	testObjectID := primitive.NewObjectID()

	tests := []struct {
		name             string
		mockGetPosts     func() ([]internal.Post, error)
		wantErr          bool
		expectedStatus   int
		expectedResponse []internal.Post
	}{
		{
			name: "Happy path",
			mockGetPosts: func() ([]internal.Post, error) {
				return []internal.Post{{ID: testObjectID, Audio: "Test Post"}}, nil
			},
			wantErr:          false,
			expectedStatus:   http.StatusOK,
			expectedResponse: []internal.Post{{ID: testObjectID, Audio: "Test Post"}},
		},
		{
			name: "Repository error",
			mockGetPosts: func() ([]internal.Post, error) {
				return nil, fmt.Errorf("database error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mockRepository{getPosts: tt.mockGetPosts}
			h := NewHandler(r)

			req, _ := http.NewRequest("GET", "/api/v1/posts", nil)
			rr := httptest.NewRecorder()

			h.GetPostsHandler(rr, req)

			if tt.wantErr {
				if rr.Code != http.StatusInternalServerError {
					t.Errorf("expected status 500, got %v", rr.Code)
				}
				return
			}

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}

			var res []internal.Post
			_ = json.NewDecoder(rr.Body).Decode(&res)

			if !reflect.DeepEqual(res, tt.expectedResponse) {
				t.Errorf("expected response %v, got %v", tt.expectedResponse, res)
			}
		})
	}
}

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postData       map[string]string
		mockCreatePost func(internal.Post) error
		expectedStatus int
	}{
		{
			name:           "Happy path",
			postData:       map[string]string{"audio": "Test Post"},
			mockCreatePost: func(p internal.Post) error { return nil },
			expectedStatus: http.StatusCreated,
		},
		{
			name:     "Missing audio",
			postData: map[string]string{"invalid_key": "Test Post"},
			mockCreatePost: func(p internal.Post) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Repository error",
			postData: map[string]string{"audio": "Test Post"},
			mockCreatePost: func(p internal.Post) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mockRepository{createPost: tt.mockCreatePost}
			h := NewHandler(r)

			jsonData, _ := json.Marshal(tt.postData)
			req, _ := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(jsonData))
			rr := httptest.NewRecorder()

			h.CreatePostHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeletePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockDeletePost func(primitive.ObjectID) error
		expectedStatus int
	}{
		{
			name:   "Happy path",
			postID: primitive.NewObjectID().Hex(),
			mockDeletePost: func(id primitive.ObjectID) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "Repository error",
			postID: primitive.NewObjectID().Hex(),
			mockDeletePost: func(id primitive.ObjectID) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepository{deletePost: tt.mockDeletePost}
			h := NewHandler(repo)

			req, _ := http.NewRequest("DELETE", "/api/v1/posts/"+tt.postID, nil)
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Delete("/api/v1/posts/{id}", h.DeletePostHandler)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestCreateCommentHandler(t *testing.T) {
	tests := []struct {
		name              string
		postID            string
		comment           internal.Comment
		mockCreateComment func(internal.Comment) error
		expectedStatus    int
	}{
		{
			name:    "Happy path",
			postID:  primitive.NewObjectID().Hex(),
			comment: internal.Comment{Audio: "Test comment"},
			mockCreateComment: func(c internal.Comment) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:    "Invalid post ID",
			postID:  "123",
			comment: internal.Comment{},
			mockCreateComment: func(c internal.Comment) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Missing audio",
			postID:  primitive.NewObjectID().Hex(),
			comment: internal.Comment{},
			mockCreateComment: func(c internal.Comment) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Repository error",
			postID:  primitive.NewObjectID().Hex(),
			comment: internal.Comment{Audio: "Test comment"},
			mockCreateComment: func(c internal.Comment) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepository{createComment: tt.mockCreateComment}
			h := NewHandler(repo)

			jsonData, _ := json.Marshal(tt.comment)
			req, _ := http.NewRequest("POST", "/api/v1/posts/"+tt.postID+"/comments", bytes.NewBuffer(jsonData))
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/api/v1/posts/{id}/comments", h.CreateCommentHandler)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeleteCommentHandler(t *testing.T) {
	tests := []struct {
		name              string
		postID            string
		commentID         string
		mockDeleteComment func(primitive.ObjectID, primitive.ObjectID) error
		expectedStatus    int
	}{
		{
			name:      "Happy path",
			postID:    primitive.NewObjectID().Hex(),
			commentID: primitive.NewObjectID().Hex(),
			mockDeleteComment: func(postID, commentID primitive.ObjectID) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:      "Invalid post ID",
			postID:    "123",
			commentID: primitive.NewObjectID().Hex(),
			mockDeleteComment: func(postID, commentID primitive.ObjectID) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Invalid comment ID",
			postID:    primitive.NewObjectID().Hex(),
			commentID: "123",
			mockDeleteComment: func(postID, commentID primitive.ObjectID) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Repository error",
			postID:    primitive.NewObjectID().Hex(),
			commentID: primitive.NewObjectID().Hex(),
			mockDeleteComment: func(postID, commentID primitive.ObjectID) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepository{deleteComment: tt.mockDeleteComment}
			h := NewHandler(repo)

			req, _ := http.NewRequest("DELETE", "/api/v1/posts/"+tt.postID+"/comments/"+tt.commentID, nil)
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Delete("/api/v1/posts/{postId}/comments/{commentId}", h.DeleteCommentHandler)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestLikePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockCreateLike func(internal.Like) error
		expectedStatus int
	}{
		{
			name:   "Happy path",
			postID: primitive.NewObjectID().Hex(),
			mockCreateLike: func(l internal.Like) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Invalid post ID",
			postID: "123",
			mockCreateLike: func(l internal.Like) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Repository error",
			postID: primitive.NewObjectID().Hex(),
			mockCreateLike: func(l internal.Like) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mockRepository{createLike: tt.mockCreateLike}
			h := NewHandler(repo)

			req, _ := http.NewRequest("POST", "/api/v1/posts/"+tt.postID+"/like", nil)
			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/api/v1/posts/{id}/like", h.LikePostHandler)
			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}
		})
	}
}
