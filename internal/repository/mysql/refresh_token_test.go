package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/suhriar/blog-mono-api/model"
)

func TestInsertRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &userRepository{db: db}

	ctx := context.Background()
	refreshToken := model.RefreshToken{
		UserID:       1,
		RefreshToken: "random_token",
		ExpiredAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    "system",
		UpdatedBy:    "system",
	}

	mock.ExpectExec(`INSERT INTO refresh_tokens`).
		WithArgs(refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiredAt, refreshToken.CreatedAt, refreshToken.UpdatedAt, refreshToken.CreatedBy, refreshToken.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(1, 1))

	lastInsertID, err := repo.InsertRefreshToken(ctx, refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastInsertID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &userRepository{db: db}

	ctx := context.Background()
	userID := int64(1)
	now := time.Now()

	expectedToken := model.RefreshToken{
		ID:           1,
		UserID:       userID,
		RefreshToken: "random_token",
		ExpiredAt:    now.Add(24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    "system",
		UpdatedBy:    "system",
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expired_at", "created_at", "updated_at", "created_by", "updated_by"}).
		AddRow(expectedToken.ID, expectedToken.UserID, expectedToken.RefreshToken, expectedToken.ExpiredAt, expectedToken.CreatedAt, expectedToken.UpdatedAt, expectedToken.CreatedBy, expectedToken.UpdatedBy)

	mock.ExpectQuery(`SELECT id, user_id, refresh_token, expired_at, created_at, updated_at, created_by, updated_by FROM refresh_tokens WHERE user_id = \? AND expired_at >= \?`).
		WithArgs(userID, now).
		WillReturnRows(rows)

	token, err := repo.GetRefreshToken(ctx, userID, now)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	assert.NoError(t, mock.ExpectationsWereMet())
}
