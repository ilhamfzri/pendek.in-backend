package app

import (
	"database/sql"
	"time"

	"github.com/ilhamfzri/pendek.in/internal/helper"
)

func NewDB() *sql.DB {
	/**
	TODO :
	- [] integrate parameters database with config loader
	- [] implement logging
	**/

	connStr := "postgres://postgres:example@localhost:5432/pendekin_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	helper.PanicIfError(err)

	err = db.Ping()
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db
}

func NewDBTest() *sql.DB {
	connStr := "postgres://postgres:example@localhost:5432/pendekin_db_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	helper.PanicIfError(err)

	err = db.Ping()
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db
}
