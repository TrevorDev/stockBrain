package database

import (
	"log"
	"./../../lib/config"
	_ "github.com/lib/pq"
	"database/sql"
)

var DB *sql.DB

func InitDatabase(){
	var err error
	DB, err = sql.Open("postgres", config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDatabase()*sql.DB{
	return DB
}