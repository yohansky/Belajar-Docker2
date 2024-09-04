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

	var result []interface{}

	// for _, user := range users {
	// 	ambassador := models.Ambassador(user)
	// 	ambassador.CalcualteRevenue(database.DB)

	// 	result = append(result, fiber.Map{
	// 		user.Name(): ambassador.Revenue,
	// 	})
	// }

	rangkings, err := database.Cache.ZRangeByScoreWithScores(context.Background(), "rangkings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return err
	}

	// result := make(map[string]float64)

	for _, rangking := range rangkings {
		member, ok := rangking.Member.(string)
		if !ok {
			c.Status(fiber.StatusInternalServerError)

			return c.JSON(fiber.Map{
				"Message": "failed to convert member to string!",
			})
		}
		// result[member] = rangking.Score
		result = append(result, fiber.Map{
			member: rangking.Score,
		})
	}

	return c.JSON(result)
}
