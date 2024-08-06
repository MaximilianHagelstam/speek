package mock

import "github.com/maximilianhagelstam/speek/internal"

var _ internal.Repository = &Repository{}

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetPosts() ([]internal.Post, error) {
	return []internal.Post{{ID: "acb123", Audio: "Test"}}, nil
}

func (r *Repository) CreatePost(audio string) error {
	return nil
}

func (r *Repository) DeletePost(postID string) error {
	return nil
}
