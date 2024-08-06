package internal

type Repository interface {
	GetPosts() ([]Post, error)
	CreatePost(post Post) error
	DeletePost(postID string) error
	CreateComment(comment Comment) error
	DeleteComment(postID, commentID string) error
	CreateLike(like Like) error
	DeleteLike(postID, likeID string) error
}
