package database

import (
	"database/sql"
	"fmt"
	"log"
	"task/config"

	_ "github.com/lib/pq"
)

func OpenConnection(config *config.Config) *sql.DB {
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

	log.Println("Database connection open")
	return db
}

func CloseConnection(db *sql.DB) {
	err := db.Close()

	if err != nil {
		log.Panic(err)
	}

	log.Println("Database connection closed")
}
