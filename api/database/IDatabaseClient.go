package database

import "time"

type IDatabaseClient interface {

	// state handling
	Initialize()
	Close()

	// short urls
	ResolveShortUrl(url string) (string, error)
	CreateShortForUrl(short string,url string,ttl time.Duration)  (string, error)

	// rate limiting
	GetRateLimitForIp(ip string) (int, error)
	SetRateLimitForIp(ip string, limit int, ttl time.Duration) error
	DecrementRateLimitForIp(ip string) (error)
	GetLeftRateLimitTime(ip string) (time.Duration, error)
}