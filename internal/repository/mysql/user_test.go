package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/suhriar/blog-mono-api/model"
)

// Test GetUser
func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &userRepository{db: db}

	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"
	userID := int64(1)

	mockUser := model.User{
		ID:        userID,
		Email:     email,
		Password:  "hashedpassword",
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "admin",
		UpdatedBy: "admin",
	}

	// Mock DB response
	rows := sqlmock.NewRows([]string{"id", "email", "password", "username", "created_at", "updated_at", "created_by", "updated_by"}).
		AddRow(mockUser.ID, mockUser.Email, mockUser.Password, mockUser.Username, mockUser.CreatedAt, mockUser.UpdatedAt, mockUser.CreatedBy, mockUser.UpdatedBy)

	mock.ExpectQuery(`SELECT id, email, password, username, created_at, updated_at, created_by, updated_by FROM users`).
		WithArgs(email, username, userID).
		WillReturnRows(rows)

	// Execute function
	user, err := repo.GetUser(ctx, email, username, userID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test CreateUser
func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &userRepository{db: db}

	ctx := context.Background()
	mockUser := model.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Username:  "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: "admin",
		UpdatedBy: "admin",
	}

	// Mocking insert result
	mock.ExpectExec(`INSERT INTO users`).
		WithArgs(mockUser.Email, mockUser.Password, mockUser.Username, mockUser.CreatedAt, mockUser.UpdatedAt, mockUser.CreatedBy, mockUser.UpdatedBy).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulasi last insert ID = 1

	// Execute function
	lastID, err := repo.CreateUser(ctx, mockUser)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(1), lastID)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
