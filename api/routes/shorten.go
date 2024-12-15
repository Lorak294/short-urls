package routes

import "time"

type request struct {
	Url				string			`json:"url"`
	CustomShort		string			`json:"short"`
	Expiry			time.Duration	`json:"expiry"`
}

type response struct {
	Url				string			`json:"url"`
	CustomShort		string			`json:"short"`
	Expiry			time.Duration	`json:"expiry"`
	XRateRemaining	int				`json:"rate_limit"`
	XRateLimitRest	time.Duration	`json:"rate_limit_reset"`
}