package main

import (
	"database/sql"
	"log"

	"github.com/alexhool2/TimeCard/cmd/api"
	"github.com/alexhool2/TimeCard/config"
	"github.com/alexhool2/TimeCard/db"
	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DbUser,
		Passwd:               config.Envs.DbPassword,
		Addr:                 config.Envs.DbAddress,
		DBName:               config.Envs.DbName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)

	server := api.NewAPIServer(":8081", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected")
}
