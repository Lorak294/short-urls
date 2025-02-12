package database

import "time"

type DatabaseClient interface {

	// state handling
	Initialize() error
	Close() error

	// short urls
	ResolveShortUrl(url string) string

	// rate limiting
	GetRateLimitForIp(ip string) (int, error)
	SetRateLimitForIp(ip string, limit int, ttl time.Duration) error
	DecrementRateLimitForIp(ip string) (error)
}