package main

import (
	"beerpaws/bot"
	"beerpaws/config"
	"beerpaws/service"
	"beerpaws/storage"
	"log"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	dbConn, err := connectToDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	pointStorage := storage.NewPointsStorage(dbConn)
	pointService := service.NewPointsService(pointStorage)
	userStorage := storage.NewUserStorage(dbConn)
	userService := service.NewUserService(userStorage)

	bot.Run(pointService, userService)
	<-make(chan struct{})
	return
}
