package config

import (
	"database/sql"
	"log"
	"os"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME")+"?parseTime=true&loc=Asia%2FJakarta")

	if err != nil {
		log.Fatal(err)
	}
	return db
}
