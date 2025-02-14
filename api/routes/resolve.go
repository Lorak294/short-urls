package routes

import (
	"short-urls/constants"
	"short-urls/database"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Handler for the Resolve endpoint.
func ResolveURL(ctx *fiber.Ctx) error {

	// get the url from params
	url := ctx.Params("url")
	
	// create database client and close it in the end
	dbClient := database.CreateDatabaseClient()
	defer dbClient.Close()

	// retreive the value for the urlkey
	targetUrl, err := dbClient.ResolveShortUrl(url)
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": constants.ERROR_SHORT_RESOLVE_FAIL})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": constants.ERROR_SERVER_GENERAL_ERROR})
	}

	// redirect to targetUrl
	return ctx.Redirect(targetUrl,301);
}