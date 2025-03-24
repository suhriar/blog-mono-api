package mysql

import (
	"context"
	"strings"

	"github.com/suhriar/blog-mono-api/model"
)

func (r *postRepository) CreatePost(ctx context.Context, model model.Post) (lastInsertID int64, err error) {
	query := `INSERT INTO posts(user_id, post_title, post_content, post_hashtags, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, model.UserID, model.PostTitle, model.PostContent, model.PostHashtags, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return
	}

	lastInsertID, err = res.LastInsertId()
	if err != nil {
		return
	}

	return
}

func (r *postRepository) GetAllPost(ctx context.Context, limit, offset int) (resp model.GetAllPostResponse, err error) {
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags 
	FROM posts p JOIN users u ON p.user_id = u.id ORDER BY p.updated_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	data := []model.PostDetail{}
	for rows.Next() {
		var (
			post     model.Post
			username string
		)
		err = rows.Scan(&post.ID, &post.UserID, &username, &post.PostTitle, &post.PostContent, &post.PostHashtags)
		if err != nil {
			return
		}
		data = append(data, model.PostDetail{
			ID:           post.ID,
			UserID:       post.UserID,
			Username:     username,
			PostTitle:    post.PostTitle,
			PostContent:  post.PostContent,
			PostHashtags: strings.Split(post.PostHashtags, ","),
		})
	}
	resp.Data = data
	resp.Pagination = model.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	return
}

func (r *postRepository) GetPostByID(ctx context.Context, id int64) (resp model.PostDetail, err error) {
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags, uv.is_liked 
	FROM posts p JOIN users u ON p.user_id = u.id 
	JOIN user_activities uv ON uv.post_id = p.id 
	WHERE p.id = ?`

	var (
		post     model.Post
		username string
		isLiked  bool
	)
	row := r.db.QueryRowContext(ctx, query, id)

	err = row.Scan(&post.ID, &post.UserID, &username, &post.PostTitle, &post.PostContent, &post.PostHashtags, &isLiked)
	if err != nil {
		return
	}

	resp = model.PostDetail{
		ID:           post.ID,
		UserID:       post.UserID,
		Username:     username,
		PostTitle:    post.PostTitle,
		PostContent:  post.PostContent,
		PostHashtags: strings.Split(post.PostHashtags, ","),
		IsLiked:      isLiked,
	}
	return
}
