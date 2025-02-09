package routes

import (
	"os"
	"short-urls/database"
	"short-urls/helpers"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse JSON"})
	}

	// get API quota constants
	API_QUOTA, err := strconv.Atoi(os.Getenv("API_QUOTA"))
	if  err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Something went wrong."})
	}
	API_QUOTA_RESET, err := strconv.Atoi(os.Getenv("API_QUOTA_RESET_TIME"))
	if  err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Something went wrong."})
	}

	// rate limiting - checking the user ip
	dbClientCounter := database.CreateDatabaseClient(database.COUNTER_DB_NR)
	defer dbClientCounter.Close()
	if err := CheckRateLimit(dbClientCounter,ctx,API_QUOTA,API_QUOTA_RESET); err != nil {
		return err;
	}

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

	// generate short url id
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()
	} else {
		id = body.CustomShort
	}

	// create clinet for the url db
	dbClientUrl := database.CreateDatabaseClient(database.SHORT_URLS_DB_NR)
	defer dbClientUrl.Close()

	// check if given short url is already in use
	targetUrl, _ := dbClientUrl.Get(database.Ctx, id).Result()
	if targetUrl != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map {
			"error" : "Given URL custom short is already in use.",
		})
	}

	// set the expiry for the url
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// set the url mapping in the db
	err = dbClientUrl.Set(database.Ctx,id,body.Url, body.Expiry * time.Hour).Err()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"error" : "Somethiing went wrong.",
		})
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
	dbClientCounter.Decr(database.Ctx,ctx.IP())
	
	// update the response rate limit fileds
	val, _ := dbClientCounter.Get(database.Ctx,ctx.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := dbClientCounter.TTL(database.Ctx,ctx.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	// update the response custom url
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func CheckRateLimit(dbClient *redis.Client, ctx *fiber.Ctx, api_quota int, api_quota_reset int) error {
	
	counterForIpStr, err := dbClient.Get(database.Ctx,ctx.IP()).Result()
	if err == redis.Nil {
		// no ip in db -> set up new quota record
		err = dbClient.Set(
			database.Ctx,
			ctx.IP(),
			api_quota,
			time.Duration(api_quota_reset) * time.Minute).Err()
		if  err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Something went wrong."})
		}
	} else {
		// ip in db -> check the quota retreived for the ip
		counterForIp, err := strconv.Atoi(counterForIpStr)
		if  err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Something went wrong."})
		}
		if counterForIp <= 0 {
			limit, _ := dbClient.TTL(database.Ctx,ctx.IP()).Result()

			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":"Request rate limit exceeded.",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}

	}
	return nil
}