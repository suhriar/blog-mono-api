package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/suhriar/blog-mono-api/model"
)

func (r *userRepository) InsertRefreshToken(ctx context.Context, model model.RefreshToken) (lastInsertID int64, err error) {
	query := `INSERT INTO refresh_tokens (user_id, refresh_token, expired_at, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, query, model.UserID, model.RefreshToken, model.ExpiredAt, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return
	}

	lastInsertID, err = res.LastInsertId()
	if err != nil {
		return
	}
	return
}

func (r *userRepository) GetRefreshToken(ctx context.Context, userID int64, now time.Time) (resp model.RefreshToken, err error) {
	query := `SELECT id, user_id, refresh_token, expired_at, created_at, updated_at, created_by, updated_by FROM refresh_tokens WHERE user_id = ? AND expired_at >= ?`

	row := r.db.QueryRowContext(ctx, query, userID, now)
	err = row.Scan(&resp.ID, &resp.UserID, &resp.RefreshToken, &resp.ExpiredAt, &resp.CreatedAt, &resp.UpdatedAt, &resp.CreatedBy, &resp.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, nil
		}
		return resp, err
	}
	return resp, nil
}
