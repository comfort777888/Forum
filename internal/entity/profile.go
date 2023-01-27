package entity

type Profile struct {
	User             UserModel
	ProfileUser      UserModel
	Post             Post
	PostsCount       int
	Posts            []Post
	PostLikes        []string
	PostDislikes     []string
	Comments         []Comments
	CommentsLikes    map[int][]string
	CommentsDislikes map[int][]string
}
