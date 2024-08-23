package controller

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var product models.Product

func AllProducts(c *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	database.DB.Find(&product)

	return c.JSON(product)
}

// func GetProductByUserId(c *fiber.Ctx) error {
// 	id, _ := strconv.Atoi(c.Params("id"))

// 	product.Id = uint(id)

// 	database.DB.Find(&product)

// 	return c.JSON(product)
// }

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var product models.Product

	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	database.DB.Delete(&product)

	return c.JSON(fiber.Map{
		"Message": "Delete Complete",
	})
}
