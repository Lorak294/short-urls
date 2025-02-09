package routes

import (
	"short-urls/database"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Handler for the Resolve endpoint.
func ResolveURL(ctx *fiber.Ctx) error {

	// get the url from params
	url := ctx.Params("url")
	
	// create database client and close it in the end
	dbClient := database.CreateDatabaseClient(database.SHORT_URLS_DB_NR)
	defer dbClient.Close()

	// retreive the value for the urlkey
	targetUrl, err := dbClient.Get(database.Ctx,url).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error" : "shortUrl was not found.",
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Something went wrong.",
		})
	}

	// create client for counter db
	dbClientCounter := database.CreateDatabaseClient(database.COUNTER_DB_NR)
	defer dbClientCounter.Close()
	
	// increase the counter
	_ = dbClientCounter.Incr(database.Ctx,"counter")

	// redirect to targetUrl
	return ctx.Redirect(targetUrl,301);
}