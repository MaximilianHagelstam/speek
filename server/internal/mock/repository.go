package mock

import (
	"github.com/maximilianhagelstam/speek/internal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ internal.Repository = &Repository{}

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetPosts() ([]internal.Post, error) {
	return []internal.Post{{ID: primitive.NewObjectID(), Audio: "Test"}}, nil
}

func (r *Repository) CreatePost(post internal.Post) error {
	return nil
}

func (r *Repository) DeletePost(postID string) error {
	return nil
}

func (r *Repository) CreateComment(comment internal.Comment) error {
	return nil
}

func (r *Repository) DeleteComment(postID, commentID string) error {
	return nil
}
