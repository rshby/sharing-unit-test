package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

func CreateConnectionDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:@tcp(%v:%v)/%v?parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))

	if err != nil {
		log.Fatalf("error cant connect to database")
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(3 * time.Hour)
	db.SetConnMaxIdleTime(1 * time.Hour)

	log.Println("success connection to database")
	return db
}
