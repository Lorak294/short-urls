package db_redis

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const SHORT_URLS_DB_NR = 0
const COUNTER_DB_NR = 1
var CTX = context.Background()

type RedisDbClient struct {
	shortsDbClient *redis.Client
	rateLimitDbClient *redis.Client
}

// state handling
func (x *RedisDbClient)Initialize() {
	x.shortsDbClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB: SHORT_URLS_DB_NR,
	})
	x.rateLimitDbClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB: SHORT_URLS_DB_NR,
	})
}

func (x *RedisDbClient)Close() {
	x.shortsDbClient.Close()
	x.rateLimitDbClient.Close()
}

// short urls
func (x *RedisDbClient)ResolveShortUrl(url string) (string, error) {
	return x.shortsDbClient.Get(CTX,url).Result()
}
func (x *RedisDbClient)CreateShortForUrl(short string,url string,ttl time.Duration)  (string, error) {
	return x.shortsDbClient.Set(CTX,short,url, ttl).Result()
}

// rate limiting
func (x *RedisDbClient)GetRateLimitForIp(ip string) (int, error) {
	res_str, err := x.rateLimitDbClient.Get(CTX,ip).Result()
	if  err != nil {
		return -1, err
	}
	return strconv.Atoi(res_str)
}

func (x *RedisDbClient)GetLeftRateLimitTime(ip string) (time.Duration, error) {
	res := x.rateLimitDbClient.TTL(CTX,ip)
	return res.Result()
}

func (x *RedisDbClient)SetRateLimitForIp(ip string, limit int, ttl time.Duration) error{
	return x.rateLimitDbClient.Set(CTX,ip,limit,ttl).Err()
}

func (x *RedisDbClient)DecrementRateLimitForIp(ip string) (error) {
	return x.rateLimitDbClient.Decr(CTX,ip).Err()
}