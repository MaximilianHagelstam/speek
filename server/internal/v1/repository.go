package v1

import (
	"context"

	"github.com/maximilianhagelstam/speek/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ internal.Repository = &Repository{}

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPosts() ([]internal.Post, error) {
	opts := options.Find().SetSort(bson.D{{Key: "date_ordered", Value: 1}})
	cursor, err := r.db.Collection("posts").Find(context.Background(), bson.D{{}}, opts)
	if err != nil {
		return []internal.Post{}, err
	}

	posts := []internal.Post{}
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repository) CreatePost(post *internal.Post) error {
	_, err := r.db.Collection("posts").InsertOne(context.Background(), *post)
	if err != nil {
		return err
	}
	return nil
}
