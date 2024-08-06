package internal

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Audio     string             `bson:"audio,omitempty" json:"audio"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Comments  []Comment          `bson:"comments" json:"comments"`
}

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PostID    primitive.ObjectID `bson:"post_id" json:"post_id"`
	Audio     string             `bson:"audio" json:"audio"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
