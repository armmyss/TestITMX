package unit

import (
	"TestITMX/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

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

func TestInitData(t *testing.T) {
	db := setupTestDB()

	err := model.InitData(db)
	assert.NoError(t, err, "initData() failed")

	var count int64
	db.Model(&model.Customers{}).Count(&count)

	assert.Equal(t, int64(3), count, "Must have 3 customers")
}

func TestGetCustomer(t *testing.T) {
	db := setupTestDB()

	err := model.InitData(db)
	assert.NoError(t, err, "initData() failed")

	t.Run("Get customer by id success", func(t *testing.T) {
		result, err := model.GetCustomer(db, 1)
		assert.NoError(t, err, "getCustomer() failed")
		assert.Equal(t, uint(1), result.ID, "Must be customer id: 1")
		assert.Equal(t, "Arm", result.Name, "Must be customer name: Arm")
		assert.Equal(t, uint(21), result.Age, "Must be customer Age: 21")
	})

	t.Run("Fail when getting customer by does not exist id", func(t *testing.T) {
		_, err := model.GetCustomer(db, 999)
		assert.Error(t, err, "Expected error for non-existing customer")
		assert.Equal(t, "customer not found", err.Error(), "Expected error 'customer not found' when getting non-existing customer")
	})
}

func TestAddCustomer(t *testing.T) {
	db := setupTestDB()

	t.Run("add customer successfully", func(t *testing.T) {
		newCustomer := &model.Customers{
			Name: "Test Customer",
			Age:  99,
		}

		err := model.AddCustomers(db, newCustomer)
		assert.NoError(t, err, "addCustomers() failed")

		var lastCustomer model.Customers
		result := db.Last(&lastCustomer) // จะใช้วิธีดึงข้อมูลล่าสุดเพื่อให้แน่ใจว่า Add เข้าไปแล้วจริงๆ
		assert.NoError(t, result.Error, "Error fetching last customer")

		assert.Equal(t, newCustomer.Name, lastCustomer.Name, "Must be customer name: Test Customer")
		assert.Equal(t, newCustomer.Age, lastCustomer.Age, "Must be customer Age: 99")
	})
}

func TestUpdateCustomer(t *testing.T) {
	db := setupTestDB()

	t.Run("update customer successfully", func(t *testing.T) {
		newCustomer := &model.Customers{
			Name: "Test Customer",
			Age:  99,
		}

		err := model.AddCustomers(db, newCustomer)
		assert.NoError(t, err, "addCustomers() failed")

		var lastCustomer model.Customers
		result := db.Last(&lastCustomer)
		assert.NoError(t, result.Error, "Error fetching last customer")
		assert.Equal(t, newCustomer.Name, lastCustomer.Name, "Must be customer name: Test Customer")
		assert.Equal(t, newCustomer.Age, lastCustomer.Age, "Must be customer Age: 99")

		lastCustomer.Name = "New Customer Name"
		lastCustomer.Age = 88

		err = model.UpdateCustomer(db, &lastCustomer)
		assert.NoError(t, err, "updateCustomer() failed")

		customerUpdated, err := model.GetCustomer(db, lastCustomer.ID)
		assert.NoError(t, err, "updateCustomer() failed")
		assert.Equal(t, lastCustomer.Name, customerUpdated.Name, "Must be customer name: New Customer Name")
		assert.Equal(t, lastCustomer.Age, customerUpdated.Age, "Must be customer Age: 88")
	})

	t.Run("Fail to update non exist customer", func(t *testing.T) {
		nonExistingCustomer := &model.Customers{
			ID:   9999,
			Name: "nonExistingCustomer",
			Age:  99,
		}
		err := model.UpdateCustomer(db, nonExistingCustomer)
		assert.Error(t, err, "Customer not found")
		assert.Equal(t, "customer not found", err.Error(), "customer not found")
	})
}

func TestDeleteCustomer(t *testing.T) {
	db := setupTestDB()
	t.Run("delete customer successfully", func(t *testing.T) {
		newCustomer := &model.Customers{
			Name: "Test Delete Customer",
			Age:  99,
		}
		err := model.AddCustomers(db, newCustomer)
		assert.NoError(t, err, "addCustomers() failed")

		var lastCustomer model.Customers
		result := db.Last(&lastCustomer)
		assert.NoError(t, result.Error, "Error fetching last customer")

		customer, err := model.GetCustomer(db, lastCustomer.ID)
		assert.NoError(t, err, "getCustomer() failed")
		assert.Equal(t, lastCustomer.Name, customer.Name, "Must be customer name: Test Delete Customer")
		assert.Equal(t, lastCustomer.Age, customer.Age, "Must be customer Age: 99")

		err = model.DeleteCustomer(db, lastCustomer.ID)
		assert.NoError(t, err, "deleteCustomer() failed")

		customer, err = model.GetCustomer(db, lastCustomer.ID)
		assert.Error(t, err, "Customer not found")
		assert.Equal(t, "customer not found", err.Error(), "customer not found")
	})

	t.Run("Fail to delete non exist customer", func(t *testing.T) {
		err := model.DeleteCustomer(db, 9999)
		assert.Error(t, err, "Customer not found")
		assert.Equal(t, "customer not found", err.Error(), "customer not found")
	})
}
