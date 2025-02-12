package database

type DatabaseClient interface {
	Close() error
	ResolveShortUrl(url string) string
	GetRateLimitForIp(ip string) int
	IncrementRateLimitForIp(ip string) int
}