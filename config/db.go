package config

import (
	"database/sql"
	"log"
	"os"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp("+os.Getenv("DB_ADDRESS")+":3306)/"+os.Getenv("DB_TABLE_NAME")+"?parseTime=true&loc=Asia%2FJakarta")

	if err != nil {
		log.Fatal(err)
	}
	return db
}
