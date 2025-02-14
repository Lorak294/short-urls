package routes

import (
	"os"
	"short-urls/constants"
	"short-urls/contracts"
	"short-urls/database"
	"short-urls/helpers"
	"short-urls/validation"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler for the Shorten endpoint.
func ShortenURL(ctx *fiber.Ctx) error {

	// parse request body
	req := new(contracts.ShortenRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": constants.ERROR_BODY_PARSE})
	}

	// validate request body
	if err := validation.ValidateShortenRequest(*req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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

	// enforce https, SSL
	req.Url = helpers.EnforceHTTP(req.Url)

	// generate short url id
	var id string
	if req.CustomShort == "" {
		id = uuid.New().String()
	} else {
		id = req.CustomShort
	}

	// check if given short url is already in use
	targetUrl, _ := dbClient.ResolveShortUrl(id)
	if targetUrl != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error" : constants.ERROR_SHORT_IN_USE})
	}

	
	// set the url mapping in the db
	_ , err =  dbClient.CreateShortForUrl(id,req.Url,req.Expiry * time.Hour)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error" : constants.ERROR_SERVER_GENERAL_ERROR})
	}
	
	// set the expiry for the url
	if req.Expiry == 0 {
		req.Expiry = constants.DEFAULT_EXPIRY_DURATION
	}

	// prepare the response
	resp := contracts.ShortenResponse{
		Url: 				req.Url,
		CustomShort: 		"",
		Expiry: 			req.Expiry,
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