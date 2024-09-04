package controller

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func Ambassador(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.JSON(users)
}

func Admin(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Where("is_ambassador = false").Find(&users)

	return c.JSON(users)
}

func Rankings(c *fiber.Ctx) error {
	// var users []models.User

	// database.DB.Find(&users, models.User{
	// 	IsAmbassador: true,
	// })

	// var result []interface{}

	// for _, user := range users {
	// 	ambassador := models.Ambassador(user)
	// 	ambassador.CalcualteRevenue(database.DB)

	// 	result = append(result, fiber.Map{
	// 		user.Name(): ambassador.Revenue,
	// 	})
	// }
	rangkings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rangkings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return err
	}

	result := make(map[string]float64)

	for _, rangking := range rangkings {
		result[rangking.Member.(string)] = rangking.Score
	}

	return c.JSON(result)
}
