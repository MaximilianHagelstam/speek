package v1

import (
	"context"
	"fmt"

	"github.com/maximilianhagelstam/speek/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *Repository) CreatePost(audio string) error {
	newPost := internal.Post{Audio: audio}
	_, err := r.db.Collection("posts").InsertOne(context.Background(), newPost)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeletePost(postID string) error {
	objectId, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	opts := options.Delete().SetHint(bson.D{{Key: "_id", Value: 1}})
	result, err := r.db.Collection("posts").DeleteOne(context.Background(), filter, opts)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("post id %s doesn't exist", postID)
	}
	return nil
}
