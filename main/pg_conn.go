package main

import (
	"beerpaws/config"
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func connectToDB() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBConf.Host, config.DBConf.Port, config.DBConf.User, config.DBConf.Password, config.DBConf.DBName)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, nil
}
