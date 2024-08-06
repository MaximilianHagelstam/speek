package internal

type Repository interface {
	GetPosts() ([]Post, error)
	CreatePost(audio string) error
	DeletePost(postID string) error
}
