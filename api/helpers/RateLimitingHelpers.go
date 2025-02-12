package helpers

import (
	"errors"
	"os"
	"short-urls/constants"
	"short-urls/database"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Cheks if the given ip has rate limit left to use and returns http status code and error containing error message
func CheckRateLimit(dbClient *redis.Client, senderIp string) (int,error)  {
	// Get env variables
	API_QUOTA, err := strconv.Atoi(os.Getenv("API_QUOTA"))
	if  err != nil {
		return fiber.StatusInternalServerError,errors.New(constants.ERROR_SERVER_GENERAL_ERROR)
	}
	API_QUOTA_RESET, err := strconv.Atoi(os.Getenv("API_QUOTA_RESET_TIME"))
	if  err != nil {
		return fiber.StatusInternalServerError,errors.New(constants.ERROR_SERVER_GENERAL_ERROR)
	}

	counterForIpStr, err := dbClient.Get(database.Ctx,senderIp).Result()
	if err == redis.Nil {
		// no ip in db -> set up new quota record
		err = dbClient.Set(
			database.Ctx,
			senderIp,
			API_QUOTA,
			time.Duration(API_QUOTA_RESET) * time.Minute).Err()
		if  err != nil {
			return fiber.StatusInternalServerError,errors.New(constants.ERROR_SERVER_GENERAL_ERROR)
		}
	} else {
		// ip in db -> check the quota retreived for the ip
		counterForIp, err := strconv.Atoi(counterForIpStr)
		if  err != nil {
			return fiber.StatusInternalServerError,errors.New(constants.ERROR_SERVER_GENERAL_ERROR)
		}
		if counterForIp <= 0 {
			
			//limit, _ := dbClient.TTL(database.Ctx,senderIp).Result()
			return fiber.StatusServiceUnavailable,errors.New(constants.ERROR_RATE_LIMIT_EXCEEDED)
		}

	}
	return fiber.StatusOK,nil
}