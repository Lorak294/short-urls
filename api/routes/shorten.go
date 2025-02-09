package routes

import (
	"os"
	"short-urls/database"
	"short-urls/helpers"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
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
	XRateLimitRest	time.Duration	`json:"rate_limit_reset"`
}

// Handler for the Shorten endpoint.
func ShortenURL(ctx *fiber.Ctx) error {

	// parse request body
	body := new(request)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse JSON"})
	}

	// rate limiting - checking the user ip
	dbClient := database.CreateDatabaseClient(database.COUNTER_DB_NR)
	defer dbClient.Close()
	if err := CheckRateLimit(dbClient,ctx); err != nil {
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

	// decrement the database key for ip
	dbClient.Decr(database.Ctx,ctx.IP())
	
	return nil
}

func CheckRateLimit(dbClient *redis.Client, ctx *fiber.Ctx ) error {
	
	counterForIpStr, err := dbClient.Get(database.Ctx,ctx.IP()).Result()
	if err == redis.Nil {
		// no ip in db -> set up new quota record
		quotaDuration, err := strconv.Atoi(os.Getenv("API_QUOTA_RESET_TIME"))
		if  err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Something went wrong."})
		}

		_ = dbClient.Set(
			database.Ctx,
			ctx.IP(),
			os.Getenv("API_QUOTA"),
			time.Duration(quotaDuration) * time.Minute).Err()
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
				"rate_limit_rest": limit / time.Nanosecond / time.Minute,
			})
		}

	}
	return nil
}