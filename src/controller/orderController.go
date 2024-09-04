package controller

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"fmt"
	"net/smtp"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	// "github.com/stripe/stripe-go/v79"
	// "github.com/stripe/stripe-go/v79/checkout/session"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	link := models.Link{
		Code: request.Code,
	}

	database.DB.Preload("User").First(&link)

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"Message": "Invalid Link!",
		})
	}

	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	// making temporary table
	// transaction
	tx := database.DB.Begin()

	// database.DB.Create(&order)
	if err := tx.Create(&order).Error; err != nil {
		// everything insert in transaction will not added to db
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"Message": err.Error(),
		})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, requestProduct := range request.Products {
		product := models.Product{}
		product.Id = uint(requestProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(requestProduct["quantity"])

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}
		if err := tx.Create(&item).Error; err != nil {
			// everything insert in transaction will not added to db
			tx.Rollback()
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"Message": err.Error(),
			})
		}

		// v79
		// lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
		// 	PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
		// 		Currency:   stripe.String("usd"),
		// 		UnitAmount: stripe.Int64(100 * int64(product.Price)),
		// 		ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
		// 			Name:        stripe.String(product.Title),
		// 			Description: stripe.String(product.Description),
		// 			Images:      []*string{stripe.String(product.Image)},
		// 		},
		// 	},
		// 	Quantity: stripe.Int64(int64(requestProduct["quantity"])),
		// })

		// v72
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(100 * int64(product.Price)),
			Currency:    stripe.String("usd"),
			Quantity:    stripe.Int64(int64(requestProduct["quantity"])),
		})
	}

	if len(lineItems) == 0 {
		tx.Rollback()
		return fmt.Errorf("no items in the cart")
	}

	stripe.Key = "sk_test_51PuD4L07mJ6xZcEwq69cvdNHWv2YvH58rQ8DSKYliGj8c8qv2Bg7nVtTK83aVJ8ZJzifveSckmmxgH7QNVkoV1RM00bnr3hxeQ"

	params := &stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String("http://localhost:5000/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String("payment"),
	}

	source, err := session.New(params)

	if err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"Message": err.Error(),
		})
	}

	order.TransactionId = source.ID

	if err := tx.Save(&order).Error; err != nil {
		// everything insert in transaction will not added to db
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"Message": err.Error(),
		})
	}

	// memasukkan data ke database (dari temporary tadi ke DB)
	tx.Commit()

	return c.JSON(order)
}

func CompleteOrder(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	order := models.Order{}

	database.DB.Preload("OrderItems").First(&order, models.Order{
		TransactionId: data["source"],
	})

	if order.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"Message": "order not found",
		})
	}

	order.Complete = true
	database.DB.Save(&order)

	go func(order models.Order) {
		ambassadorRevenue := 0.0
		adminRevenue := 0.0

		for _, item := range order.OrderItems {
			ambassadorRevenue += item.AmbassadorRevenue
			adminRevenue += item.AdminRevenue
		}

		user := models.User{}
		user.Id = order.UserId

		database.DB.First(&user)

		database.Cache.ZIncrBy(context.Background(), "rangkings", ambassadorRevenue, user.Name())

		// Sending Email, "host.docker.internal:1025" = localhost on docker
		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", ambassadorRevenue, order.Code))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{order.AmbassadorEmail}, ambassadorMessage)

		// Sending Email, "host.docker.internal:1025" = localhost on docker
		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed", order.Id, adminRevenue))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)

	}(order)

	return c.JSON(fiber.Map{
		"Message": "Success",
	})
}
