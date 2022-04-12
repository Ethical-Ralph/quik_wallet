package main

import (
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/cache"
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database"
	"github.com/Ethical-Ralph/quik_wallet/repository"
	"github.com/Ethical-Ralph/quik_wallet/server"
	"github.com/Ethical-Ralph/quik_wallet/service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	redis, err := cache.ConnectRedis()
	if err != nil {
		panic(err)
	}

	databaseConn, err := database.Connect()
	if err != nil {
		panic(err)
	}
	mysqlRepo := repository.NewMySQLRepository(databaseConn)

	cache := cache.NewCache(redis)
	repository := repository.NewRepository(mysqlRepo)

	service := service.NewService(&repository, &cache)

	server := server.NewServer(service, router)

	server.Run()
}
