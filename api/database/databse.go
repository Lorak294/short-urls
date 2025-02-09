package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func CreateDatabaseClient(dbNumber int) *redis.Client {
	dbClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB: dbNumber,
	})
	return dbClient
}