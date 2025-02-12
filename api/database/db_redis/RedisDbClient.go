package db_redis

import (
	"context"
	"time"
)

const SHORT_URLS_DB_NR = 0
const COUNTER_DB_NR = 1
var CTX = context.Background()

type RedisDbClient struct {

}

// state handling
func (x RedisDbClient)Initialize() error{

}

func (x RedisDbClient)Close() error {

}

// short urls
func (x RedisDbClient)ResolveShortUrl(url string) string {

}

// rate limiting
func (x RedisDbClient)GetRateLimitForIp(ip string) (int, error) {

}

func (x RedisDbClient)SetRateLimitForIp(ip string, limit int, ttl time.Duration) error{

}

func (x RedisDbClient)DecrementRateLimitForIp(ip string) (error) {

}