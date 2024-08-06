package internal

type Post struct {
	ID    string `bson:"_id,omitempty"`
	Audio string `bson:"audio,omitempty" json:"audio"`
}
