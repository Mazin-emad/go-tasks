package main

import (
	"database/sql"
	"log"

	"github.com/Mazin-emad/todo-backend/cmd/api"
	"github.com/Mazin-emad/todo-backend/config"
	"github.com/Mazin-emad/todo-backend/db"
	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.NewMysqlStorage(mysql.Config{
		User: config.ConfigAmigoo.DBUser,
		Passwd: config.ConfigAmigoo.DBPassword,
		Net: "tcp",
		Addr: config.ConfigAmigoo.DBAddress,
		DBName: config.ConfigAmigoo.DBName,
		AllowNativePasswords: true,
		ParseTime: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewApiServer(":4001", db)
	if err := server.Run(); err !=nil {
		log.Fatal(err)
	} 
}

func initStorage(db *sql.DB) {
  err := db.Ping()
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Connected to the database")
}