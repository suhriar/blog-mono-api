package mysql

import (
	"context"
	"database/sql"

	"github.com/suhriar/blog-mono-api/model"
)

func (r *userRepository) GetUser(ctx context.Context, email, username string, userID int64) (user model.User, err error) {
	query := `SELECT id, email, password, username, created_at, updated_at, created_by, updated_by
	FROM users WHERE email = ? OR username = ? OR id = ?`
	row := r.db.QueryRowContext(ctx, query, email, username, userID)

	err = row.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return
	}

	return
}

func (r *userRepository) CreateUser(ctx context.Context, model model.User) (lastInsertID int64, err error) {
	query := `INSERT INTO users (email, password, username, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, model.Email, model.Password, model.Username, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return
	}

	lastInsertID, err = res.LastInsertId()
	if err != nil {
		return
	}

	return
}
