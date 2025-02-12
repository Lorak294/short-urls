package database

import "time"

type DatabaseClient interface {

	// state handling
	Initialize() DatabaseClient
	Close()

	// short urls
	ResolveShortUrl(url string) (string, error)

	// rate limiting
	GetRateLimitForIp(ip string) (int, error)
	SetRateLimitForIp(ip string, limit int, ttl time.Duration) error
	DecrementRateLimitForIp(ip string) (error)
	GetLeftRateLimitTime(ip string) (time.Duration, error)
}