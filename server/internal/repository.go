package internal

type Repository interface {
	GetPosts() ([]Post, error)
	CreatePost(post *Post) error
}
