package internal

type Post struct {
	ID      string `bson:"_id,omitempty"`
	Caption string `bson:"caption,omitempty" json:"caption"`
}
