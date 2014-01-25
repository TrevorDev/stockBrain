package database

import (
	"log"
	"./../../lib/config"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

var DB *sql.DB

func InitDatabase(){
	var err error
	log.Println("connection")
	DB, err = sql.Open("mysql", config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	//THIS IS SOO BROKE ADDING THIS CAUSES DEADLOCK I SWEAR ITS A GOLANG BUG
	//DB.SetMaxIdleConns(9)
	//DB.SetMaxOpenConns(9)
}

func GetDatabase()*sql.DB{
	return DB
}