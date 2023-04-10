package database

import (
	"task/config"
	"testing"
)

func TestOpenConnection(t *testing.T) {
	conf := config.LoadConfig("C:/Users/nebyk/GolandProjects/task/config.yml") //error occurs if you pass the path as config.yml, ../config.yml etc
	db := OpenConnection(conf)
	defer CloseConnection(db)
	err := db.Ping()
	if err != nil {
		t.Errorf("Error pinging database: %v", err)
	}
}

func TestCloseConnection(t *testing.T) {
	conf := config.LoadConfig("C:/Users/nebyk/GolandProjects/task/config.yml") //error occurs if you pass the path as config.yml, ../config.yml etc
	db := OpenConnection(conf)
	CloseConnection(db)
	err := db.Ping()
	if err == nil {
		t.Errorf("Error must be not nil if db connection has been closed: %v", err)
	}
}
