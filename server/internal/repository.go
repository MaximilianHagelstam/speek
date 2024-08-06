package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

type Repository interface {
	GetPosts() ([]Post, error)
	CreatePost(post Post) error
	DeletePost(postID primitive.ObjectID) error
	CreateComment(comment Comment) error
	DeleteComment(postID, commentID primitive.ObjectID) error
	CreateLike(like Like) error
	DeleteLike(postID, likeID primitive.ObjectID) error
}
