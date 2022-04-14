package main

import (
	"os"

	"github.com/Ethical-Ralph/quik_wallet/config"
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/cache"
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database"
	"github.com/Ethical-Ralph/quik_wallet/repository"
	"github.com/Ethical-Ralph/quik_wallet/server"
	"github.com/Ethical-Ralph/quik_wallet/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	config := config.Load()

	redis, err := cache.ConnectRedis(&redis.Options{
		Addr:     config.REDIS_ADDRESS,
		Password: config.REDIS_PASSWORD,
		DB:       0,
	})
	if err != nil {
		panic(err)
	}

	databaseConn, err := database.Connect(config.DATABASE_URL)
	if err != nil {
		panic(err)
	}

	mysqlRepo := repository.NewMySQLRepository(databaseConn)
	repository := repository.NewRepository(mysqlRepo)

	cache := cache.NewCache(redis)
	service := service.NewService(repository, cache)

	server := server.NewServer(service, router, logger)

	server.Run()
}
