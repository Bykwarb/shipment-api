package database

import (
	"database/sql"
	"fmt"
	"log"

	"task/iternal/model"

	_ "github.com/lib/pq"
)

func OpenConnection(config *model.Config) *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Panic(err)
	}

	log.Println("database connection open")
	return db
}

func CloseConnection(db *sql.DB) {
	err := db.Close()

	if err != nil {
		log.Panic(err)
	}

	log.Println("database connection closed")
}
