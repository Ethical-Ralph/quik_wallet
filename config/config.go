package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DATABASE_URL   string
	REDIS_ADDRESS  string
	REDIS_PASSWORD string
}

func Load() *config {
	c := &config{}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		log.Fatalf("DATABASE_URL MISSING IN ENV FILE")

	}
	redis_addr := os.Getenv("REDIS_ADDRESS")
	if redis_addr == "" {
		log.Fatalf("REDIS_ADDRESS MISSING IN ENV FILE")

	}

	redis_password := os.Getenv("REDIS_PASSWORD")

	c.DATABASE_URL = db_url
	c.REDIS_ADDRESS = redis_addr
	c.REDIS_PASSWORD = redis_password

	return c
}
