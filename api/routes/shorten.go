package routes

import (
	"short-urls/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
)

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

// Handler for the Shorten endpoint.
func ShortenURL(ctx *fiber.Ctx) error {

	// parse request body
	body := new(request)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse JSON"})
	}

	// TODO: rate limiting - checking the user ip

	// check if the input is an actual URL
	if !govalidator.IsURL(body.Url) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid Url"})
	}

	// check for domain error 
	if !helpers.RemoveDomainError(body.Url){
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"This Url cannot be shortened."})
	}

	// enforce https, SSL
	body.Url = helpers.EnforceHTTP(body.Url)
}