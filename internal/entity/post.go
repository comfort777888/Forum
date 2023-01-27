package entity

import "time"

type Post struct {
	PostId       int
	PostAuthor   string
	Title        string
	Content      string
	CreationTime time.Time
	Category     []string
	Likes        int
	Dislikes     int
}
