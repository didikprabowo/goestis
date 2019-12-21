package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/k/config"
	"time"
)

var (
	db *sql.DB
)

// GetConnection
func GetConnection() *sql.DB {

	co := config.GetConfig()
	c := co.DB
	urlConnection := c.URLConfig()

	db, err := sql.Open("mysql", urlConnection)

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(1000)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err != nil {
		fmt.Println("gagal konek mysql", err.Error())
	}

	return db
}
