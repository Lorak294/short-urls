package database

import (
	"context"
	"short-urls/database/db_redis"
)


const SHORT_URLS_DB_NR = 0
const COUNTER_DB_NR = 1
var Ctx = context.Background()

func CreateDatabaseClient() IDatabaseClient {
	// current implementation -> RedisDbCllient
	var client db_redis.RedisDbClient
	client.Initialize()
	return &client
}