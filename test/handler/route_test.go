package handler_test

import (
	"TestITMX/handler"
	"TestITMX/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http/httptest"
	"testing"
)

// เป็น func ที่ setup ตัว db ที่รันบน Memory ขึ้นมา
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	db.AutoMigrate(&model.Customers{})
	return db
}

func TestHandleGetCustomerByID(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()

	handler.SetupRoutes(app, db)

	err := model.InitData(db) //เป็น func ที่สร้างไว้ initData หากว่าไม่มีข้อมูลบน DB
	assert.NoError(t, err, "initData() failed")

	t.Run("Get customer by id success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/customers/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Fail when getting customer with non-existing id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/customers/999", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestHandleCreateCustomer(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()

	handler.SetupRoutes(app, db)

	t.Run("Create customer successfully", func(t *testing.T) {
		customer1 := &model.Customers{
			Name: "TestHandler",
			Age:  20,
		}
		customer, err := json.Marshal(customer1)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", "/customers", bytes.NewBuffer(customer))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("Fail when creating customer with no body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/customers", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestHandleUpdateCustomer(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	handler.SetupRoutes(app, db)

	newCustomer := &model.Customers{
		Name: "TestUpdateHandler",
		Age:  20,
	}
	err := model.AddCustomers(db, newCustomer)
	assert.NoError(t, err)
	updateCustomer := &model.Customers{
		Name: "UpdatedHandler",
		Age:  21,
	}
	t.Run("Update customer successfully", func(t *testing.T) {
		body, err := json.Marshal(updateCustomer)
		req := httptest.NewRequest("PUT", "/customers/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Get updated customer by id success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/customers/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		var updatedCustomer model.Customers
		err = db.First(&updatedCustomer, 1).Error
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedHandler", updatedCustomer.Name)
		assert.Equal(t, uint(21), updatedCustomer.Age)
	})

	t.Run("Fail when updating customer with non-existing id", func(t *testing.T) {
		body, err := json.Marshal(updateCustomer)
		req := httptest.NewRequest("PUT", "/customers/999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}

func TestHandleDeleteCustomer(t *testing.T) {
	db := setupTestDB()
	app := fiber.New()
	handler.SetupRoutes(app, db)

	newCustomer := &model.Customers{
		Name: "TestUpdateHandler",
		Age:  20,
	}
	err := model.AddCustomers(db, newCustomer)
	assert.NoError(t, err)

	t.Run("Delete customer successfully", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/customers/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)
	})

	t.Run("Test get user after delete should return 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/customers/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("fail when deleting not exist customer should return 404", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/customers/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})
}
