package models

type Comment struct {
	CommentId string
	UserId    string
	Likes     int
	Dislikes  int
	Replies   []Comment
}
