package mock

import "github.com/maximilianhagelstam/speek/internal"

var _ internal.Repository = &Repository{}

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetPosts() ([]internal.Post, error) {
	return []internal.Post{{ID: "acb123", Caption: "Test"}}, nil
}

func (r *Repository) CreatePost(post *internal.Post) error {
	return nil
}
