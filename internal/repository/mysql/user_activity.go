package mysql

import (
	"context"
	"database/sql"

	"github.com/suhriar/blog-mono-api/model"
)

func (r *postRepository) GetUserActivity(ctx context.Context, model model.UserActivity) (resp model.UserActivity, err error) {
	query := `SELECT id, post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by FROM user_activities WHERE post_id = ? AND user_id = ?`

	row := r.db.QueryRowContext(ctx, query, model.PostID, model.UserID)

	err = row.Scan(&resp.ID, &resp.PostID, &resp.UserID, &resp.IsLiked, &resp.CreatedAt, &resp.UpdatedAt, &resp.CreatedBy, &resp.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}
		return resp, err
	}
	return resp, nil
}

func (r *postRepository) CreateUserActivity(ctx context.Context, model model.UserActivity) (lastInsertID int64, err error) {
	query := `INSERT INTO user_activities (post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, model.PostID, model.UserID, model.IsLiked, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return
	}

	lastInsertID, err = res.LastInsertId()
	if err != nil {
		return
	}

	return
}

func (r *postRepository) UpdateUserActivity(ctx context.Context, req model.UserActivity) (err error) {
	query := `UPDATE user_activities SET is_liked = ?, updated_at = ?, updated_by = ? WHERE post_id = ? AND user_id = ?`
	_, err = r.db.ExecContext(ctx, query, req.IsLiked, req.UpdatedAt, req.UpdatedBy, req.PostID, req.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) CountLikeByPostID(ctx context.Context, postID int64) (count int, err error) {
	query := `SELECT COUNT(id) FROM user_activities WHERE post_id = ? AND is_liked = true`

	row := r.db.QueryRowContext(ctx, query, postID)

	err = row.Scan(&count)
	if err != nil {
		return
	}
	return
}
