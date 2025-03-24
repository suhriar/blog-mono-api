// File: config/database/mysql.go
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/suhriar/blog-mono-api/config"
)

// NewMySQLConnection creates and returns a new MySQL database connection
func NewMySQLConnection(cfg *config.Config) (db *sql.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error occured %+v", r)
		}
	}()
	dbConfig := cfg.MySql

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
