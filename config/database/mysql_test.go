package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suhriar/blog-mono-api/config"
	"github.com/suhriar/blog-mono-api/config/database"
)

func TestNewMySQLConnection_Integration(t *testing.T) {
	cfg := &config.Config{
		MySql: config.MySqlConfig{
			User:     "root",
			Password: "root",
			Host:     "localhost",
			Port:     "3306",
			Name:     "user-db",
		},
	}

	db, err := database.NewMySQLConnection(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Pastikan koneksi bisa digunakan
	err = db.Ping()
	assert.NoError(t, err)

	// Tutup koneksi setelah selesai
	defer db.Close()
}
