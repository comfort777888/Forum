package entity

type Comments struct {
	CommentId int
	PostId    int
	Content   string
	Author    string
	Likes     int
	Dislikes  int
}
