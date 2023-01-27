package entity

type Likes struct {
	LikeAuthor  string
	LikedPostId int
	Liked       int // bool?
	Disliked    int
}
