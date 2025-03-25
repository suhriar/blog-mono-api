package model

import "time"

type Post struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	PostTitle    string    `json:"post_title" db:"post_title"`
	PostContent  string    `json:"post_content" db:"post_content"`
	PostHashtags string    `json:"post_hashtags" db:"post_hashtags"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
	UpdatedBy    string    `json:"updated_by" db:"updated_by"`
}

type CreatePostRequest struct {
	PostTitle    string   `json:"postTitle"`
	PostContent  string   `json:"postContent"`
	PostHashtags []string `json:"postHashtags"`
}

type GetAllPostResponse struct {
	Data       []PostDetail `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

type PostDetail struct {
	ID           int64    `json:"id"`
	UserID       int64    `json:"user_id"`
	Username     string   `json:"username"`
	PostTitle    string   `json:"post_title"`
	PostContent  string   `json:"post_content"`
	PostHashtags []string `json:"post_hashtags"`
	IsLiked      bool     `json:"isLiked"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetPostResponse struct {
	PostDetail PostDetail        `json:"post_detail"`
	LikeCount  int               `json:"like_count"`
	Comments   []CommentResponse `json:"comments"`
}

type CommentResponse struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	Username       string `json:"username"`
	CommentContent string `json:"comment_content"`
}
