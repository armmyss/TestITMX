package handler

import (
	"TestITMX/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/customers", HandleGetCustomers(db))
	app.Get("/customers/:id", HandleGetCustomerByID(db))
	app.Post("/customers", HandleCreateCustomer(db))
	app.Put("/customers/:id", HandleUpdateCustomer(db))
	app.Delete("/customers/:id", HandleDeleteCustomer(db))
}

func HandleGetCustomers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customers, err := model.GetCustomers(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(customers)
	}
}

func HandleGetCustomerByID(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "must be a number",
			})
		}
		customer, err := model.GetCustomer(db, uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "customer not found",
			})
		}
		return c.JSON(customer)
	}
}

func HandleCreateCustomer(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customer := new(model.Customers)
		if err := c.BodyParser(customer); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err := model.AddCustomers(db, customer)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(customer)
	}
}

func HandleUpdateCustomer(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid customer ID",
			})
		}

		customer := new(model.Customers)
		if err := c.BodyParser(customer); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse request body",
			})
		}

		customer.ID = uint(id)

		err = model.UpdateCustomer(db, customer)
		if err != nil {
			if err.Error() == "customer not found" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Customer not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update customer",
			})
		}

		return c.Status(fiber.StatusOK).JSON(customer)
	}
}

func HandleDeleteCustomer(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = model.DeleteCustomer(db, uint(id))
		if err != nil {
			if err.Error() == "customer not found" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
