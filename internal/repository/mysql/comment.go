package mysql

import (
	"context"

	"github.com/suhriar/blog-mono-api/model"
)

func (r *postRepository) CreateComment(ctx context.Context, model model.Comment) (lastInsertID int64, err error) {
	query := `INSERT INTO comments(post_id, user_id, comment_content, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, model.PostID, model.UserID, model.CommentContent, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return
	}

	lastInsertID, err = res.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (r *postRepository) GetCommentsByPostID(ctx context.Context, postID int64) (comments []model.CommentResponse, err error) {
	query := `SELECT c.id, c.user_id, c.comment_content, u.username 
	FROM comments c JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			comment  model.CommentResponse
			username string
		)
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.CommentContent, &username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, model.CommentResponse{
			ID:             comment.ID,
			UserID:         comment.UserID,
			CommentContent: comment.CommentContent,
			Username:       username,
		})
	}
	return
}
