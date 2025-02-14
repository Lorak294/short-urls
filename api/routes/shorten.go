package routes

import (
	"os"
	"short-urls/constants"
	"short-urls/database"
	"short-urls/helpers"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	XRateLimitReset	time.Duration	`json:"rate_limit_reset"`
}

// Handler for the Shorten endpoint.
func ShortenURL(ctx *fiber.Ctx) error {

	// parse request body
	body := new(request)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": constants.ERROR_BODY_PARSE})
	}

	// get API quota constants
	API_QUOTA, err := strconv.Atoi(os.Getenv("API_QUOTA"))
	if  err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": constants.ERROR_SERVER_GENERAL_ERROR})
	}
	API_QUOTA_RESET, err := strconv.Atoi(os.Getenv("API_QUOTA_RESET_TIME"))
	if  err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": constants.ERROR_SERVER_GENERAL_ERROR})
	}

	// rate limiting - checking the user ip
	dbClient := database.CreateDatabaseClient()
	defer dbClient.Close()
	err_code, err := helpers.CheckRateLimit(dbClient,ctx.IP())
	if err != nil {
		return ctx.Status(err_code).JSON(fiber.Map{"error":err.Error()})
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.Url) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":constants.ERROR_INVALID_URL})
	}

	// check for domain error 
	if !helpers.RemoveDomainError(body.Url){
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":constants.ERROR_FORBIDDEN_URL})
	}

	// enforce https, SSL
	body.Url = helpers.EnforceHTTP(body.Url)

	// generate short url id
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()
	} else {
		id = body.CustomShort
	}

	// check if given short url is already in use
	targetUrl, _ := dbClient.ResolveShortUrl(id)
	if targetUrl != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error" : constants.ERROR_SHORT_IN_USE})
	}

	// set the expiry for the url
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// set the url mapping in the db
	_ , err =  dbClient.CreateShortForUrl(id,body.Url,body.Expiry * time.Hour)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error" : constants.ERROR_SERVER_GENERAL_ERROR})
	}

	// prepare the response
	resp := response{
		Url: 				body.Url,
		CustomShort: 		"",
		Expiry: 			body.Expiry,
		XRateRemaining: 	API_QUOTA,
		XRateLimitReset: 	time.Duration(API_QUOTA_RESET)*time.Minute,
	}

	// decrement the database key for ip
	dbClient.DecrementRateLimitForIp(ctx.IP())
	
	// update the response rate limit fileds
	resp.XRateRemaining, _ =  dbClient.GetRateLimitForIp(ctx.IP())
	ttl, _ := dbClient.GetLeftRateLimitTime(ctx.IP())
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	// update the response custom url
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return ctx.Status(fiber.StatusOK).JSON(resp)
}